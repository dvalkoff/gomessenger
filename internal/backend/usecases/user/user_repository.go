package user

import "database/sql"

type UserRow struct {
	Nickname string
	Name string
	HashedPassword []byte
}

type UserRepository interface {
	SaveUser(UserRow) error
	FindUsersByNickname(string) ([]UserRow, error)
	AddFriend(nickname, friendsNickname string) error
	GetFriends(nickname string) ([]UserRow, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (repository *userRepository) SaveUser(user UserRow) error {
	sql := `
	INSERT INTO messenger.users(nickname, name, password)
	VALUES($1, $2, $3);`
	_, err := repository.db.Exec(sql, user.Nickname, user.Name, user.HashedPassword)
	return err
}

func (repository *userRepository) FindUsersByNickname(nicknameSubstring string) ([]UserRow, error) {
	sql := `
	SELECT name, nickname, password FROM messenger.users
	WHERE nickname = $1;`
	rows, err := repository.db.Query(sql, nicknameSubstring)
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

func (repository *userRepository) AddFriend(nickname, friendsNickname string) error {
	sql := `INSERT INTO messenger.friends(nickname, friends_nickname)
	VALUES ($1, $2)`
	_, err := repository.db.Exec(sql, nickname, friendsNickname)
	return err
}

func (repository *userRepository) GetFriends(nickname string) ([]UserRow, error) {
	sql := `
		SELECT u.name, u.nickname FROM messenger.users u
		JOIN messenger.friends f ON u.nickname = f.friends_nickname
		WHERE f.nickname = $1;`
	rows, err := repository.db.Query(sql, nickname)
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