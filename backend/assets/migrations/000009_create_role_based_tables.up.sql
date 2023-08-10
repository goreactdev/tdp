-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
    id bigserial PRIMARY KEY,
    name text NOT NULL UNIQUE,
    description text
);

-- Create roles_permissions table
CREATE TABLE IF NOT EXISTS roles_permissions (
    role_id bigint NOT NULL REFERENCES roles ON DELETE CASCADE,
    permission_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

-- Create users_roles table
CREATE TABLE IF NOT EXISTS users_roles (
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    role_id bigint NOT NULL REFERENCES roles ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

CREATE UNIQUE INDEX idx_users_roles_unique ON users_roles (user_id);


-- add default admin role
INSERT INTO roles (name, description)
VALUES
('admin', 'Admin role');

-- add all permissions to admin role
INSERT INTO roles_permissions (role_id, permission_id)
SELECT roles.id, permissions.id
FROM roles, permissions
WHERE roles.name = 'admin';