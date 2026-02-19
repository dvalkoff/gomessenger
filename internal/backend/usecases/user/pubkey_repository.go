package user

import (
	"database/sql"

	"github.com/google/uuid"
)

type PubKey []byte

type KeyType int

const (
	OneTimePrekey = iota
	SignedPrekey
	IdentityKey
)

type PubKeyRow struct {
	Id     int
	Key    PubKey
	Type   KeyType
	UserId uuid.UUID
}

type PubKeyRepository interface {
	SaveKeys(*sql.Tx, []PubKeyRow) error
}

type pubKeyRepository struct {
	db *sql.DB
}

func NewPubKeyRepository(db *sql.DB) PubKeyRepository {
	return &pubKeyRepository{db: db}
}

func (repository *pubKeyRepository) SaveKeys(tx *sql.Tx, keys []PubKeyRow) error {
	sql := `
	INSERT INTO messenger.users_pub_keys(pub_key, key_type, user_id)
	VALUES($1, $2, $3)`
	statement, err := tx.Prepare(sql)
	if err != nil {
		return err
	}
	defer statement.Close()
	for _, key := range keys {
		_, err = statement.Exec(key.Key, key.Type, key.UserId)
		if err != nil {
			return err
		}
	}
	return nil
}
