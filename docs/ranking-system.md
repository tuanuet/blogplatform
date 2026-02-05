# User Ranking System

This document describes the user ranking system that calculates and tracks user rankings based on velocity scores.

## Overview

The ranking system evaluates user activity and growth to assign rank positions. It considers follower growth (60% weight) and blog post activity (40% weight) to compute a composite velocity score.

## How Rankings Work

### Data Collection

The system captures daily snapshots of follower counts and monitors blog publishing activity. These metrics feed into the velocity calculation that determines user rankings.

### Velocity Score Calculation

The composite score formula combines two key metrics:

**Composite Score = (Follower Growth Rate × 0.6) + (Blog Post Velocity × 0.4)**

#### Follower Growth Rate (60%)

Calculated over a 30-day window with these safeguards:

- Minimum follower requirement: 100 followers for normalization
- Growth rate cap: Maximum 10x (1000%) to prevent manipulation
- Special handling for small follower counts and missing historical data

**Calculation Details:**

The growth rate is calculated as a ratio (not percentage):

```
If previous count = 0:
    Growth Rate = current / 100

If previous count < 100:
    Growth Rate = (current - previous) / 100

Otherwise:
    Growth Rate = (current - previous) / previous
    (capped at 10.0)
```

**Example:**

A user had 1000 followers 30 days ago and now has 1200 followers:

```
Growth Rate = (1200 - 1000) / 1000 = 0.20 (20% as ratio)
```

#### Blog Post Velocity (40%)

Measures publishing activity as posts per day over the last 30 days.

**Formula:**

```
Velocity = Number of Posts in 30 Days / 30
```

**Example:**

A user published 45 posts in the last 30 days:

```
Velocity = 45 / 30 = 1.5 posts per day
```

### Ranking Assignment

Users receive rank positions based on their composite scores in descending order:

1. Highest composite score gets rank #1
2. Scores are sorted and positions assigned sequentially
3. Rank positions update during recalculation
4. Previous rankings are archived to history before recalculation

## Architecture

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│   User Activity │────▶│  Daily Snapshot  │────▶│ Velocity Score  │
│  (Followers,    │     │   (Followers)    │     │  Calculation    │
│   Blog Posts)   │     │                  │     │                 │
└─────────────────┘     └──────────────────┘     └────────┬────────┘
                                                          │
┌─────────────────┐     ┌──────────────────┐              │
│   API Response  │◀────│  Rank Assignment │◀─────────────┘
│  (Trending, Top │     │   (Manual/Job)   │
│   User Detail)  │     │                  │
└─────────────────┘     └──────────────────┘
```

## Database Schema

### Core Tables

**user_velocity_scores**
Stores current velocity metrics and rankings.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | Reference to users table |
| follower_count | int | Current follower count |
| follower_growth_rate | float | 30-day growth ratio (e.g., 0.2 for 20%) |
| blog_post_velocity | float | Posts per day |
| composite_score | float | Weighted combination |
| rank_position | int | Current rank (nullable) |
| calculation_date | timestamp | When calculated |
| created_at | timestamp | Record creation time |
| updated_at | timestamp | Last update time |

**user_ranking_history**
Archives historical ranking snapshots.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | Reference to user |
| rank_position | int | Rank at archive time |
| composite_score | float | Score at archive time |
| follower_count | int | Follower count at archive time |
| recorded_at | timestamp | Archive timestamp |

**user_follower_snapshots**
Daily follower count snapshots for growth calculation.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | Reference to user |
| follower_count | int | Count at snapshot time |
| snapshot_date | date | Date of snapshot |

### Indexes

- Composite index on `composite_score DESC` for ranking queries
- Index on `rank_position` for rank-based lookups
- Index on `(user_id, snapshot_date)` for growth calculations

## API Endpoints

All ranking endpoints are available at `/api/v1/rankings`.

### Get Trending Users

`GET /api/v1/rankings/trending`

Returns users ranked by velocity score with pagination.

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | int | 1 | Page number |
| pageSize | int | 20 | Results per page (max 100) |
| category | string | - | Optional category filter |

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "johndoe",
      "displayName": "John Doe",
      "avatarUrl": "https://example.com/avatar.jpg",
      "followerCount": 1500,
      "followerGrowthRate": 0.25,
      "blogPostVelocity": 1.5,
      "compositeScore": 15.6,
      "rank": 1,
      "calculationDate": "2026-02-05T00:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "pageSize": 20,
    "total": 150,
    "totalPages": 8
  }
}
```

### Get Top Users

`GET /api/v1/rankings/top`

