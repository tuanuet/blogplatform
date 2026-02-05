-- Migration: Create blog_versions and blog_version_tags tables
-- Description: Creates tables for blog versioning functionality

-- =============================================
-- Table: blog_versions
-- =============================================
CREATE TABLE IF NOT EXISTS blog_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    blog_id UUID NOT NULL REFERENCES blogs(id) ON DELETE CASCADE,
    version_number INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    excerpt TEXT,
    content TEXT NOT NULL,
    thumbnail_url VARCHAR(500),
    status blog_status NOT NULL,
    visibility blog_visibility NOT NULL,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    editor_id UUID NOT NULL REFERENCES users(id),
    change_summary TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_blog_versions_blog_id_version_number UNIQUE (blog_id, version_number)
);

-- Indexes for blog_versions
CREATE INDEX IF NOT EXISTS idx_blog_versions_blog_id ON blog_versions(blog_id);
CREATE INDEX IF NOT EXISTS idx_blog_versions_blog_id_created_at ON blog_versions(blog_id, created_at DESC);

-- =============================================
-- Table: blog_version_tags (many-to-many)
-- =============================================
CREATE TABLE IF NOT EXISTS blog_version_tags (
    version_id UUID NOT NULL REFERENCES blog_versions(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (version_id, tag_id)
);
