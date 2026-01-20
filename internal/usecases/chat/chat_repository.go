package chat

import "database/sql"

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
	CreateChat(name string, creatorNickname string) (*ChatRow, error)
	AddUserToChat(tx *sql.Tx, chatId int, nickname string, role string) (ChatUserRow, error)
	GetChat(chatId int) (*ChatRow, error)
	GetChatIdsByUser(nickname string) ([]int, error)
	GetNicknamesByChatId(chatId int) ([]string, error)
	GetChatIds() ([]int, error)
}

type chatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) ChatRepository {
	return &chatRepository{db: db}
}

func (repository *chatRepository) CreateChat(name string, creatorNickname string) (*ChatRow, error) {
	tx, err := repository.db.Begin() // TODO: how to I specify timeout?
	if err != nil {
		return nil, err
	}
	defer tx.Commit() // TODO: should I close tx this way? should I catch err?

	chatId, err := repository.insertChatAndGetId(tx, name)
	if err != nil {
		return nil, err
	}
	chatCreatorRole := "admin"
	chatUserRow, err := repository.AddUserToChat(tx, chatId, creatorNickname, chatCreatorRole)
	if err != nil {
		return nil, err
	}
	return &ChatRow{
		id: int(chatId),
		name: name,
		users: []ChatUserRow{chatUserRow},
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

func (repository *chatRepository) AddUserToChat(tx *sql.Tx, chatId int, nickname string, role string) (ChatUserRow, error) {
	sql := `INSERT INTO messenger.chats_users(user_nickname, chat_id, role)
	VALUES($1, $2, $3)`
	var err error
	if tx != nil {
		_, err = tx.Exec(sql, nickname, chatId, role) 
	} else {
		_, err = repository.db.Exec(sql, nickname, chatId, role)
	}
	if err != nil {
		return ChatUserRow{}, err
	}
	return ChatUserRow{nickname: nickname, role: role}, nil
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