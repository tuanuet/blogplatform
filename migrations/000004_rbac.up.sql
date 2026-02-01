-- RBAC with Bitmask Permissions Migration

-- Roles table
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Role permissions table (resource-based permissions with bitmask)
-- Permission values: READ=1, CREATE=2, UPDATE=4, DELETE=8
CREATE TABLE IF NOT EXISTS role_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    resource VARCHAR(50) NOT NULL,  -- e.g., 'blogs', 'categories', 'users'
    permissions INTEGER NOT NULL DEFAULT 0,  -- bitmask value
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(role_id, resource)
);

-- User roles junction table
CREATE TABLE IF NOT EXISTS user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, role_id)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_resource ON role_permissions(resource);
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);

-- Insert default roles
INSERT INTO roles (name, description) VALUES
    ('admin', 'Full access to all resources'),
    ('editor', 'Can read, create, and update content'),
    ('contributor', 'Can read and create content'),
    ('viewer', 'Read-only access')
ON CONFLICT (name) DO NOTHING;

-- Insert default permissions for admin role (full access = 15)
INSERT INTO role_permissions (role_id, resource, permissions)
SELECT id, resource, 15
FROM roles, (VALUES ('blogs'), ('categories'), ('tags'), ('comments'), ('users')) AS resources(resource)
WHERE roles.name = 'admin'
ON CONFLICT (role_id, resource) DO NOTHING;

-- Insert default permissions for editor role (READ + CREATE + UPDATE = 7)
INSERT INTO role_permissions (role_id, resource, permissions)
SELECT id, resource, 7
FROM roles, (VALUES ('blogs'), ('categories'), ('tags'), ('comments')) AS resources(resource)
WHERE roles.name = 'editor'
ON CONFLICT (role_id, resource) DO NOTHING;

-- Insert default permissions for contributor role (READ + CREATE = 3)
INSERT INTO role_permissions (role_id, resource, permissions)
SELECT id, resource, 3
FROM roles, (VALUES ('blogs'), ('comments')) AS resources(resource)
WHERE roles.name = 'contributor'
ON CONFLICT (role_id, resource) DO NOTHING;

-- Insert default permissions for viewer role (READ = 1)
INSERT INTO role_permissions (role_id, resource, permissions)
SELECT id, resource, 1
FROM roles, (VALUES ('blogs'), ('categories'), ('tags'), ('comments')) AS resources(resource)
WHERE roles.name = 'viewer'
ON CONFLICT (role_id, resource) DO NOTHING;

-- Comments
COMMENT ON TABLE roles IS 'User roles for RBAC';
COMMENT ON TABLE role_permissions IS 'Resource permissions per role using bitmask (READ=1, CREATE=2, UPDATE=4, DELETE=8)';
COMMENT ON TABLE user_roles IS 'Many-to-many relationship between users and roles';
COMMENT ON COLUMN role_permissions.permissions IS 'Bitmask: READ=1, CREATE=2, UPDATE=4, DELETE=8';
