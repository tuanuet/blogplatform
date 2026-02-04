-- Migration: Add reaction counts to blogs table
-- Description: Adds upvote_count and downvote_count columns to support blog reactions

ALTER TABLE blogs
ADD COLUMN IF NOT EXISTS upvote_count INT NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS downvote_count INT NOT NULL DEFAULT 0;
