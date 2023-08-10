package database

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)


const (
	ProviderGithub = "github"
	ProviderTelegram = "telegram"
)

var (
	ErrDuplicateUsername = errors.New("duplicate username")
	ErrRecordNotFound	= errors.New("record not found")
)

var AnonymousUser = &User{}


func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type GithubUser struct {
    Login string `json:"login"`
    ID int `json:"id"`
    NodeID string `json:"node_id"`
    AvatarURL string `json:"avatar_url"`
    GravatarID string `json:"gravatar_id"`
    URL string `json:"url"`
    HTMLURL string `json:"html_url"`
    FollowersURL string `json:"followers_url"`
    FollowingURL string `json:"following_url"`
    GistsURL string `json:"gists_url"`
    StarredURL string `json:"starred_url"`
    SubscriptionsURL string `json:"subscriptions_url"`
    OrganizationsURL string `json:"organizations_url"`
    ReposURL string `json:"repos_url"`
    EventsURL string `json:"events_url"`
    ReceivedEventsURL string `json:"received_events_url"`
    Type string `json:"type"`
    SiteAdmin bool `json:"site_admin"`
    Name interface{} `json:"name"`
    Company interface{} `json:"company"`
    Blog string `json:"blog"`
    Location interface{} `json:"location"`
    Email interface{} `json:"email"`
    Hireable interface{} `json:"hireable"`
    Bio interface{} `json:"bio"`
    TwitterUsername interface{} `json:"twitter_username"`
    PublicRepos int `json:"public_repos"`
    PublicGists int `json:"public_gists"`
    Followers int `json:"followers"`
    Following int `json:"following"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}


type LinkedAccount struct {
	ID         int64  `db:"id" json:"id"`
	UserID     int64  `db:"user_id" json:"user_id"`
	TelegramUserID *int64  `db:"telegram_user_id" json:"telegram_user_id"`
	Provider   string `db:"provider" json:"provider"`
	AvatarURL  string `db:"avatar_url" json:"avatar_url"`
	Login      string `db:"login" json:"login"`
	AccessToken string `db:"access_token" json:"-"`
	CreatedAt  uint64 `db:"created_at" json:"created_at"`
	UpdatedAt  uint64 `db:"updated_at" json:"updated_at"`
	Version    int    `db:"version" json:"version"`
}

// array of linked accounts
type LinkedAccounts []*LinkedAccount

type User struct {
	ID             int64       `db:"id" json:"id"`
	FirstName      string      `db:"first_name" json:"first_name"`
	LastName       string      `db:"last_name" json:"last_name"`
	Rating 	   float64     `db:"rating" json:"rating"`
	Username       string      `db:"username" json:"username"`
	FriendlyAddress string      `db:"friendly_address" json:"friendly_address"`
	RawAddress        string      `db:"raw_address" json:"raw_address"`
	Job            *string      `db:"job" json:"job"`
	Bio            *string      `db:"bio" json:"bio"`
	Languages      pq.StringArray `db:"languages" json:"languages"`
	Certifications pq.StringArray `db:"certifications" json:"certifications"`
	CreatedAt      uint64   `db:"created_at" json:"created_at"`
	UpdatedAt      uint64   `db:"updated_at" json:"updated_at"`
	AvatarURL      *string     `db:"avatar_url" json:"avatar_url"`
	AwardsCount    int         `db:"awards_count" json:"awards_count"`
	MessagesCount  int64         `db:"messages_count" json:"messages_count,omitempty"`
	LastAwardAt    *uint64     `db:"last_award_at" json:"last_award_at"`
	LinkedAccounts LinkedAccounts `db:"linked_accounts" json:"linked_accounts"`
	Permissions     []Permission `json:"permissions,omitempty"`
	Role 		  *Role       `json:"role,omitempty"`	
	Version        int       `db:"version" json:"version"`
}


type UserModel struct {
	DB *sqlx.DB
}



// get for token

func (m *UserModel) GetForToken(tokenScope string, tokenString string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenString))

	query := ` 
		SELECT 
			users.id,
			users.first_name,
			users.last_name,
			users.username,
			users.friendly_address,
			users.raw_address,
			users.job,
			users.bio,
			users.languages,
			users.certifications,
			users.created_at,
			users.updated_at,
			users.avatar_url,
			users.awards_count,
			users.messages_count,
			users.last_award_at,
			users.version
	FROM users
		INNER JOIN tokens
		ON users.id = tokens.user_id
		WHERE tokens.hash = $1
		AND tokens.scope = $2
		AND tokens.expiry > $3`
	
	args := []interface{}{tokenHash[:], tokenScope, time.Now().Unix()}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()


	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.FriendlyAddress,
		&user.RawAddress,
		&user.Job,
		&user.Bio,
		&user.Languages,
		&user.Certifications,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.AvatarURL,
		&user.AwardsCount,
		&user.MessagesCount,
		&user.LastAwardAt,
		&user.Version,
	)


	if err != nil {
		return nil, err
	}

	return &user, nil



}

func (m *UserModel) Insert(user *User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// generate random string with faker

	query := `
		INSERT INTO users (first_name,last_name, username, raw_address, friendly_address, rating, created_at, updated_at, version)
		VALUES ($1, $2, generate_unique_username(), $3, $4, 0, $5, $6, $7)
		RETURNING *
		`

	row := m.DB.QueryRowContext(ctx, query, gofakeit.FirstName(), gofakeit.LastName(), user.RawAddress, user.FriendlyAddress, time.Now().Unix(), time.Now().Unix(), 1)

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.RawAddress,
		&user.FriendlyAddress,
		&user.Job,
		&user.Bio,
		&user.Languages,
		&user.Certifications,
		&user.AvatarURL,
		&user.AwardsCount,
		&user.MessagesCount,
		&user.Rating,
		&user.LastAwardAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Version,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

//  take the rating from the user

func (m *UserModel) TakeRating(userID int64, rating int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		UPDATE users
		SET rating = rating - $1
		WHERE id = $2
		`

	_, err := m.DB.ExecContext(ctx, query, 1, userID)
	if err != nil {
		return err
	}

	return nil
}

