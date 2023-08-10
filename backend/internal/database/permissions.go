package database

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)


type Role struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Permissions []Permission  `json:"permissions"`
	Tags        []int64 `json:"tags"`
}

type Permission struct {
	ID        int64  `db:"id" json:"id"`
	Name 	  string `db:"name" json:"name"`
	Route 	  string `db:"route" json:"route"`
	Method 	  string `db:"method" json:"method"`
}


type PermissionModel struct {
	DB *sqlx.DB
}

type Roles []Role

// insert new role
func (m PermissionModel) InsertRole(role Role) (int64, error) {
	query := `INSERT INTO roles (name, description) VALUES ($1, $2) RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int64

	err := m.DB.QueryRowContext(ctx, query, role.Name, role.Description).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// insert role permissions
func (m PermissionModel) InsertRolePermissions(roleID int64, permissions []string) error {
	query := `INSERT INTO roles_permissions (role_id, permission_id) VALUES ($1, (SELECT id FROM permissions WHERE code = $2))`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for _, permission := range permissions {
		_, err := m.DB.ExecContext(ctx, query, roleID, permission)
		if err != nil {
			return err
		}
	}

	return nil
}

// insert role permissions by array of permission ids
func (m PermissionModel) InsertRolePermissionsByIDs(roleID int64, permissionIDs []int64) error {
	query := `INSERT INTO roles_permissions (role_id, permission_id) VALUES ($1, $2)`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for _, permissionID := range permissionIDs {
		_, err := m.DB.ExecContext(ctx, query, roleID, permissionID)
		if err != nil {
			return err
		}
	}

	return nil
}

// insert role to user
func (m PermissionModel) InsertUserRole(userID int64, roleID int64) error {
	query := `INSERT INTO users_roles (user_id, role_id) VALUES ($1, $2)`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, userID, roleID)
	if err != nil {
		return err
	}

	return nil
}

// get all roles
func (m PermissionModel) GetAllRoles() (Roles, error) {
	query := `SELECT id, name, description FROM roles ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var roles Roles

	for rows.Next() {
		var role Role

		err := rows.Scan(&role.ID, &role.Name, &role.Description)
		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

// get role by id
func (m PermissionModel) GetRoleByID(id int64) (*Role, error) {
	query := `SELECT id, name, description FROM roles WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)

	var role Role

	err := row.Scan(&role.ID, &role.Name, &role.Description)

	if err != nil {
		return nil, err
	}

	return &role, nil
}



// get role by name
func (m PermissionModel) GetRoleByName(name string) (*Role, error) {
	query := `SELECT id, name, description FROM roles WHERE name = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, name)

	var role Role

	err := row.Scan(&role.ID, &role.Name, &role.Description)

	if err != nil {
		return nil, err
	}

	return &role, nil
}

// get role permissions
func (m PermissionModel) GetRolePermissions(roleID int64) ([]Permission, error) {
	query := `SELECT permissions.id, permissions.route, permissions.method, permissions.name FROM permissions INNER JOIN roles_permissions ON roles_permissions.permission_id = permissions.id WHERE roles_permissions.role_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var permissions []Permission

	for rows.Next() {
		var permission Permission

		err := rows.Scan(&permission.ID, &permission.Route, &permission.Method, &permission.Name)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

// get user roles
func (m PermissionModel) GetUserRoles(userID int64) (*Role, error) {
	query := `SELECT roles.id, roles.name, roles.description FROM roles INNER JOIN users_roles ON users_roles.role_id = roles.id WHERE users_roles.user_id = $1 LIMIT 1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var role Role

	for rows.Next() {
		err := rows.Scan(&role.ID, &role.Name, &role.Description)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &role, nil

}

// get user permissions by user id using roles
func (m PermissionModel) GetUserPermissions(userID int64) ([]Permission, error) {
	query := `SELECT permissions.code FROM permissions INNER JOIN roles_permissions ON roles_permissions.permission_id = permissions.id INNER JOIN users_roles ON users_roles.role_id = roles_permissions.role_id WHERE users_roles.user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var permissions []Permission

	for rows.Next() {
		var permission Permission

		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}


// delete role
func (m PermissionModel) DeleteRole(id int64) error {
	query := `DELETE FROM roles WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// delete role permissions
func (m PermissionModel) DeleteRolePermissions(roleID int64) error {
	query := `DELETE FROM roles_permissions WHERE role_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, roleID)
	if err != nil {
		return err
	}

	return nil
}

// delete user roles
func (m PermissionModel) DeleteUserRoles(userID int64) error {
	query := `DELETE FROM users_roles WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

// update role
func (m PermissionModel) UpdateRole(role *Role) error {
	query := `UPDATE roles SET name = $1, description = $2 WHERE id = $3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, role.Name, role.Description, role.ID)
	if err != nil {
		return err
	}

	return nil

}

// update role permissions
func (m PermissionModel) UpdateRolePermissions(roleID int64, permissions []string) error {
	query := `INSERT INTO roles_permissions (role_id, permission_id) VALUES ($1, (SELECT id FROM permissions WHERE code = $2))`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for _, permission := range permissions {
		_, err := m.DB.ExecContext(ctx, query, roleID, permission)
		if err != nil {
			return err
		}
	}

	return nil
}

// update user roles
func (m PermissionModel) UpdateUserRoles(userID int64, roles []int64) error {
	query := `INSERT INTO users_roles (user_id, role_id) VALUES ($1, $2)`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for _, role := range roles {
		_, err := m.DB.ExecContext(ctx, query, userID, role)
		if err != nil {
			return err
		}

	}

	return nil
}

// get all permissions
func (m PermissionModel) GetAllPermissions() ([]Permission, error) {
	query := `SELECT * FROM permissions`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}


	defer rows.Close()

	var permissions []Permission

	for rows.Next() {
		var permission Permission

		err := rows.Scan(&permission.ID, &permission.Name, &permission.Route, &permission.Method)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil

	

}







// func (m PermissionModel) GetAllForUser(userID int64) (Permissions, error) {
// 	query := `
// 		SELECT permissions.code
// 		FROM permissions
// 		INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
// 		INNER JOIN users ON users_permissions.user_id = users.id
// 		WHERE users.id = $1`

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	rows, err := m.DB.QueryContext(ctx, query, userID)
// 	if err != nil {		
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var permissions Permissions

// 	for rows.Next() {
// 		var permission string

// 		err := rows.Scan(&permission)
// 		if err != nil {		
// 			return nil, err
// 		}

// 		permissions = append(permissions, permission)
// 	}
// 	if err = rows.Err(); err != nil {		
// 		return nil, err
// 	}

// 	return permissions, nil
// }
