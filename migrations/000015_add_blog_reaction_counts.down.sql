-- Migration: Remove reaction counts from blogs table
-- Description: Removes upvote_count and downvote_count columns

ALTER TABLE blogs
DROP COLUMN IF EXISTS upvote_count,
DROP COLUMN IF EXISTS downvote_count;
