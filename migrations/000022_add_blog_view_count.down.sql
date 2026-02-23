-- Migration: Remove view count from blogs table
-- Description: Removes view_count column and top-viewed query index

DROP INDEX IF EXISTS idx_blogs_top_viewed;

ALTER TABLE blogs
DROP CONSTRAINT IF EXISTS chk_blogs_view_count_non_negative;

ALTER TABLE blogs
DROP COLUMN IF EXISTS view_count;
