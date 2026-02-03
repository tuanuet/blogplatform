# Refined Spec: Payment Integration (Momo & VNPay) - v2

## User Stories
- **As a User**, I can subscribe to a **specific Author** (e.g., 1 Month) via Momo/VNPay to access their exclusive content.
- **As a User**, I can buy a specific **Series** via Momo/VNPay to own it forever.
- **As a User**, I can **Donate** to an author to support their work.
- **As a System**, I must verify the **integrity** (signature) of callbacks from Payment Gateways.

## Data Model Changes (High Level)
1. **New: `Transaction`**
   - `Type`: `SUBSCRIPTION` (was Membership), `SERIES`, `DONATION`.
   - `TargetID`:
     - If `SUBSCRIPTION` -> `AuthorID`.
     - If `SERIES` -> `SeriesID`.
     - If `DONATION` -> `AuthorID`.
2. **New: `UserSeriesPurchase`**
   - `UserID` (PK)
   - `SeriesID` (PK)
3. **Modified: `Subscription`** (Existing in `internal/domain/entity/subscription.go`)
   - Currently acts as a "Follow" relationship.
   - **Requirement:** Add `ExpiresAt` (Nullable) to this table.
     - `ExpiresAt` is NULL = Free Follower.
     - `ExpiresAt` > NOW = Paid Subscriber.
   - This avoids creating a duplicate `UserAuthorSubscription` table.

## API Contract (Updated)
`POST /api/v1/payments`
- `type`: "SUBSCRIPTION", "SERIES", "DONATION"
- `targetId`: The ID of the Author (for Sub/Donate) or Series.
- `plan`: (Optional) "1_MONTH", "3_MONTHS" (if type=SUBSCRIPTION).
- `amount`: (Optional) Custom amount for DONATION.

## Tech Stack
- **Language:** Go 1.24
- **Framework:** Gin, GORM
- **Database:** PostgreSQL
- **Existing Entities:** `User`, `Series`, `Subscription`
