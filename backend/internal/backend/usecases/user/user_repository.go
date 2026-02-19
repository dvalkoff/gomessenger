package user

import (
	"database/sql"

	"github.com/google/uuid"
)

type UserRow struct {
	Id              uuid.UUID
	Nickname        string
	HashedPassword  []byte
	IdentityPubKey  PubKey
	SignedPubPrekey PubKey
}

type UserRepository interface {
	SaveUser(tx *sql.Tx, user UserRow) error
	FindUsersByNickname(string) ([]UserRow, error)
	FindUserById(uuid.UUID) (UserRow, error)
	AddContact(userId, contactId uuid.UUID) error
	GetContacts(userId uuid.UUID) ([]UserRow, error)
	StartTx() (*sql.Tx, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (repository *userRepository) SaveUser(tx *sql.Tx, user UserRow) error {
	sql := `
	INSERT INTO messenger.users(id, nickname, password, identity_pub_key, signed_pub_key)
	VALUES($1, $2, $3, $4, $5);`
	_, err := tx.Exec(sql,
		user.Id,
		user.Nickname,
		user.HashedPassword,
		user.IdentityPubKey,
		user.SignedPubPrekey,
	)
	return err
}

func (repository *userRepository) FindUsersByNickname(nicknameSubstring string) ([]UserRow, error) {
	sql := `
	SELECT id, nickname FROM messenger.users
	WHERE nickname LIKE $1 || '%'`
	rows, err := repository.db.Query(sql, nicknameSubstring)
	if err != nil {
		return nil, err
	}

	userRows := make([]UserRow, 0)
	for rows.Next() {
		userRow := UserRow{}
		rows.Scan(&userRow.Id, &userRow.Nickname)
		userRows = append(userRows, userRow)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return userRows, nil
}

func (repository *userRepository) FindUserById(id uuid.UUID) (UserRow, error) {
	sql := `
	SELECT id, nickname, password FROM messenger.users
	WHERE id = $1`
	user := UserRow{}
	row := repository.db.QueryRow(sql, id)
	err := row.Scan(&user.Id, &user.Nickname, &user.HashedPassword)
	if err != nil {
		return UserRow{}, err
	}
	return user, err
}

func (repository *userRepository) AddContact(userId, contactId uuid.UUID) error {
	sql := `INSERT INTO messenger.friends(user_id, contact_user_id)
	VALUES ($1, $2)`
	_, err := repository.db.Exec(sql, userId, contactId)
	return err
}

func (repository *userRepository) GetContacts(userId uuid.UUID) ([]UserRow, error) {
	sql := `
		SELECT u.id, u.nickname FROM messenger.users u
		JOIN messenger.contacts c ON u.user_id = c.contact_user_id
		WHERE c.user_id = $1;`
	rows, err := repository.db.Query(sql, userId)
	if err != nil {
		return nil, err
	}

	userRows := make([]UserRow, 0)
	for rows.Next() {
		userRow := UserRow{}
		rows.Scan(&userRow.Id, &userRow.Nickname)
		userRows = append(userRows, userRow)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return userRows, nil
}

func (repository *userRepository) StartTx() (*sql.Tx, error) {
	return repository.db.Begin()
}
