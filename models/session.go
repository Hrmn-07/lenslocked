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
	// insert a new session into the database, but if there is a conflict with the user_id
	// meaning that user already has a session, we instead update the user’s existing session
	// and set its token_hash to the new value.
	row := ss.DB.QueryRow(`
	INSERT INTO sessions (user.id, token_hash)
	VALUES ($1, $2) ON CONFLICT (user_id) DO
	UPDATE
	SET token_hash = $2
	RETURNING id;`, session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)
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
	SELECT users.id,
    users.email,
    users.password_hash
    FROM sessions
    JOIN users ON users.id = sessions.user_id
    WHERE sessions.token_hash = $1;`, tokenHash)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	// 3. Return the user
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
