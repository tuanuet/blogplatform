package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
)

const (
	// Notification expiration time
	notificationExpirationHours = 24 * 7 // 1 week
)

// NotificationDispatcher orchestrates notification sending with preferences, aggregation, and push notifications
type NotificationDispatcher interface {
	// Notify sends a notification to a user following the full flow:
	// 1. Check user preferences -> If disabled, return early
	// 2. Check rate limit -> If exceeded, skip
	// 3. Check aggregation -> If similar exists, update
	// 4. Save notification to database
	// 5. Increment rate limit counter
	// 6. Get user's device tokens
	// 7. Send FCM push if channel enabled
	// 8. Handle errors gracefully
	Notify(ctx context.Context, userID uuid.UUID, notifType entity.NotificationType, data map[string]interface{}) error
}

// notificationDispatcher implements the NotificationDispatcher interface
type notificationDispatcher struct {
	notifRepo  repository.NotificationRepository
	prefRepo   repository.NotificationPreferenceRepository
	tokenRepo  repository.DeviceTokenRepository
	aggregator NotificationAggregator
	firebase   FirebaseAdapter
	email      EmailService
	taskRunner TaskRunner
}

// NewNotificationDispatcher creates a new NotificationDispatcher instance
func NewNotificationDispatcher(
	notifRepo repository.NotificationRepository,
	prefRepo repository.NotificationPreferenceRepository,
	tokenRepo repository.DeviceTokenRepository,
	aggregator NotificationAggregator,
	firebase FirebaseAdapter,
	email EmailService,
	taskRunner TaskRunner,
) NotificationDispatcher {
	return &notificationDispatcher{
		notifRepo:  notifRepo,
		prefRepo:   prefRepo,
		tokenRepo:  tokenRepo,
		aggregator: aggregator,
		firebase:   firebase,
		email:      email,
		taskRunner: taskRunner,
	}
}

