# Notification Feature

## Overview
The Notification feature manages the lifecycle of user notifications, including delivery, storage, read status tracking, and user preferences. It supports both in-app notifications and push notifications via device tokens.

## Architecture
- **UseCase**: `NotificationUseCase` handles the business logic for retrieving notifications, managing preferences, and registering devices.
- **Repositories**:
  - `NotificationRepository`: Stores and retrieves notification history.
  - `DeviceTokenRepository`: Manages user device tokens for push notifications.
  - `NotificationPreferenceRepository`: Stores user settings for notification delivery.
- **Adapter**: `NotificationAdapter` abstracts the external notification delivery service (e.g., FCM, APNs).

## Core Logic & Features
### Notification Management
- **Listing**: Retrieves paginated lists of notifications for a user (`List`).
- **Read Status**: Tracks whether a notification has been read (`IsRead`). Supports marking individual or all notifications as read.
- **Unread Count**: Quickly retrieves the count of unread notifications for badge displays.

### User Preferences
Users can configure which notifications they receive and through which channels.
- **Granularity**: Settings are per `NotificationType` (e.g., `new_follower`, `blog_like`) and `Channel` (e.g., `in_app`, `email`, `push`).
- **Defaults**: System assumes default enablement if no specific preference exists.

### Device Registration
- **Token Management**: Registers and updates device tokens (`RegisterDeviceToken`) to enable push notifications on mobile/web clients.
- **Platform Tracking**: Stores the platform (iOS, Android, Web) associated with each token.

## Data Model

### Notification
Represents a single notification event.
```go
type Notification struct {
    ID           uuid.UUID
    UserID       uuid.UUID
    Type         NotificationType // e.g., "new_follower", "blog_like"
    Category     NotificationCategory // "social", "content", "system"
    Title        string
    Body         string
    Data         map[string]interface{} // Payload for deep linking/UI
    IsRead       bool
    CreatedAt    time.Time
}
```

### NotificationPreference
User-specific settings.
```go
type NotificationPreference struct {
    UserID           uuid.UUID
    NotificationType NotificationType
    Channel          string // "in_app", "push", "email"
    Enabled          bool
}
```

### UserDeviceToken
Stored tokens for push notification services.
```go
type UserDeviceToken struct {
    UserID      uuid.UUID
    DeviceToken string
    Platform    string
    LastSeenAt  time.Time
}
```

## API Reference (Internal)
### NotificationUseCase
- `List(ctx, userID, page, pageSize)`: Get paginated notifications.
- `GetUnreadCount(ctx, userID)`: Get number of unread items.
- `MarkAsRead(ctx, userID, notificationID)`: Mark specific item as read.
- `MarkAllAsRead(ctx, userID)`: Mark all items as read.
- `GetPreferences(ctx, userID)`: Get user notification settings.
- `UpdatePreferences(ctx, userID, req)`: Update settings.
- `RegisterDeviceToken(ctx, userID, req)`: Register/Update a push token.
