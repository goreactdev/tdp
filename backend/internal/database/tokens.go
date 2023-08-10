package database

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ton-developer-program/internal/validator"
)

const (
	ScopeAuthentication = "authentication"
)

type Token struct {
	Plaintext string `json:"token"`
	Hash      []byte `json:"-"`
	UserID    int64  `json:"-"`
	Expiry    uint64 `json:"expiry"`
	Scope     string `json:"-"`
}


type TokensModel struct {
	DB *sqlx.DB
}


func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {

	token := &Token{
		UserID: userID,
		Expiry: uint64(time.Now().Add(ttl).Unix()),
		Scope:  scope,
	}

	// Increase the random byte size to 32
	randomBytes := make([]byte, 32)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	// Use base64 encoding instead of base32
	token.Plaintext = base64.RawURLEncoding.EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", " token must be provided")
	v.Check(len(tokenPlaintext) == 43, "token must be 43 bytes long")
}

func (m *TokensModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {		
		return nil, err
	}

	err = m.Insert(token)
	return token, err
}

func (m *TokensModel) Insert(token *Token) error {
	query := `
	    INSERT INTO tokens (hash, user_id, expiry, scope)
	    VALUES($1,$2,$3,$4)`

	args := []interface{}{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m *TokensModel) DeleteAllForUser(scope string, userID int64) error {
	query := `
	    DELETE FROM tokens
	    WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, scope, userID)
	return err
}
