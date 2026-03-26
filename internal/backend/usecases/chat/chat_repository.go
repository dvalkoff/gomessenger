package chat

import (
	"context"
	"database/sql"
	"time"
)

const(
	 txTimeout time.Duration = 1 * time.Second
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
	CreateChat(ctx context.Context, createChatInfo ChatRow) (ChatRow, error)
	AddUsersToChat(ctx context.Context, chatId int, users []ChatUserRow) error
	GetChat(ctx context.Context, chatId int) (ChatRow, error)
	GetNicknamesByChatId(ctx context.Context, chatId int) ([]string, error)
	GetChatIds(ctx context.Context) ([]int, error)
	GetChatsNoUsersByNickname(ctx context.Context, nickname string) ([]ChatRow, error)
}

type chatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) ChatRepository {
	return &chatRepository{db: db}
}

func (repository *chatRepository) CreateChat(ctx context.Context, createChatInfo ChatRow) (ChatRow, error) {
	ctx, cancel := context.WithTimeout(ctx, txTimeout)
	defer cancel()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return ChatRow{}, err
	}
	defer tx.Rollback()

	err = repository.insertChatAndGetId(ctx, tx, &createChatInfo)
	if err != nil {
		return ChatRow{}, err
	}

	err = repository.addUsersToChatInTx(ctx, tx, createChatInfo.id, createChatInfo.users)
	if err != nil {
		return ChatRow{}, err
	}
	if err := tx.Commit(); err != nil {
		return ChatRow{}, err
	}
	return createChatInfo, nil
}

func (repository *chatRepository) insertChatAndGetId(ctx context.Context, tx *sql.Tx, chat *ChatRow) error {
	sql := `INSERT INTO messenger.chats(name) VALUES ($1) RETURNING id`
	err := tx.QueryRowContext(ctx, sql, chat.name).Scan(&chat.id)
	return err
}

func (repository *chatRepository) AddUsersToChat(ctx context.Context, chatId int, users []ChatUserRow) error {
	ctx, cancel := context.WithTimeout(ctx, txTimeout)
	defer cancel()
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = repository.addUsersToChatInTx(ctx, tx, chatId, users)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (repository *chatRepository) addUsersToChatInTx(ctx context.Context, tx *sql.Tx, chatId int, users []ChatUserRow) error {
	sql := `INSERT INTO messenger.chats_users(user_nickname, chat_id, role)
	VALUES($1, $2, $3)`
	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		return err
	}
	for _, user := range users {
		if _, err := stmt.ExecContext(ctx, user.nickname, chatId, user.role); err != nil {
			return err
		}
	}
	return nil
}

func (repository *chatRepository) GetChat(ctx context.Context, chatId int) (ChatRow, error) {
	ctx, cancel := context.WithTimeout(ctx, txTimeout)
	defer cancel()
	getChatSql := `SELECT id, name FROM messenger.chats WHERE id = $1`
	chat := ChatRow{}
	err := repository.db.QueryRowContext(ctx, getChatSql, chatId).Scan(&chat.id, &chat.name)
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

func (repository *chatRepository) GetNicknamesByChatId(ctx context.Context, chatId int) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, txTimeout)
	defer cancel()
	sql := `SELECT user_nickname FROM messenger.chats_users WHERE chat_id = $1`
	rows, err := repository.db.QueryContext(ctx, sql, chatId)
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

func (repository *chatRepository) GetChatIds(ctx context.Context) ([]int, error) {
	ctx, cancel := context.WithTimeout(ctx, txTimeout)
	defer cancel()
	sql := `SELECT id FROM messenger.chats`
	rows, err := repository.db.QueryContext(ctx, sql)
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

func (repository *chatRepository) GetChatsNoUsersByNickname(ctx context.Context, nickname string) ([]ChatRow, error) {
	ctx, cancel := context.WithTimeout(ctx, txTimeout)
	defer cancel()
	sql := `SELECT c.id, c.name FROM messenger.chats c
			JOIN messenger.chats_users cu ON c.id = cu.chat_id
			WHERE cu.user_nickname = $1`
	rows, err := repository.db.QueryContext(ctx, sql, nickname)
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
