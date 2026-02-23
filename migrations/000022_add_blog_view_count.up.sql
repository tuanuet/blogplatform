-- Migration: Add view count to blogs table
-- Description: Adds view_count column and index support for top-viewed query

ALTER TABLE blogs
ADD COLUMN IF NOT EXISTS view_count INT NOT NULL DEFAULT 0;

ALTER TABLE blogs
ADD CONSTRAINT chk_blogs_view_count_non_negative CHECK (view_count >= 0);

CREATE INDEX IF NOT EXISTS idx_blogs_top_viewed ON blogs (status, visibility, view_count DESC, published_at DESC)
WHERE deleted_at IS NULL;