// get top users by rating

func (m *UserModel) GetTopUsers(pagination *Pagination) ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		SELECT	
			users.id,
			users.first_name,
			users.last_name,
			users.username,
			users.raw_address,
			users.friendly_address,
			users.job,
			users.bio,
			users.languages,
			users.certifications,
			users.avatar_url,
			users.awards_count,
			users.messages_count,
			users.rating,
			users.last_award_at,
			users.created_at,
			users.updated_at,
			users.version
		FROM users
		ORDER BY rating DESC
		LIMIT $1
		OFFSET $2
		`

	rows, err := m.DB.QueryContext(ctx, query, pagination.End-pagination.Start, pagination.Start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []*User{}

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Username,
			&user.RawAddress,
			&user.FriendlyAddress,
			&user.Job,
			&user.Bio,
			&user.Languages,
			&user.Certifications,
			&user.AvatarURL,
			&user.AwardsCount,
			&user.MessagesCount,
			&user.Rating,
			&user.LastAwardAt,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Version,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// get user position based on rating and all user counts
func (m *UserModel) GetUserPosition(userID int64) (int64, int64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
    defer cancel()

    var position, count int64

    // Get user rating
    var rating float64
    err := m.DB.QueryRowContext(ctx, `SELECT rating FROM users WHERE id = $1`, userID).Scan(&rating)
    if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, nil
		}
        return 0, 0, err
    }

    // Count users with a higher rating
    err = m.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE rating > $1`, rating).Scan(&position)
    if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, nil
		}
        return 0, 0, err
    }
    position++ // Increase by 1 to get the position of current user

    // Count total users
    err = m.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&count)
    if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, nil
		}
        return 0, 0, err
    }

    return position, count, nil
}


func (m *UserModel) Update(user *User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		UPDATE users SET
			first_name = $1,
			last_name = $2,
			username = $3,
			raw_address = $4,
			friendly_address = $5,
			job = $6,
			bio = $7,
			languages = $8,
			certifications = $9,
			avatar_url = $10,
			awards_count = $11,
			messages_count = $12,
			last_award_at = $13,
			updated_at = $14,	
			version = version + 1
		WHERE id = $15 AND version = $16
		RETURNING *
		`

	row := m.DB.QueryRowContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Username,
		user.RawAddress,
		user.FriendlyAddress,
		user.Job,
		user.Bio,
		pq.Array(user.Languages),
		pq.Array(user.Certifications),
		user.AvatarURL,
		user.AwardsCount,
		user.MessagesCount,
		user.LastAwardAt,
		time.Now().Unix(),
		user.ID,
		user.Version,
	)

	var updatedUser User
	err := row.Scan(
		&updatedUser.ID,
		&updatedUser.FirstName,
		&updatedUser.LastName,
		&updatedUser.Username,
		&updatedUser.RawAddress,
		&updatedUser.FriendlyAddress,
		&updatedUser.Job,
		&updatedUser.Bio,
		&updatedUser.Languages,
		&updatedUser.Certifications,
		&updatedUser.AvatarURL,
		&updatedUser.AwardsCount,
		&updatedUser.MessagesCount,
		&updatedUser.Rating,
		&updatedUser.LastAwardAt,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
		&updatedUser.Version,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("concurrent update conflict: user ID %d with version %d", user.ID, user.Version)
		}
		return nil, err
	}

	return &updatedUser, nil
}


// delete user by id
func (m *UserModel) Delete(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		DELETE FROM users
		WHERE id = $1
		`

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}

