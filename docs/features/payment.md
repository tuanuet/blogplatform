# Payment Feature

## Overview
The Payment feature facilitates financial transactions within the platform, enabling users to pay for subscriptions, purchase premium content (Series), or make donations. It currently integrates with **SePay** to support popular Vietnamese payment methods like VietQR and bank transfers.

## Architecture
- **UseCase**:
  - `CreatePaymentUseCase`: Handles payment initialization requests.
  - `ProcessWebhookUseCase`: Handles asynchronous payment confirmation callbacks.
- **Service**: Delegates core payment logic and SePay API interaction to `domainService.PaymentService`.
- **Entity**:
  - `Transaction`: Central record of all payment attempts and their states.
  - `UserSeriesPurchase`: Records content ownership upon successful payment.

## Core Logic & Features
### Payment Initialization
The `CreatePayment` operation prepares a transaction:
1.  **Validation**: Verifies user eligibility and request validity.
2.  **Transaction Record**: Creates a `Transaction` entity with status `PENDING`.
3.  **Gateway Integration**: Generates necessary payment data (QR Code URL, Bank Account details, Reference Code) via SePay.
4.  **Response**: Returns formatted data for the frontend to display a payment UI (e.g., a QR code to scan).

### Webhook Processing
The `ProcessWebhook` operation handles real-time updates from the payment gateway:
1.  **Verification**: Validates the incoming webhook payload from SePay.
2.  **State Transition**: Updates the `Transaction` status to `SUCCESS` or `FAILED`.
3.  **Fulfillment**:
    - **Subscription**: Activates the user's subscription plan.
    - **Series**: Creates a `UserSeriesPurchase` record, granting access to the content.
    - **Donation**: Logs the donation.

## Data Model

### Transaction
The central ledger for payments.
```go
type Transaction struct {
    ID             uuid.UUID
    UserID         uuid.UUID
    Amount         decimal.Decimal
    Currency       string // Default "VND"
    Provider       TransactionProvider // "SEPAY"
    Type           TransactionType // "SUBSCRIPTION", "SERIES", "DONATION"
    Status         TransactionStatus // "PENDING", "SUCCESS", "FAILED"
    TargetID       *uuid.UUID // ID of Series or other item being bought
    PlanID         *string // ID of Subscription Plan
    ReferenceCode  string // Unique code for bank transfer matching
    SePayID        string // External ID from gateway
}
```

### UserSeriesPurchase
Grants permanent access to specific content.
```go
type UserSeriesPurchase struct {
    UserID    uuid.UUID
    SeriesID  uuid.UUID
    Amount    decimal.Decimal
    CreatedAt time.Time
}
```

## API Reference (Internal)
### CreatePaymentUseCase
- `Execute(ctx, req)`: Initialize a payment.
  - Input: `CreatePaymentRequest` (Amount, Type, Gateway, TargetID/PlanID).
  - Output: `CreatePaymentResponse` (QRDataURL, AccountNo, ReferenceCode).

### ProcessWebhookUseCase
- `Execute(ctx, req)`: Process gateway callback.
  - Input: `ProcessWebhookRequest` (SePay payload).
  - Output: Updated `Transaction`.
