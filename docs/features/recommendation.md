# Recommendation Feature

## Overview
The Recommendation feature is the discovery engine of the platform, designed to surface relevant content to users. It supports both passive discovery (popular items, related content) and active personalization (user interest-based feeds).

## Architecture
- **UseCase**: `RecommendationUseCase` orchestrates discovery requests.
- **Domain Service**: `RecommendationService` contains the core algorithms for content scoring, matching, and retrieval.

## Core Logic & Features
### Content Discovery
- **Popular Tags**: Retrieves currently trending topics (`GetPopularTags`) to help users explore hot discussions.
- **Related Content**: Analyzes a specific blog post (`GetRelatedBlogs`) to suggest similar articles, likely based on tag overlap, category, or author.

### Personalization
- **Personalized Feed**: The `GetPersonalizedFeed` operation builds a custom content stream for the user. This is likely driven by:
  - Explicit interests (tags the user follows).
  - Implicit signals (reading history, interactions).
- **Interest Management**: Users can manually curate their feed by selecting preferred topics (`UpdateInterests`).

## Data Model
This feature primarily acts as a filter/sorter for existing content entities.

### Key DTOs
- **BlogListResponse**: Standardized blog post representation including author, category, and engagement metrics.
- **TagResponse**: Tag metadata for interest selection and display.

## API Reference (Internal)
### RecommendationUseCase
- `GetPopularTags(ctx, limit)`: Get trending tags.
- `GetRelatedBlogs(ctx, blogID, limit)`: Get contextually similar posts.
- `GetPersonalizedFeed(ctx, userID, page, pageSize)`: Get the user's "For You" feed.
- `UpdateInterests(ctx, userID, tagIDs)`: Save user topic preferences.
