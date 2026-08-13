package repository

import (
	"context"
	"translator/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type buttonRepository struct {
	db *pgxpool.Pool
}

func NewButtonRepository(db *pgxpool.Pool) ButtonRepository {
	return &buttonRepository{db: db}
}

func (b *buttonRepository) GetGroupByNameAndLang(name model.ButtonGroupName, lang model.Language) ([]model.Button, error) {
	query := `
		select id, language_id, group_name, data, title 
		from button 
		where group_name = $1 and language_id = $2
		order by sort_order
	`

	rows, err := b.db.Query(context.Background(), query, name, lang.ID)
	if err != nil {
		return nil, err
	}

	var buttons []model.Button
	for rows.Next() {
		var button model.Button

		if err := rows.Scan(&button.ID, &button.LanguageID, &button.GroupName, &button.Data, &button.Title); err != nil {
			return nil, err
		}
		buttons = append(buttons, button)
	}

	return buttons, rows.Err()
}
