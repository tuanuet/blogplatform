# Feature: password reset (forgot/reset)

This document describes the password reset feature for the API. It covers the
endpoints, expected behavior, and high-level implementation impact so the
feature is ready for implementation.

## Objective

Let users recover access to their account by requesting a password reset email
and setting a new password via a time-limited token. The API must not leak
whether an email exists.

## Requirements

- [ ] Add `POST /api/v1/auth/forgot-password`.
- [ ] Always return `200` with a generic message from `forgot-password`.
- [ ] If the user exists and is eligible, generate a reset token stored in Redis
      with a TTL (default: 30 minutes) and send a password reset email.
- [ ] Rate-limit `forgot-password` (reuse existing rate limit middleware).
- [ ] Add `GET /api/v1/auth/reset-password/validate?token=...`.
- [ ] Return JSON `200` when the token is valid and JSON `400` when the token is
      invalid or expired.
- [ ] Add `POST /api/v1/auth/reset-password`.
- [ ] Validate the token and set a new `password_hash` (bcrypt) for the user.
- [ ] Invalidate the reset token after a successful password reset.
- [ ] Use a dedicated token prefix, for example `reset:<uuid>`, when storing the
      token in Redis via `SessionRepository`.
- [ ] Build reset links using `app.public_url`.
- [ ] Enforce password rules consistent with registration (minimum length 8).
- [ ] Do not log raw tokens.

## Non-goals

- Frontend reset-password UI or redirects.
- Multi-factor authentication.
- Account lockout or risk-based reset policies.
- Changing the existing cookie session model.

## Acceptance criteria

1. `POST /api/v1/auth/forgot-password` returns `200` with a generic message for
   all requests, regardless of whether the email exists.
2. When the email exists, the API creates a `reset:` token in Redis with a TTL
   of about 30 minutes and attempts to send a password reset email.
3. `GET /api/v1/auth/reset-password/validate?token=...` returns `200` for a valid
   token and `400` for an invalid or expired token.
4. `POST /api/v1/auth/reset-password` with a valid token updates the user's
   password hash and invalidates the token.
5. After a successful reset, the old password no longer works and the new
   password works.
6. Swagger documentation includes the new endpoints and request/response DTOs.

## Technical context

### Impacted areas

- `internal/application/usecase/auth/usecase.go`
- `internal/interfaces/http/handler/auth/auth_handler.go`
- `internal/interfaces/http/router/auth_routes.go`
- `internal/domain/service/email_service.go`
- `internal/domain/service/email_service_impl.go`
- `internal/infrastructure/email/templates/password_reset.html`
- `internal/application/dto/auth.go`

### Data changes

- No database schema changes.
- Redis keys reuse `SessionRepository` and add a reset prefix in the token
  string (for example, `session:reset:<uuid>`).
