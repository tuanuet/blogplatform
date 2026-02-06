# Feature: Profile Description

## Objective

Add a 'description' field to the user profile for detailed text (up to 5000 chars), visible on public profiles.

## Requirements

- [ ] Add `Description` field to `User` entity.
  - Type: `*string`
  - Database: `TEXT` or `VARCHAR(5000)`
  - Validation: Max 5000 characters
- [ ] Update `UpdateProfileRequest` DTO.
  - Add `Description *string` with validation.
- [ ] Update `ProfileResponse` DTO.
  - Add `Description string`.
- [ ] Update `PublicProfileResponse` DTO.
  - Add `Description string`.
- [ ] Update `ProfileUseCase`.
  - Handle `Description` in `UpdateProfile`.
  - Map `Description` in `toProfileResponse`.
  - Map `Description` in `toPublicProfileResponse`.

## Technical Context

- Impacted Services: Profile Service (Monolith)
- New Data/Fields: `users.description`

## Acceptance Criteria

1. User can update profile with a description up to 5000 chars.
2. User can see their description in their profile.
3. Public visitors can see the user's description.
