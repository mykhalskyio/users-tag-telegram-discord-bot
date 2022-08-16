package entity

type ChatUser struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
	ChatId   int    `db:"chat_id"`
}
