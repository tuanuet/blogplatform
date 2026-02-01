CREATE TABLE IF NOT EXISTS series (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_series_slug ON series(slug) WHERE deleted_at IS NULL;
CREATE INDEX idx_series_author_id ON series(author_id);

CREATE TABLE IF NOT EXISTS series_blogs (
    series_id UUID NOT NULL REFERENCES series(id) ON DELETE CASCADE,
    blog_id UUID NOT NULL REFERENCES blogs(id) ON DELETE CASCADE,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (series_id, blog_id)
);

CREATE INDEX idx_series_blogs_blog_id ON series_blogs(blog_id);

DO $$
DECLARE
    admin_role_id UUID;
BEGIN
    SELECT id INTO admin_role_id FROM roles WHERE name = 'admin';
    
    IF admin_role_id IS NOT NULL THEN
        INSERT INTO role_permissions (id, role_id, resource, permissions, created_at, updated_at)
        VALUES (gen_random_uuid(), admin_role_id, 'series', 15, NOW(), NOW());
    END IF;

    SELECT id INTO admin_role_id FROM roles WHERE name = 'user';
    
    IF admin_role_id IS NOT NULL THEN
        INSERT INTO role_permissions (id, role_id, resource, permissions, created_at, updated_at)
        VALUES (gen_random_uuid(), admin_role_id, 'series', 15, NOW(), NOW());
    END IF;
END $$;

