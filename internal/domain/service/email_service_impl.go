package service

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/aiagent/internal/infrastructure/adapter"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type emailServiceImpl struct {
	userRepo    repository.UserRepository
	provider    adapter.EmailProvider
	taskRunner  TaskRunner
	templateDir string
	templates   map[string]*template.Template
}

// NewEmailServiceImpl creates a new implementation of EmailService
func NewEmailServiceImpl(
	userRepo repository.UserRepository,
	provider adapter.EmailProvider,
	taskRunner TaskRunner,
	templateDir string,
) EmailService {
	s := &emailServiceImpl{
		userRepo:    userRepo,
		provider:    provider,
		taskRunner:  taskRunner,
		templateDir: templateDir,
		templates:   make(map[string]*template.Template),
	}

	// Pre-parse templates
	s.loadTemplates()

	return s
}

func (s *emailServiceImpl) loadTemplates() {
	layoutPath := filepath.Join(s.templateDir, "layout.html")
	files, err := filepath.Glob(filepath.Join(s.templateDir, "*.html"))
	if err != nil {
		log.Error().Err(err).Msg("failed to glob templates")
		return
	}

	for _, file := range files {
		name := filepath.Base(file)
		if name == "layout.html" {
			continue
		}

		tmpl, err := template.ParseFiles(layoutPath, file)
		if err != nil {
			log.Error().Err(err).Str("template", name).Msg("failed to parse template")
			continue
		}
		s.templates[name] = tmpl
	}
}

func (s *emailServiceImpl) SendNotification(ctx context.Context, userID uuid.UUID, notifType entity.NotificationType, data map[string]interface{}) error {
	// Async dispatch as requested
	s.taskRunner.Submit(func(taskCtx context.Context) {
		user, err := s.userRepo.FindByID(taskCtx, userID)
		if err != nil {
			log.Error().
				Err(err).
				Str("user_id", userID.String()).
				Str("notif_type", string(notifType)).
				Msg("failed to find user for email notification")
			return
		}

		if user == nil || user.Email == "" {
			log.Warn().
				Str("user_id", userID.String()).
				Str("notif_type", string(notifType)).
				Msg("user not found or has no email for notification")
			return
		}

		subject := fmt.Sprintf("New Notification: %s", notifType)

		// Map data for template
		tmplData := map[string]interface{}{
			"Message":   data["message"],
			"ActionURL": data["action_url"],
		}

		htmlBody, textBody, err := s.renderTemplate("notification.html", tmplData)
		if err != nil {
			log.Error().
				Err(err).
				Str("user_id", userID.String()).
				Str("notif_type", string(notifType)).
				Msg("failed to render notification email template")
			return
		}

		err = s.provider.Send(taskCtx, []string{user.Email}, subject, htmlBody, textBody)
		if err != nil {
			log.Error().
				Err(err).
				Str("user_id", userID.String()).
				Str("notif_type", string(notifType)).
				Msg("failed to send notification email")
		}
	})

	return nil
}

func (s *emailServiceImpl) SendWelcomeEmail(ctx context.Context, userID uuid.UUID, email string, name string) error {
	subject := "Welcome to AI Agent!"

	tmplData := map[string]interface{}{
		"Name":   name,
		"AppURL": "https://aiagent.com/get-started",
	}

	htmlBody, textBody, err := s.renderTemplate("welcome.html", tmplData)
	if err != nil {
		return fmt.Errorf("failed to render welcome email: %w", err)
	}

	return s.provider.Send(ctx, []string{email}, subject, htmlBody, textBody)
}

func (s *emailServiceImpl) SendVerificationEmail(ctx context.Context, userID uuid.UUID, email string, token string) error {
	subject := "Verify your email address"

	tmplData := map[string]interface{}{
		"VerificationURL": fmt.Sprintf("https://aiagent.com/verify?token=%s", token),
	}

	htmlBody, textBody, err := s.renderTemplate("verification.html", tmplData)
	if err != nil {
		return fmt.Errorf("failed to render verification email: %w", err)
	}

	return s.provider.Send(ctx, []string{email}, subject, htmlBody, textBody)
}

func (s *emailServiceImpl) renderTemplate(tmplName string, data interface{}) (string, string, error) {
	tmpl, ok := s.templates[tmplName]
	if !ok {
		return "", "", fmt.Errorf("template %s not found", tmplName)
	}

	var htmlBody bytes.Buffer
	if err := tmpl.ExecuteTemplate(&htmlBody, "layout", data); err != nil {
		return "", "", fmt.Errorf("failed to execute html template: %w", err)
	}

	var textBody bytes.Buffer
	if err := tmpl.ExecuteTemplate(&textBody, "content", data); err != nil {
		return "", "", fmt.Errorf("failed to execute text template: %w", err)
	}

	return htmlBody.String(), textBody.String(), nil
}