// Notify sends a notification to a user
func (d *notificationDispatcher) Notify(
	ctx context.Context,
	userID uuid.UUID,
	notifType entity.NotificationType,
	data map[string]interface{},
) error {
	// Step 1: Check user preferences for all channels
	inAppEnabled, err := d.prefRepo.IsEnabled(ctx, userID, notifType, "in_app")
	if err != nil {
		inAppEnabled = true // Fail open
	}

	pushEnabled, err := d.prefRepo.IsEnabled(ctx, userID, notifType, "push")
	if err != nil {
		pushEnabled = true // Fail open
	}

	emailEnabled, err := d.prefRepo.IsEnabled(ctx, userID, notifType, "email")
	if err != nil {
		emailEnabled = true // Fail open
	}

	if !inAppEnabled && !pushEnabled && !emailEnabled {
		return nil
	}

	// Step 2: Check rate limit
	if d.aggregator != nil {
		allowed, err := d.aggregator.CheckRateLimit(ctx, userID, notifType)
		if err != nil {
			log.Printf("Warning: failed to check rate limit for user %s: %v", userID, err)
		} else if !allowed {
			log.Printf("Info: rate limit exceeded for user %s, notification type %s", userID, notifType)
			return nil
		}
	}

	// Step 3: Check aggregation
	var existingNotif *entity.Notification
	if d.aggregator != nil {
		targetID := extractTargetID(data)
		existingNotif, err = d.aggregator.ShouldAggregate(ctx, userID, notifType, targetID)
		if err != nil {
			log.Printf("Warning: aggregator check failed for user %s: %v", userID, err)
		}
	}

	// Step 4: Prepare notification
	notif := d.prepareNotification(userID, notifType, data)
	if existingNotif != nil {
		notif = existingNotif
		d.updateAggregatedNotification(notif, data)
	}

	// Step 5: Save notification to database (In-app)
	if inAppEnabled {
		if err := d.notifRepo.Save(ctx, notif); err != nil {
			return fmt.Errorf("failed to save notification: %w", err)
		}

		// Increment rate limit counter
		if d.aggregator != nil {
			if err := d.aggregator.IncrementRateLimit(ctx, userID, notifType); err != nil {
				log.Printf("Warning: failed to increment rate limit for user %s: %v", userID, err)
			}
		}
	}

	// Step 6: Asynchronous dispatch to other channels
	channels := []string{"push", "email"}
	for _, channel := range channels {
		switch channel {
		case "push":
			if pushEnabled && d.firebase != nil && d.taskRunner != nil {
				d.taskRunner.Submit(func(ctx context.Context) {
					tokens, err := d.tokenRepo.FindByUserID(ctx, userID)
					if err != nil {
						log.Printf("Error: failed to get device tokens for user %s: %v", userID, err)
						return
					}
					if len(tokens) > 0 {
						title := d.generateTitle(notifType, data)
						body := d.generateBody(notifType, data)

						// Enrich data map for Firebase validation
						enrichedData := make(map[string]interface{})
						for k, v := range data {
							enrichedData[k] = v
						}
						if _, ok := enrichedData["category"]; !ok {
							enrichedData["category"] = string(d.getCategory(notifType))
						}
						if _, ok := enrichedData["target_type"]; !ok {
							// Infer target type from notification type if missing
							switch notifType {
							case entity.NotificationTypeBlogLike, entity.NotificationTypeBlogComment, entity.NotificationTypeNewBlogFromFollowing:
								enrichedData["target_type"] = "blog"
							case entity.NotificationTypeCommentReply:
								enrichedData["target_type"] = "comment"
							case entity.NotificationTypeNewFollower, entity.NotificationTypeMention:
								enrichedData["target_type"] = "user"
							default:
								enrichedData["target_type"] = "user"
							}
						}
						// Ensure target_id is present for validation
						if _, ok := enrichedData["target_id"]; !ok {
							enrichedData["target_id"] = userID.String()
						}

						if err := d.firebase.SendPushToUser(ctx, userID, title, body, enrichedData); err != nil {
							log.Printf("Error: failed to send FCM push for user %s: %v", userID, err)
						}
					}
				})
			}
		case "email":
			if emailEnabled && d.email != nil && d.taskRunner != nil {
				d.taskRunner.Submit(func(ctx context.Context) {
					if err := d.email.SendNotification(ctx, userID, notifType, data); err != nil {
						log.Printf("Error: failed to send email notification for user %s: %v", userID, err)
					}
				})
			}
		}
	}

	return nil
}

