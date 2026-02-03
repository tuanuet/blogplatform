# Feature: Notification System (MVP)

## Objective

Xay dung he thong notification real-time cho phep users nhan thong bao ve social interactions (followers, likes, comments) va content updates (new blogs tu authors ho follow), voi kha nang customize preferences chi tiet. MVP tap trung vao in-app notifications su dung Firebase Cloud Messaging.

## Problem Statement

Users hien tai khong co cach nao biet real-time khi co activities lien quan den ho. Dieu nay dan den:
- Giam engagement va response time
- Bo lo content tu authors yeu thich
- Lower retention rates

## Requirements

### Functional

- [ ] Users receive real-time in-app notifications for social interactions (new follower, likes, comments, mentions)
- [ ] Users receive notifications when followed authors publish new content or update series
- [ ] Similar notifications are grouped (e.g., "Linh and 5 others liked your blog")
- [ ] Users can view notification history (last 30 days) with read/unread status
- [ ] Users can mark notifications as read (individual or all)
- [ ] Users can configure notification preferences per type (on/off)
- [ ] Clicking notification navigates to relevant content
- [ ] Quick actions available from notification (follow back, like, etc.)
- [ ] Integrate with existing system alerts (bot detection, badge status)

### Non-Functional

- [ ] Real-time delivery (< 2 seconds latency)
- [ ] Notifications retained for 30 days
- [ ] Handle notification flood gracefully (rate limiting + grouping)
- [ ] Graceful degradation if Firebase unavailable (fallback to polling)

## User Stories

### Story 1: Receive Social Notifications

**As a** user
**I want to** receive real-time notifications when someone follows me, likes my blog, or comments on my content
**So that** I can engage back promptly

**Acceptance Criteria:**

- [ ] Given I am logged in, When someone follows me, Then I see a notification within 2 seconds
- [ ] Given multiple people like my blog within 5 minutes, When I view notifications, Then I see a grouped notification
- [ ] Given I click on a notification, Then I navigate to the relevant content

### Story 2: Receive Content Update Notifications

**As a** user
**I want to** be notified when authors I follow publish new content
**So that** I don't miss content I care about

**Acceptance Criteria:**

- [ ] Given I follow an author, When they publish a new blog, Then I receive a notification
- [ ] Given I subscribe to a series, When a new chapter is added, Then I receive a notification

### Story 3: Manage Notification Preferences

**As a** user
**I want to** control which types of notifications I receive
**So that** I'm not overwhelmed with unwanted alerts

**Acceptance Criteria:**

- [ ] Given I am in Settings, When I navigate to Notifications, Then I see all notification types grouped by category
- [ ] Given I toggle off a notification type, Then I stop receiving that type of notification
- [ ] Given I change preferences, When a relevant event occurs, Then the system respects my preferences

### Story 4: View Notification History

**As a** user
**I want to** see a list of my notifications with read/unread status
**So that** I can catch up on what I missed

**Acceptance Criteria:**

- [ ] Given I click the bell icon, Then I see my notification list with unread count
- [ ] Given I have unread notifications, Then they are visually distinguished
- [ ] Given I click "Mark all as read", Then all notifications are marked read

## Technical Context

### Impacted Services

- `internal/domain/service` - New NotificationService, extend existing
- `internal/domain/entity` - New Notification, NotificationPreference entities
- `internal/domain/repository` - New NotificationRepository
- `internal/application/usecase` - New notification usecases
- `internal/interfaces/http/handler` - New notification handlers
- `internal/infrastructure/adapter` - Firebase adapter

### New Data/Fields

#### notifications table

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | Recipient user ID (FK to users) |
| type | VARCHAR | Notification type enum |
| category | VARCHAR | social, content, system |
| title | VARCHAR | Rendered title |
| body | TEXT | Rendered body |
| data | JSONB | Actor IDs, target ID, target type, etc. |
| grouped_count | INT | 1 if single, >1 if merged |
| is_read | BOOL | Read status |
| created_at | TIMESTAMP | Creation time |
| expires_at | TIMESTAMP | created_at + 30 days |

#### notification_preferences table

| Column | Type | Description |
|--------|------|-------------|
| user_id | UUID | FK to users |
| notification_type | VARCHAR | Type enum |
| channel | VARCHAR | in_app, email, push |
| enabled | BOOL | On/off |

#### user_device_tokens table

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | FK to users |
| device_token | VARCHAR | FCM token |
| platform | VARCHAR | ios, android, web |
| created_at | TIMESTAMP | Creation time |

### Notification Types

| Category | Type | Template Example |
|----------|------|------------------|
| social | new_follower | "{user} da follow ban" |
| social | blog_like | "{user} va 5 nguoi khac thich bai viet cua ban" |
| social | blog_comment | "{user} da comment: \"{preview}...\"" |
| social | comment_reply | "{user} da reply comment cua ban" |
| social | mention | "{user} da nhac den ban trong {context}" |
| content | new_blog_from_following | "{author} da dang bai moi: \"{title}\"" |
| content | series_update | "{author} da cap nhat series \"{title}\"" |
| system | bot_follower_detected | "Phat hien {count} followers dang ngo" |
| system | badge_status_change | "Badge cua ban da thay doi" |

### API Changes

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/notifications | List notifications (paginated) |
| GET | /api/v1/notifications/unread-count | Get unread count |
| PUT | /api/v1/notifications/:id/read | Mark as read |
| PUT | /api/v1/notifications/read-all | Mark all as read |
| GET | /api/v1/notifications/preferences | Get preferences |
| PUT | /api/v1/notifications/preferences | Update preferences |
| POST | /api/v1/devices/token | Register device token for push |

### Dependencies

- Firebase Cloud Messaging (FCM) for push delivery
- Firebase Admin SDK for Go (`firebase.google.com/go/v4`)

## User Flows

### Happy Path - Receive Notification

```
1. User A likes User B's blog
2. System creates notification event
3. Check User B's preferences -> allowed
4. Check for similar recent notifications -> group if exists
5. Save to database
6. Push to Firebase Cloud Messaging
7. Firebase delivers to User B's device
8. Client shows notification (bell badge + toast/popup)
9. User B clicks -> navigates to blog
10. Notification marked as read via API
```

### Happy Path - Manage Preferences

```
1. User opens Settings -> Notifications
2. GET /api/v1/notifications/preferences
3. UI shows notification types with toggles
4. User toggles off "new_follower"
5. PUT /api/v1/notifications/preferences
6. Future new_follower events are not sent to this user
```

## Edge Cases

| Scenario | Handling |
|----------|----------|
| User unfollows after notification sent | Notification van hien thi (historical record) |
| Firebase delivery fails | Retry voi exponential backoff, fallback to in-app |
| User disables notifications | Respect preferences, van store in DB cho history |
| Notification flood (viral content) | Rate limiting + smart grouping |
| Deleted content | Notification van show, navigate to 404 gracefully |

## Out of Scope (Future Phases)

- Email notifications
- Push notifications (browser/mobile) - MVP chi lam in-app + Firebase integration san
- SMS notifications
- Admin system announcements
- Notification templates management (admin)
- A/B testing notification content
- Analytics dashboard

## Open Questions

- FCM configuration: Can setup Firebase project va credentials
- Grouping window: 5 minutes hay configurable?
- Rate limiting threshold: Max bao nhieu notifications/hour per user?

---

**Status**: Approved
**Created**: 2026-02-03
**Source**: /brainstorm session
**Next**: /pipeline (Gatekeeper Phase)
