# Subscription Feature

## Overview
The Subscription feature manages the social and economic relationships between users. It unifies the concept of "Following" (free updates) and "Subscribing" (paid tier-based access) into a single coherent system.

## Architecture
- **UseCase**: `SubscriptionUseCase` provides high-level operations for relationship management.
- **Domain Service**: `SubscriptionService` enforces rules (e.g., self-subscription prevention, duplicate checks) and manages persistence.
- **Entities**:
  - `Subscription`: The runtime relationship record.
  - `SubscriptionPlan`: Configuration for paid tiers defined by content creators.

## Core Logic & Features
### Relationship Types
The system distinguishes between two primary relationship states using the same `Subscription` entity:
1.  **Follower (Free)**: A standard social follow. `ExpiresAt` is null, and `Tier` is `FREE`.
2.  **Subscriber (Paid)**: A premium relationship. `ExpiresAt` is set to a future date, and `Tier` indicates the level (Bronze, Silver, Gold).

### Lifecycle Management
- **Subscribe**: Establishes a connection. Used for "Following". Paid subscriptions are typically initiated via the Payment feature but result in a `Subscription` record update.
- **Unsubscribe**: Terminates the relationship.
- **Expiration**: Paid subscriptions automatically degrade or expire based on the `ExpiresAt` timestamp.

### Plan Management
Content creators can define `SubscriptionPlan`s (Bronze, Silver, Gold) with specific prices and durations, allowing them to monetize their content.

## Data Model

### Subscription
Represents the link between a Subscriber and an Author.
```go
type Subscription struct {
    ID           uuid.UUID
    SubscriberID uuid.UUID
    AuthorID     uuid.UUID
    Tier         string     // "FREE", "BRONZE", "SILVER", "GOLD"
    ExpiresAt    *time.Time // Null for free followers
    CreatedAt    time.Time
}
```

### SubscriptionPlan
Defines a tier available for purchase.
```go
type SubscriptionPlan struct {
    ID           uuid.UUID
    AuthorID     uuid.UUID
    Tier         SubscriptionTier // Enum: BRONZE, SILVER, GOLD
    Price        decimal.Decimal
    DurationDays int
    IsActive     bool
}
```

## API Reference (Internal)
### SubscriptionUseCase
- `Subscribe(ctx, subscriberID, authorID)`: Create a free follow relationship.
- `Unsubscribe(ctx, subscriberID, authorID)`: Remove relationship.
- `IsSubscribed(ctx, subscriberID, authorID)`: Check if link exists.
- `GetSubscriptions(ctx, subscriberID, page, size)`: List who a user follows.
- `GetSubscribers(ctx, authorID, page, size)`: List who follows a user.
- `CountSubscribers(ctx, authorID)`: Get total follower count.
