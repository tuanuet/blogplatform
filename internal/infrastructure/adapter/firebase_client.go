package adapter

import (
	"context"
	"errors"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

// firebaseClient implements FirebaseClient using Firebase FCM SDK
type firebaseClient struct {
	app    *firebase.App
	client *messaging.Client
}

// NewFirebaseClient creates a new Firebase client instance
// If Firebase is disabled in config, returns nil client (no-op behavior)
func NewFirebaseClient(ctx context.Context, projectID, serviceAccountPath string, enabled bool) (FirebaseClient, error) {
	// Return nil client if Firebase is disabled
	if !enabled {
		log.Warn().Msg("Firebase FCM is disabled in config")
		return &noOpFirebaseClient{}, nil
	}

	// Validate required configuration
	if projectID == "" {
		return nil, errors.New("firebase project_id is required when Firebase is enabled")
	}

	var opts []option.ClientOption

	// Load service account from file if path is provided
	if serviceAccountPath != "" {
		// Check if file exists
		if _, err := os.Stat(serviceAccountPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("firebase service account file not found at path: %s", serviceAccountPath)
		}
		opts = append(opts, option.WithCredentialsFile(serviceAccountPath))
	} else {
		// Try to use default credentials (e.g., when running in GCP)
		log.Info().Msg("no service account path provided, using default credentials")
	}

	// Initialize Firebase app
	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		return nil, fmt.Errorf("error initializing Firebase app: %w", err)
	}

	// Initialize FCM client
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing Firebase Messaging client: %w", err)
	}

	log.Info().
		Str("project_id", projectID).
		Msg("Firebase FCM client initialized successfully")

	return &firebaseClient{
		app:    app,
		client: client,
	}, nil
}

// Send sends a push notification via Firebase FCM
func (f *firebaseClient) Send(ctx context.Context, tokens []string, title, body string, data map[string]interface{}) error {
	// Build notification message
	notification := &messaging.Notification{
		Title: title,
		Body:  body,
	}

	// Convert data to map[string]string as required by Firebase
	stringData := convertDataMap(data)

	// Build message for each token
	var messages []*messaging.Message
	for _, token := range tokens {
		message := &messaging.Message{
			Token:        token,
			Notification: notification,
			Data:         stringData,
			Android: &messaging.AndroidConfig{
				Priority: "high",
				Notification: &messaging.AndroidNotification{
					Title: title,
					Body:  body,
					Sound: "default",
				},
			},
			APNS: &messaging.APNSConfig{
				Payload: &messaging.APNSPayload{
					Aps: &messaging.Aps{
						Alert: &messaging.ApsAlert{
							Title: title,
							Body:  body,
						},
						Sound: "default",
					},
				},
			},
		}
		messages = append(messages, message)
	}

	// Send messages individually (Firebase doesn't support sending to multiple tokens in one call for legacy tokens)
	// Use batch sending for better performance when sending to many tokens
	batchSize := 500 // Firebase API limit
	var errors []error

	for i := 0; i < len(messages); i += batchSize {
		end := i + batchSize
		if end > len(messages) {
			end = len(messages)
		}

		batch := messages[i:end]
		batchErrors := f.sendBatch(ctx, batch)
		errors = append(errors, batchErrors...)

		// Log batch results
		if len(batchErrors) > 0 {
			log.Error().
				Int("batch_start", i).
				Int("batch_end", end).
				Int("errors", len(batchErrors)).
				Msg("some messages failed to send")
		}
	}

	// Return error if any send failed
	if len(errors) > 0 {
		return fmt.Errorf("failed to send %d out of %d messages", len(errors), len(messages))
	}

	return nil
}

// sendBatch sends a batch of messages and returns any errors
func (f *firebaseClient) sendBatch(ctx context.Context, messages []*messaging.Message) []error {
	var errors []error

	for _, message := range messages {
		_, err := f.client.Send(ctx, message)
		if err != nil {
			// Log error for individual message
			log.Error().
				Err(err).
				Str("token", message.Token).
				Msg("failed to send FCM message")
			errors = append(errors, err)
		}
	}

	return errors
}

// Close closes the Firebase client
func (f *firebaseClient) Close() error {
	// Firebase app doesn't have explicit Close method
	// Connection pooling is managed internally
	return nil
}

// noOpFirebaseClient is a no-op implementation used when Firebase is disabled
type noOpFirebaseClient struct{}

func (n *noOpFirebaseClient) Send(ctx context.Context, tokens []string, title, body string, data map[string]interface{}) error {
	log.Debug().
		Str("title", title).
		Int("token_count", len(tokens)).
		Msg("Firebase disabled, skipping push notification")
	return nil
}

// convertDataMap converts map[string]interface{} to map[string]string for Firebase
// Values are converted to their string representation
func convertDataMap(data map[string]interface{}) map[string]string {
	result := make(map[string]string, len(data))
	for key, value := range data {
		switch v := value.(type) {
		case string:
			result[key] = v
		case int, int32, int64:
			result[key] = fmt.Sprintf("%d", v)
		case float32, float64:
			result[key] = fmt.Sprintf("%f", v)
		case bool:
			result[key] = fmt.Sprintf("%t", v)
		default:
			// For other types, use fmt.Sprintf to convert
			result[key] = fmt.Sprintf("%v", v)
		}
	}
	return result
}
