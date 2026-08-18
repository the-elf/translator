package repository

import (
	"context"
	"errors"
	"translator/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) Get(chatID int64) (model.User, error) {
	query := `
		select u.id, u.chat_id, l.id, l.code from "user" u
		join language l on u.language_id = l.id
		where u.chat_id = $1
	`

	var user model.User
	err := u.db.QueryRow(context.Background(), query, chatID).Scan(
		&user.ID,
		&user.ChatID,
		&user.Language.ID,
		&user.Language.Code,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return user, model.ErrUserNotFound
	}
	return user, err
}

func (u *userRepository) Save(user model.User) error {
	query := `
		insert into "user" (chat_id, language_id)
		values ($1, $2)
		on conflict (chat_id)
		do update set language_id = excluded.language_id
	`
	_, err := u.db.Exec(context.Background(), query, user.ChatID, user.Language.ID)
	return err
}
