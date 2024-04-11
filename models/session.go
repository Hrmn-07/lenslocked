package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"lenslocked/rand"
)

type Session struct {
	ID     int
	UserID int
	// Token is only set when creating a new session. When looking up a session
	// this will be left empty, as we only store the hash of a session token
	// in our database and we cannot reverse it into a raw token.
	Token     string
	TokenHash string
}

type SessionService struct {
	DB            *sql.DB
	BytesPerToken int
}

// initialize our token size as constant
const (
	MinBytesPerToken = 32
)

// Create will create a new session for the user provided. The session token
// will be returned as the Token field on the Session type, but only the hashed
// session token is stored in the database.
func (ss *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}

	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: ss.Hash(token),
	}
	row := ss.DB.QueryRow(`
	UPDATE sessions
	SET token_hash = $2
	WHERE user_id = $1
	RETURNING id;`, session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)
	if err == sql.ErrNoRows {
		// if no rows exist, we need to create a new session
		// object for that user
		row = ss.DB.QueryRow(`
		INSERT INTO sessions (user_id, token_hash)
		VALUES ($1, $2)
		RETURNING id;`, session.UserID, session.TokenHash)
		// the error will be overwritten with either a new error or nil
		err = row.Scan(&session.ID)
	}
	// If the err was not sql.ErrNoRows, we need to check if it was any
	// other error. If it was sql.ErrNoRows it will be overwritten inside the if
	// block, and we still need to check for any errors
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &session, nil
}

func (ss *SessionService) Hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	// encode the resulting hash into a string
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

func (ss *SessionService) User(token string) (*User, error) {
	var user User
	// 1. Hash the session token
	tokenHash := ss.Hash(token)
	// 2. Query for the session with that hash
	row := ss.DB.QueryRow(`
	SELECT user_id
	FROM sessions
	WHERE token_hash = $1;`, tokenHash)
	err := row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	// 3. Using the UserID from the session, we need to query for that user
	row = ss.DB.QueryRow(`
	SELECT email, password_hash 
	FROM users WHERE id = $1;`, user.ID)
	err = row.Scan(&user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	// 4. Return the user
	return &user, nil
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.Hash(token)
	_, err := ss.DB.Exec(`
	DELETE FROM sessions
	WHERE token_hash = $1;`, tokenHash)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}
