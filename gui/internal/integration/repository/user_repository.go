package repository

import (
	"database/sql"

	"github.com/google/uuid"
)

type UserRow struct {
	Id           uuid.UUID
	WorkspaceId  int
	Nickname     string
	AccessToken  string
	RefreshToken string
	Keys         []KeyRow
}

func (row *UserRow) GetKeys(keyType KeyType) []KeyRow {
	result := make([]KeyRow, 0)
	for _, key := range row.Keys {
		if key.Type == keyType {
			result = append(result, key)
		}
	}
	return result
}

type KeyType int

const (
	OneTimePrekey = iota
	SignedPrekey
	IdentityKey
)

type KeyRow struct {
	UserId     uuid.UUID
	Type       KeyType
	PrivateKey []byte
}

type UserRepository interface {
	GetCurrentUser(spaceId int) (*UserRow, error)
	SaveUser(UserRow) error
	SetUserCurrent(id uuid.UUID) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetCurrentUser(spaceId int) (*UserRow, error) {
	sql := `
	SELECT id, workspace_id, nickname, access_token, refresh_token
	FROM auth
	WHERE workspace_id = $1 and is_current
	ORDER BY id
	`
	rows, err := r.db.Query(sql, spaceId)
	if err != nil {
		return nil, err
	}
	userRows := make([]UserRow, 0)
	for rows.Next() {
		row := UserRow{}
		rows.Scan(&row.Id, &row.WorkspaceId, &row.Nickname, &row.AccessToken, &row.RefreshToken)
		userRows = append(userRows, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(userRows) == 0 {
		return nil, nil
	}
	user := userRows[0]
	keys, err := r.GetKeys(user.Id)
	if err != nil {
		return nil, err
	}
	user.Keys = keys
	return &user, nil
}

func (r *userRepository) GetKeys(userId uuid.UUID) ([]KeyRow, error) {
	sql := `
	SELECT user_id, key_type, private_key
	FROM keys
	WHERE user_id = $1
	`
	rows, err := r.db.Query(sql, userId)
	if err != nil {
		return nil, err
	}
	keys := make([]KeyRow, 0)
	for rows.Next() {
		key := KeyRow{}
		rows.Scan(&key.UserId, &key.Type, &key.PrivateKey)
		keys = append(keys, key)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return keys, nil
}

func (r *userRepository) SaveUser(row UserRow) error {
	sql := `
	INSERT INTO auth(id, workspace_id, nickname, is_current)
	VALUES($1, $2, $3, $4)
	`
	_, err := r.db.Exec(sql, row.Id, row.WorkspaceId, row.Nickname, false)
	if err != nil {
		return err
	}
	return r.SaveKeys(row.Keys)
}

func (r *userRepository) SaveKeys(keys []KeyRow) error {
	sql := `
	INSERT INTO keys(user_id, key_type, private_key)
	VALUES ($1, $2, $3)
	`
	stmt, err := r.db.Prepare(sql)
	if err != nil {
		return err
	}
	for _, key := range keys {
		_, err = stmt.Exec(key.UserId, key.Type, key.PrivateKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *userRepository) SetUserCurrent(id uuid.UUID) error {
	sqlResetCurrent := `
	UPDATE auth
	SET is_current = false
	`
	_, err := r.db.Exec(sqlResetCurrent)
	if err != nil {
		return err
	}
	sqlSetCurrent := `
	UPDATE auth
	SET is_current = true
	WHERE id = $1
	`
	_, err = r.db.Exec(sqlSetCurrent, id)
	return err
}
