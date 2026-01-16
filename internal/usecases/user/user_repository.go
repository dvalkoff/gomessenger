package user

import "database/sql"

type UserRow struct {
	nickname string
	name string
	hashedPassword []byte
}

type UserRepository interface {
	saveUser(UserRow) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (repository *userRepository) saveUser(user UserRow) error {
	sql := `
	INSERT INTO messenger.users(nickname, name, password)
	VALUES($1, $2, $3);`
	_, err := repository.db.Exec(sql, user.nickname, user.name, user.hashedPassword)
	return err
}
