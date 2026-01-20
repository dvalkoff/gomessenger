package chat

import (
	"database/sql"
)

type ChatRow struct {
	id int
	name string
	users []ChatUserRow
}

type ChatUserRow struct {
	nickname string
	role string
}

type ChatRepository interface {
	CreateChat(createChatInfo CreateChatInfo) (*ChatRow, error)
	AddUsersToChat(chatId int, users []ChatUserRow) (error)
	GetChat(chatId int) (*ChatRow, error)
	GetChatIdsByUser(nickname string) ([]int, error)
	GetNicknamesByChatId(chatId int) ([]string, error)
	GetChatIds() ([]int, error)
}

type chatRepository struct {
	db *sql.DB
}
// https://blog.thibaut-rousseau.com/blog/sql-transactions-in-go-the-good-way/

func NewChatRepository(db *sql.DB) ChatRepository {
	return &chatRepository{db: db}
}

func (repository *chatRepository) CreateChat(createChatInfo CreateChatInfo) (*ChatRow, error) {
	tx, err := repository.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Commit()

	chatId, err := repository.insertChatAndGetId(tx, createChatInfo.Name)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	usersToAdd := make([]ChatUserRow, 0, len(createChatInfo.Users) + 1)
	usersToAdd = append(usersToAdd, ChatUserRow{createChatInfo.CreatorNickname, "admin"})
	for _, userToAdd := range createChatInfo.Users {
		usersToAdd = append(usersToAdd, ChatUserRow{userToAdd, "user"})
	}

	err = repository.addUsersToChatInTx(tx, chatId, usersToAdd)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return &ChatRow{
		id: int(chatId),
		name: createChatInfo.Name,
		users: usersToAdd,
	}, nil
}

func (repository *chatRepository) insertChatAndGetId(tx *sql.Tx, name string) (int, error) {
	sql := `INSERT INTO messenger.chats(name) VALUES ($1) RETURNING id`
	var chatId int
	err := tx.QueryRow(sql, name).Scan(&chatId)
	if err != nil {
		return 0, err
	}
	return chatId, nil
}

func (repository *chatRepository) AddUsersToChat(chatId int, users []ChatUserRow) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	err = repository.addUsersToChatInTx(tx, chatId, users)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (repository *chatRepository) addUsersToChatInTx(tx *sql.Tx, chatId int, users []ChatUserRow) error {
	sql := `INSERT INTO messenger.chats_users(user_nickname, chat_id, role)
	VALUES($1, $2, $3)`
	stmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}
	for _, user := range users {
		if _, err := stmt.Exec(user.nickname, chatId, user.role); err != nil {
			return err
		}
	}
	return nil
}

func (repository *chatRepository) GetChat(chatId int) (*ChatRow, error) {
	getChatSql := `SELECT id, name FROM messenger.chats WHERE id = $1`
	chat := ChatRow{}
	err := repository.db.QueryRow(getChatSql, chatId).Scan(&chat.id, &chat.name)
	if err != nil {
		return nil, err
	}
	getChatUsersSql := `SELECT user_nickname, role FROM messenger.chats_users WHERE chat_id = $1`
	rows, err := repository.db.Query(getChatUsersSql, chatId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []ChatUserRow{}
	for rows.Next() {
		user := ChatUserRow{}
		err := rows.Scan(&user.nickname, &user.role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	chat.users = users
	if rows.Err() != nil {
		return nil, err
	}
	return &chat, nil
}

func (repository *chatRepository) GetChatIdsByUser(nickname string) ([]int, error) {
	sql := `SELECT chat_id FROM messenger.chats_users WHERE user_nickname = $1`
	rows, err := repository.db.Query(sql, nickname)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	chatIds := make([]int, 0)
	for rows.Next() {
		chatId := 0
		rows.Scan(&chatId)
		chatIds = append(chatIds, chatId)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return chatIds, nil
}

func (repository *chatRepository) GetNicknamesByChatId(chatId int) ([]string, error) {
	sql := `SELECT user_nickname FROM messenger.chats_users WHERE chat_id = $1`
	rows, err := repository.db.Query(sql, chatId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	nicknames := make([]string, 0)
	for rows.Next() {
		var nickname string
		rows.Scan(&nickname)
		nicknames = append(nicknames, nickname)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return nicknames, nil
}

func (repository *chatRepository) GetChatIds() ([]int, error) {
	sql := `SELECT id FROM messenger.chats`
	rows, err := repository.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	chatIds := make([]int, 0)
	for rows.Next() {
		chatId := 0
		rows.Scan(&chatId)
		chatIds = append(chatIds, chatId)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return chatIds, nil
}