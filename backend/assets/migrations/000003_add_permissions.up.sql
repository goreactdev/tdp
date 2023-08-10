CREATE TABLE IF NOT EXISTS permissions (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    route text NOT NULL,
    method text NOT NULL
);

-- Add the two permissions to the table.
INSERT INTO permissions (name, route, method)
VALUES
('permissions:users-read', '/v1/admin/users', 'GET'),
('permissions:users-edit', '/v1/admin/users/:id', 'PATCH'),
('permissions:users-create', '/v1/admin/users', 'POST'),
('permissions:users-delete', '/v1/admin/users/:id', 'DELETE'),
('permissions:roles-read', '/v1/admin/roles', 'GET'),
('permissions:roles-edit', '/v1/admin/roles/:id', 'PATCH'),
('permissions:roles-create', '/v1/admin/roles', 'POST'),
('permissions:roles-delete', '/v1/admin/roles/:id', 'DELETE'),
('permissions:collections-read', '/v1/admin/collections', 'GET'),
('permissions:collections-edit', '/v1/admin/collections/:id', 'PATCH'),
('permissions:collections-create', '/v1/admin/collections', 'POST'),
('permissions:collections-delete', '/v1/admin/collections/:id', 'DELETE'),
('permissions:minted-nfts-read', '/v1/admin/minted-nfts', 'GET'),
('permissions:minted-nfts-edit', '/v1/admin/minted-nfts/:id', 'PATCH'),
('permissions:minted-nfts-create', '/v1/admin/minted-nfts', 'POST'),
('permissions:minted-nfts-delete', '/v1/admin/minted-nfts/:id', 'DELETE'),
('permissions:activities-read', '/v1/admin/activities', 'GET'),
('permissions:activities-edit', '/v1/admin/activities/:id', 'PATCH'),
('permissions:activities-create', '/v1/admin/activities', 'POST'),
('permissions:activities-delete', '/v1/admin/activities/:id', 'DELETE'),
('permissions:prototype-nfts-read', '/v1/admin/prototype-nfts', 'GET'),
('permissions:prototype-nfts-edit', '/v1/admin/prototype-nfts/:id', 'PATCH'),
('permissions:prototype-nfts-create', '/v1/admin/prototype-nfts', 'POST'),
('permissions:prototype-nfts-delete', '/v1/admin/prototype-nfts/:id', 'DELETE'),
('permissions:rewards-read', '/v1/admin/rewards', 'GET'),
('permissions:rewards-edit', '/v1/admin/rewards/:id', 'PATCH'),
('permissions:rewards-create', '/v1/admin/rewards', 'POST'),
('permissions:rewards-delete', '/v1/admin/rewards/:id', 'DELETE'),
('permissions:permissions-read', '/v1/admin/permissions', 'GET'),
('permissions:permissions-edit', '/v1/admin/permissions/:id', 'PATCH'),
('permissions:permissions-create', '/v1/admin/permissions', 'POST'),
('permissions:permissions-delete', '/v1/admin/permissions/:id', 'DELETE'),
('permissions:existing-collection-read', '/v1/admin/existing-collection', 'POST'),
('permissions:csv-upload-create', '/v1/admin/csv/upload', 'POST'),
('permissions:merch-create', '/v1/admin/merch', 'POST');