// create user
func (m *UserModel) Create(user *User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		INSERT INTO users (first_name,last_name, username, raw_address, friendly_address, created_at, updated_at, version)	
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING *
		`

	row := m.DB.QueryRowContext(ctx, query, user.FirstName, user.LastName, user.Username, user.RawAddress, user.FriendlyAddress, time.Now().Unix(), time.Now().Unix(), 1)

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.RawAddress,
		&user.FriendlyAddress,
		&user.Job,
		&user.Bio,
		&user.Languages,
		&user.Certifications,
		&user.AvatarURL,
		&user.AwardsCount,
		&user.MessagesCount,
		&user.LastAwardAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Version,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

//

// get users f
func (m *UserModel) GetMany(pagination *Pagination, filter, sort string) ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var users []*User

	query := fmt.Sprintf(`
		SELECT *
		FROM users
		%s
		ORDER BY %s
		LIMIT $1 OFFSET $2
		`, filter, sort)

	rows, err := m.DB.QueryContext(ctx, query, pagination.End - pagination.Start, pagination.Start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Username,
			&user.RawAddress,
			&user.FriendlyAddress,
			&user.Job,
			&user.Bio,
			&user.Languages,
			&user.Certifications,
			&user.AvatarURL,
			&user.AwardsCount,
			&user.MessagesCount,
			&user.Rating,
			&user.LastAwardAt,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Version,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	


	return users, nil
}


type UserPermissions struct {
	UserID int64 `db:"user_id" json:"user_id"`
	PermissionID int64 `db:"permission_id" json:"permission_id"`
}


// get all permissions for user
func (m *UserModel) GetAllPermissions(id int64) ([]Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var permissions []Permission

	query := `
	SELECT p.id, p.name, p.route, p.method
	FROM users_roles ur
	INNER JOIN roles_permissions rp ON rp.role_id = ur.role_id
	INNER JOIN permissions p ON p.id = rp.permission_id
	WHERE ur.user_id = $1
			`

	rows, err := m.DB.QueryContext(ctx, query, id)

	if err != nil {
		return nil, err
	}


	defer rows.Close()

	for rows.Next() {
		var permission Permission

		err := rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.Route,
			&permission.Method,
		)

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

	
// get user by address
func (m *UserModel) GetByAddress(address string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)

	defer cancel()

	var user User

	query := `SELECT * FROM users WHERE raw_address = $1`

	err := m.DB.GetContext(ctx, &user, query, address)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}


	return &user, err
}

func (m *UserModel) GetByFriendlyAddress(address string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)

	defer cancel()

	var user User

	query := `SELECT * FROM users WHERE friendly_address = $1`

	err := m.DB.GetContext(ctx, &user, query, address)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}


	return &user, err
}



// Get user by username
func (m *UserModel) GetByUsername(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)

	defer cancel()

	var user User

	query := `SELECT * FROM users WHERE username = $1::text`

	err := m.DB.GetContext(ctx, &user, query, username)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) GetByTelegramUsername(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)

	defer cancel()

	var user User

	query := `SELECT * FROM users WHERE id = (SELECT user_id FROM linked_accounts WHERE login = $1 AND provider = 'telegram') `

	err := m.DB.GetContext(ctx, &user, query, username)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		
		return nil, err
	}
	return &user, nil
}


// delete linked account
func (m *UserModel) DeleteLinkedAccountByUserId(userId int64, provider string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `DELETE FROM linked_accounts WHERE user_id = $1 AND provider = $2`

	_, err := m.DB.ExecContext(ctx, query, userId, provider)

	if err != nil {
		return err
	}

	return nil
}

// get linked accounts
func (m *UserModel) GetLinkedAccounts(userId int64) ([]*LinkedAccount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var accounts []*LinkedAccount

	query := `
		SELECT * FROM linked_accounts WHERE user_id = $1
		`

	rows, err := m.DB.QueryContext(ctx, query, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var account LinkedAccount

		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.TelegramUserID,
			&account.Provider,
			&account.AvatarURL,
			&account.Login,
			&account.AccessToken,
			&account.CreatedAt,
			&account.UpdatedAt,
			&account.Version,
		)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, &account)
	}

	return accounts, nil	
}

// get user by telegram user id
func (m *UserModel) GetByTelegramUserId(telegramUserId int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	
	defer cancel()

	var user User

	query := `SELECT * FROM users WHERE id = (SELECT user_id FROM linked_accounts WHERE telegram_user_id = $1)`

	err := m.DB.GetContext(ctx, &user, query, telegramUserId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &user, err
}

// check if user has 2 linked accounts
func (m *UserModel) HasTwoLinkedAccounts(userId int64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var count int

	query := `
		SELECT COUNT(*) FROM linked_accounts WHERE user_id = $1
		`

	err := m.DB.GetContext(ctx, &count, query, userId)

	if err != nil {
		return false, err
	}

	if count == 2 {
		return true, nil
	}

	return false, nil
}



// insert linked account
func (m *UserModel) InsertLinkedAccount(account *LinkedAccount) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		INSERT INTO linked_accounts (user_id, telegram_user_id, provider, avatar_url, login, access_token, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`

	_, err := m.DB.ExecContext(ctx, query,
		account.UserID,
		account.TelegramUserID,
		account.Provider,
		account.AvatarURL,
		account.Login,
		account.AccessToken,
		account.CreatedAt,
		account.UpdatedAt,
		account.Version,
	)

	if err != nil {
		// if already exists, do nothing
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil
		}

		return err
	}

	return nil
}





func (m *UserModel) GetById(id int64) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var user User

	query := `SELECT * FROM users WHERE id = $1`

	err := m.DB.GetContext(ctx, &user, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &user, err
}
