package chat

import (
	"database/sql"
)

type ChatRow struct {
	id    int
	name  string
	users []ChatUserRow
}

type ChatUserRow struct {
	nickname string
	role     string
}

type ChatRepository interface {
	CreateChat(createChatInfo ChatRow) (ChatRow, error)
	AddUsersToChat(chatId int, users []ChatUserRow) error
	GetChat(chatId int) (ChatRow, error)
	GetNicknamesByChatId(chatId int) ([]string, error)
	GetChatIds() ([]int, error)
	GetChatsNoUsersByNickname(nickname string) ([]ChatRow, error)
}

type chatRepository struct {
	db *sql.DB
}

// https://blog.thibaut-rousseau.com/blog/sql-transactions-in-go-the-good-way/

func NewChatRepository(db *sql.DB) ChatRepository {
	return &chatRepository{db: db}
}

func (repository *chatRepository) CreateChat(createChatInfo ChatRow) (ChatRow, error) {
	tx, err := repository.db.Begin()
	if err != nil {
		return ChatRow{}, err
	}
	defer tx.Commit()

	err = repository.insertChatAndGetId(tx, &createChatInfo)
	if err != nil {
		tx.Rollback()
		return ChatRow{}, err
	}

	err = repository.addUsersToChatInTx(tx, createChatInfo.id, createChatInfo.users)
	if err != nil {
		tx.Rollback()
		return ChatRow{}, err
	}
	return createChatInfo, nil
}

func (repository *chatRepository) insertChatAndGetId(tx *sql.Tx, chat *ChatRow) error {
	sql := `INSERT INTO messenger.chats(name) VALUES ($1) RETURNING id`
	err := tx.QueryRow(sql, chat.name).Scan(&chat.id)
	return err
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

func (repository *chatRepository) GetChat(chatId int) (ChatRow, error) {
	getChatSql := `SELECT id, name FROM messenger.chats WHERE id = $1`
	chat := ChatRow{}
	err := repository.db.QueryRow(getChatSql, chatId).Scan(&chat.id, &chat.name)
	if err != nil {
		return ChatRow{}, err
	}
	getChatUsersSql := `SELECT user_nickname, role FROM messenger.chats_users WHERE chat_id = $1`
	rows, err := repository.db.Query(getChatUsersSql, chatId)
	if err != nil {
		return ChatRow{}, err
	}
	users := []ChatUserRow{}
	for rows.Next() {
		user := ChatUserRow{}
		rows.Scan(&user.nickname, &user.role)
		users = append(users, user)
	}
	chat.users = users
	if rows.Err() != nil {
		return ChatRow{}, err
	}
	return chat, nil
}

func (repository *chatRepository) GetNicknamesByChatId(chatId int) ([]string, error) {
	sql := `SELECT user_nickname FROM messenger.chats_users WHERE chat_id = $1`
	rows, err := repository.db.Query(sql, chatId)
	if err != nil {
		return nil, err
	}
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

func (repository *chatRepository) GetChatsNoUsersByNickname(nickname string) ([]ChatRow, error) {
	sql := `SELECT c.id, c.name FROM messenger.chats c
			JOIN messenger.chats_users cu ON c.id = cu.chat_id
			WHERE cu.user_nickname = $1`
	rows, err := repository.db.Query(sql, nickname)
	if err != nil {
		return nil, err
	}
	chats := make([]ChatRow, 0)
	for rows.Next() {
		chatRow := ChatRow{}
		rows.Scan(&chatRow.id, &chatRow.name)
		chats = append(chats, chatRow)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return chats, nil
}
