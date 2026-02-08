# Email Notification System Design

## Overview
This document outlines the design for adding email notifications to the AI Agent platform. The system will support transactional emails (welcome, verification) and activity notifications (new followers, mentions) with user-configurable preferences.

## Architecture

```mermaid
graph TD
    Dispatcher[Notification Dispatcher] --> PrefCheck{Check Prefs}
    PrefCheck -->|Push Enabled?| Firebase[Firebase Adapter]
    PrefCheck -->|Email Enabled?| EmailWorker[Async Email Worker]
    
    EmailWorker --> Template[Template Engine]
    Template --> Provider[Email Provider Interface]
    
    Provider --> SMTP[SMTP (Dev)]
    Provider --> SendGrid[SendGrid (Prod)]
    Provider --> SES[AWS SES (Prod)]
```

## Core Components

### 1. Email Provider Interface
Abstracts the email sending mechanism to support multiple providers (SMTP for dev, SendGrid/SES for prod).

```go
package email

type EmailProvider interface {
    Send(ctx context.Context, to []string, subject string, htmlBody string, textBody string) error
}
```

### 2. Email Service
Handles business logic, template rendering, and async dispatch.

```go
package service

type EmailService interface {
    SendNotification(ctx context.Context, userID uuid.UUID, notifType NotificationType, data map[string]interface{}) error
    SendWelcomeEmail(ctx context.Context, userID uuid.UUID, email string, name string) error
    SendVerificationEmail(ctx context.Context, userID uuid.UUID, email string, token string) error
}
```

### 3. Template Engine
Uses Go's `html/template` to render responsive HTML emails.
- **Location**: `internal/infrastructure/email/templates/`
- **Structure**:
    - `layout.html`: Base style (logo, footer, unsubscribe link)
    - `notification.html`: Generic notification template
    - `welcome.html`: Onboarding email

## Integration Strategy

### Notification Dispatcher Update
The `NotificationDispatcher` will be refactored to be channel-agnostic.

```go
// Current
enabled, _ := d.prefRepo.IsEnabled(ctx, userID, notifType, "push")
if enabled { d.firebase.SendPushToUser(...) }

// Proposed
channels := []string{"push", "email"}
for _, channel := range channels {
    enabled, _ := d.prefRepo.IsEnabled(ctx, userID, notifType, channel)
    if !enabled { continue }

    switch channel {
    case "push":
        go d.firebase.SendPushToUser(...)
    case "email":
        go d.emailService.SendNotification(...)
    }
}
```

## Database Schema (No Changes Required)
The existing `notification_preferences` table already supports the `channel` column.

```sql
CREATE TABLE notification_preferences (
    user_id UUID,
    notification_type VARCHAR(50),
    channel VARCHAR(20), -- 'in_app', 'push', 'email'
    enabled BOOLEAN,
    PRIMARY KEY (user_id, notification_type, channel)
);
```

## Implementation Phases

### Phase 1: Foundation
- Define `EmailProvider` interface
- Implement `SMTPSender` (using `gomail` or standard lib)
- Create `EmailService` and template loading logic

### Phase 2: Integration
- Update `NotificationDispatcher` to handle multiple channels
- Add `email` preference check logic
- Inject `EmailService` into `NotificationDispatcher`

### Phase 3: Templates & User Controls
- Create HTML templates (Welcome, New Follower, Mention)
- Update User Settings API to toggle email preferences
- Add Unsubscribe link logic

## Configuration
New environment variables required:

```env
EMAIL_PROVIDER=smtp # or sendgrid, ses
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USER=
SMTP_PASS=
SMTP_FROM=noreply@aiagent.com
SENDGRID_API_KEY= # optional
```

## Future Considerations
- **Digest Emails**: Weekly/Daily summaries to reduce noise
- **Bounce Handling**: Webhook to handle invalid emails
- **Click Tracking**: Analytics for email engagement