Returns top users ranked by total follower count.

**Query Parameters:** Same as trending endpoint.

**Response:** Same format as trending endpoint, but sorted by follower count.

### Get User Ranking Details

`GET /api/v1/rankings/users/{userId}`

Returns detailed ranking information for a specific user including history.

**Response:**

```json
{
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "username": "johndoe",
  "displayName": "John Doe",
  "avatarUrl": "https://example.com/avatar.jpg",
  "followerCount": 1500,
  "followerGrowthRate": 0.25,
  "blogPostVelocity": 1.5,
  "compositeScore": 15.6,
  "rank": 5,
  "previousRank": 7,
  "rankChange": 2,
  "calculationDate": "2026-02-05T00:00:00Z",
  "history": [
    {
      "rankPosition": 7,
      "compositeScore": 14.2,
      "followerCount": 1400,
      "recordedAt": "2026-02-04T00:00:00Z"
    }
  ]
}
```

### Recalculate Scores (Admin Only)

`POST /api/v1/rankings/recalculate`

Triggers manual recalculation of all velocity scores and rankings.

**Authentication:** Requires admin privileges via Bearer token.

**Response:**

```json
{
  "success": true,
  "message": "rankings recalculated successfully"
}
```

## Configuration

### Ranking Weights

The current weights are defined in the ranking service configuration:

```go
RankingConfig{
    FollowerGrowthWeight:   0.6,    // 60%
    BlogPostVelocityWeight: 0.4,    // 40%
    TimeWindowDays:         30,
    MinFollowersForRate:    100,
}
```

**Important:** Ensure `FollowerGrowthWeight + BlogPostVelocityWeight = 1.0`

The maximum growth rate cap (10x or 1000%) is hardcoded in the calculation logic.

To modify these values, update the configuration in `internal/domain/service/ranking_service.go`.

## Admin Operations

### Manual Recalculation

Administrators can trigger manual score recalculation via the API:

```bash
curl -X POST https://api.example.com/api/v1/rankings/recalculate \
  -H "Authorization: Bearer {admin_token}"
```

This is useful when:

- Correcting data anomalies
- Testing algorithm changes
- Recovering from system outages

### Archive Management

Ranking history is automatically archived during recalculation. Old archive records can be purged periodically to manage storage:

```sql
-- Example: Remove archives older than 90 days
DELETE FROM user_ranking_history 
WHERE recorded_at < NOW() - INTERVAL '90 days';
```

### Troubleshooting

**Issue: User has no rank position**

- Verify the user has velocity score data
- Check if scores have been calculated
- Manually trigger recalculation if needed

**Issue: Rankings seem incorrect**

- Review follower snapshot data for gaps
- Check calculation logs for errors
- Verify time window settings

**Issue: Slow ranking queries**

- Ensure database indexes are present
- Consider adding composite indexes for frequent queries
- Monitor query performance in slow query logs

## Integration Points

### Follower Snapshot Creation

Follower snapshots are created by integrating with the user follower system. Ensure your follower tracking creates snapshots when counts change.

### Blog Publishing Impact

Blog posts automatically affect rankings through the velocity calculation. The system tracks publishing dates and calculates velocity based on posts within the time window.

## Security Considerations

- Growth rate capping prevents manipulation by artificially inflating follower counts
- Admin-only endpoints require proper authentication and authorization middleware
- Rate limiting recommended for public ranking endpoints
- Input validation on all query parameters

## Performance Optimization

### Database

- Indexes on composite_score and rank_position enable fast sorting
- Partitioning user_ranking_history by recorded_at improves historical queries
- Regular VACUUM and ANALYZE maintain query performance

### Caching

Consider caching popular ranking endpoints:

- Trending users (TTL: 5 minutes)
- Top 10 users (TTL: 1 minute)
- User ranking details (TTL: 10 minutes)

### Scaling

For high-traffic scenarios:

1. Use read replicas for ranking queries
2. Implement Redis caching layer
3. Consider materialized views for complex calculations
4. Batch process follower snapshots during off-peak hours

## Monitoring

Track these metrics for ranking system health:

- Score calculation execution time and success rate
- Average composite score distribution
- API endpoint response times
- Database query performance
- Cache hit rates (if implemented)

## Next Steps

- Review the API documentation in `docs/swagger.yaml` for complete endpoint specifications
- Check `internal/domain/service/ranking_service.go` for algorithm implementation details
- See migration files for database schema evolution
- Read `ranking-algorithm.md` for detailed algorithm documentation

For questions or issues, refer to the codebase or contact the development team.
