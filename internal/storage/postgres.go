package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mykhalskyio/users-tag-telegram-discord-bot/internal/config"
	"github.com/mykhalskyio/users-tag-telegram-discord-bot/internal/entity"
)

type Postgres struct {
	DB *sqlx.DB
}

func NewPostgres(cfg *config.Config) (*Postgres, error) {
	pg, err := sqlx.Connect("postgres", fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d sslmode=%s ",
		cfg.Postgres.Name,
		cfg.Postgres.User,
		cfg.Postgres.Pass,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Sslmode))
	if err != nil {
		return nil, err
	}
	if err = pg.Ping(); err != nil {
		return nil, err
	}
	defer pg.Close()
	return &Postgres{
		DB: pg,
	}, nil
}

func (psql *Postgres) Insert(username string, chatId int) error {
	_, err := psql.DB.Exec("INSERT INTO bobik(username, chat_id) VALUES($1, $2);", username, chatId)
	if err != nil {
		return err
	}
	return nil
}

func (psql *Postgres) GetAll(chatId int) (*[]entity.ChatUser, bool, error) {
	users := []entity.ChatUser{}
	err := psql.DB.Select(&users, "SELECT * FROM bobik WHERE chat_id = $1;", chatId)
	if err != nil {
		return nil, false, err
	}
	return &users, true, nil
}

func (psql *Postgres) Get(username string, chatId int) (*entity.ChatUser, bool, error) {
	user := entity.ChatUser{}
	err := psql.DB.Get(&user, "SELECT * FROM bobik WHERE username = $1 AND chat_id = $2;", username, chatId)
	if err != nil {
		return nil, false, err
	}
	return &user, true, nil
}

func (psql *Postgres) Delete(username string, chatId int) error {
	_, err := psql.DB.Exec("DELETE FROM bobik WHERE username = $1 AND chat_id = $2;", username, chatId)
	if err != nil {
		return err
	}
	return nil
}