// prepareNotification creates a new notification entity
func (d *notificationDispatcher) prepareNotification(
	userID uuid.UUID,
	notifType entity.NotificationType,
	data map[string]interface{},
) *entity.Notification {
	return &entity.Notification{
		UserID:       userID,
		Type:         notifType,
		Category:     d.getCategory(notifType),
		Title:        d.generateTitle(notifType, data),
		Body:         d.generateBody(notifType, data),
		Data:         data,
		GroupedCount: 1,
		IsRead:       false,
		ExpiresAt:    time.Now().Add(notificationExpirationHours * time.Hour),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// updateAggregatedNotification updates an existing notification with new data
func (d *notificationDispatcher) updateAggregatedNotification(notif *entity.Notification, data map[string]interface{}) {
	notif.GroupedCount++
	notif.UpdatedAt = time.Now()

	if notif.Data == nil {
		notif.Data = make(map[string]interface{})
	}

	// Update body to reflect grouping
	actorName := ""
	if name, ok := data["actor_name"].(string); ok {
		actorName = name
	}
	notif.Body = d.generateGroupedBody(actorName, notif.GroupedCount, notif.Type)
}

// getCategory returns the notification category based on type
func (d *notificationDispatcher) getCategory(notifType entity.NotificationType) entity.NotificationCategory {
	switch notifType {
	case entity.NotificationTypeNewFollower,
		entity.NotificationTypeMention:
		return entity.NotificationCategorySocial
	case entity.NotificationTypeBlogLike,
		entity.NotificationTypeBlogComment,
		entity.NotificationTypeCommentReply,
		entity.NotificationTypeNewBlogFromFollowing,
		entity.NotificationTypeSeriesUpdate:
		return entity.NotificationCategoryContent
	case entity.NotificationTypeBotFollowerDetected,
		entity.NotificationTypeBadgeStatusChange:
		return entity.NotificationCategorySystem
	default:
		return entity.NotificationCategorySocial
	}
}

// generateTitle generates notification title based on type
func (d *notificationDispatcher) generateTitle(notifType entity.NotificationType, data map[string]interface{}) string {
	switch notifType {
	case entity.NotificationTypeNewFollower:
		return "New Follower"
	case entity.NotificationTypeBlogLike:
		return "New Like"
	case entity.NotificationTypeBlogComment:
		return "New Comment"
	case entity.NotificationTypeCommentReply:
		return "Reply to Comment"
	case entity.NotificationTypeMention:
		return "You were mentioned"
	case entity.NotificationTypeNewBlogFromFollowing:
		return "New Blog"
	case entity.NotificationTypeSeriesUpdate:
		return "Series Update"
	case entity.NotificationTypeBotFollowerDetected:
		return "Bot Follower Detected"
	case entity.NotificationTypeBadgeStatusChange:
		return "Badge Status Changed"
	default:
		return "Notification"
	}
}

// generateBody generates notification body based on type and data
func (d *notificationDispatcher) generateBody(notifType entity.NotificationType, data map[string]interface{}) string {
	actorName := ""
	if name, ok := data["actor_name"].(string); ok {
		actorName = name
	}

	switch notifType {
	case entity.NotificationTypeNewFollower:
		return actorName + " started following you"
	case entity.NotificationTypeBlogLike:
		blogTitle := ""
		if title, ok := data["blog_title"].(string); ok {
			blogTitle = title
		}
		return actorName + " liked your blog: " + blogTitle
	case entity.NotificationTypeBlogComment:
		blogTitle := ""
		if title, ok := data["blog_title"].(string); ok {
			blogTitle = title
		}
		return actorName + " commented on your blog: " + blogTitle
	case entity.NotificationTypeCommentReply:
		return actorName + " replied to your comment"
	case entity.NotificationTypeMention:
		return actorName + " mentioned you"
	case entity.NotificationTypeNewBlogFromFollowing:
		blogTitle := ""
		if title, ok := data["blog_title"].(string); ok {
			blogTitle = title
		}
		return actorName + " published a new blog: " + blogTitle
	case entity.NotificationTypeSeriesUpdate:
		seriesName := ""
		if name, ok := data["series_name"].(string); ok {
			seriesName = name
		}
		return actorName + " updated series: " + seriesName
	case entity.NotificationTypeBotFollowerDetected:
		return "We detected a bot follower that may have been following you"
	case entity.NotificationTypeBadgeStatusChange:
		status := ""
		if s, ok := data["badge_status"].(string); ok {
			status = s
		}
		return "Your badge status has changed to: " + status
	default:
		return "You have a new notification"
	}
}

// generateGroupedBody generates body for aggregated notifications
func (d *notificationDispatcher) generateGroupedBody(actorName string, count int, notifType entity.NotificationType) string {
	if count <= 1 {
		return ""
	}

	otherCount := count - 1
	if otherCount == 1 {
		switch notifType {
		case entity.NotificationTypeBlogLike:
			return actorName + " and 1 other liked your blog"
		case entity.NotificationTypeNewFollower:
			return actorName + " and 1 other started following you"
		default:
			return actorName + " and 1 other"
		}
	}
	return actorName + " and " + fmt.Sprintf("%d", otherCount) + " others"
}

// extractTargetID extracts target ID from notification data
func extractTargetID(data map[string]interface{}) uuid.UUID {
	if data == nil {
		return uuid.Nil
	}
	if idStr, ok := data["target_id"].(string); ok {
		if id, err := uuid.Parse(idStr); err == nil {
			return id
		}
	}
	return uuid.Nil
}
