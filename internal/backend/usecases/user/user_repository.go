package user

import (
	"context"
	"database/sql"
	"time"
)

const (
	txTimeout time.Duration = 1 * time.Second
)

type UserRow struct {
	Nickname       string
	Name           string
	HashedPassword []byte
}

type UserRepository interface {
	SaveUser(context.Context, UserRow) error
	FindUsersByNickname(ctx context.Context, nickname string) ([]UserRow, error)
	AddFriend(ctx context.Context, nickname, friendsNickname string) error
	GetFriends(ctx context.Context, nickname string) ([]UserRow, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (repository *userRepository) SaveUser(ctx context.Context, user UserRow) error {
	ctx, cancel := context.WithTimeout(ctx, txTimeout)
	defer cancel()
	sql := `
	INSERT INTO messenger.users(nickname, name, password)
	VALUES($1, $2, $3);`
	_, err := repository.db.ExecContext(ctx, sql, user.Nickname, user.Name, user.HashedPassword)
	return err
}

func (repository *userRepository) FindUsersByNickname(ctx context.Context, nicknameSubstring string) ([]UserRow, error) {
	ctx, cancel := context.WithTimeout(ctx, txTimeout)
	defer cancel()
	sql := `
	SELECT name, nickname, password FROM messenger.users
	WHERE nickname = $1;`
	rows, err := repository.db.QueryContext(ctx, sql, nicknameSubstring)
	if err != nil {
		return nil, err
	}

	userRows := []UserRow{}
	for rows.Next() {
		userRow := UserRow{}
		rows.Scan(&userRow.Name, &userRow.Nickname, &userRow.HashedPassword)
		userRows = append(userRows, userRow)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return userRows, nil
}

func (repository *userRepository) AddFriend(ctx context.Context, nickname, friendsNickname string) error {
	ctx, cancel := context.WithTimeout(ctx, txTimeout)
	defer cancel()
	sql := `INSERT INTO messenger.friends(nickname, friends_nickname)
	VALUES ($1, $2)`
	_, err := repository.db.ExecContext(ctx, sql, nickname, friendsNickname)
	return err
}

func (repository *userRepository) GetFriends(ctx context.Context, nickname string) ([]UserRow, error) {
	ctx, cancel := context.WithTimeout(ctx, txTimeout)
	defer cancel()
	sql := `
		SELECT u.name, u.nickname FROM messenger.users u
		JOIN messenger.friends f ON u.nickname = f.friends_nickname
		WHERE f.nickname = $1;`
	rows, err := repository.db.QueryContext(ctx, sql, nickname)
	if err != nil {
		return nil, err
	}

	userRows := []UserRow{}
	for rows.Next() {
		userRow := UserRow{}
		rows.Scan(&userRow.Name, &userRow.Nickname)
		userRows = append(userRows, userRow)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return userRows, nil
}
