package user

import "database/sql"

type UserRow struct {
	nickname string
	name string
	hashedPassword []byte
}

type UserRepository interface {
	saveUser(UserRow) error
	findUsersByNickname(string) ([]UserRow, error)
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

func (repository *userRepository) findUsersByNickname(nicknameSubstring string) ([]UserRow, error) {
	sql := `
	SELECT name, nickname FROM messenger.users
	WHERE nickname = $1;`
	rows, err := repository.db.Query(sql, nicknameSubstring)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // TODO: should I call Close()? what should I do with err that Close() returns?
	
	userRows := []UserRow{}
	for rows.Next() {
		userRow := UserRow{}
		rows.Scan(&userRow.name, &userRow.nickname)
		userRows = append(userRows, userRow)
	}
	if err := rows.Err(); err != nil { // TODO: when do I call Err()?
		return nil, err
	}
	return userRows, nil
}
