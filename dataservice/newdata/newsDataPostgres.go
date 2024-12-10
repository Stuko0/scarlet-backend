package newdata

import (
	"database/sql"
	"scarlet_backend/model"
)

type NewsDataPostgres struct{ db *sql.DB }

func NewNewsDataPostgres(db *sql.DB) *NewsDataPostgres { return &NewsDataPostgres{db: db} }

func (n *NewsDataPostgres) GetNews() ([]model.New, error) {
	rows, err := n.db.Query("select newId,title, description,image,url,active,created_at,created_by from new_data.news")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var news []model.New
	for rows.Next() {
		var newItem model.New
		if err = rows.Scan(&newItem.Title, &newItem.Description, &newItem.Image, &newItem.URL, &newItem.Active, &newItem.CreatedAt, &newItem.CreatedBy); err != nil {
			return nil, err
		}
		news = append(news, newItem)
	}
	return news, nil
}

func (n *NewsDataPostgres) AddNew(newItem model.New) error {
	query := `insert into new_data.news (title,description,image,url,created_at,created_by)values($1,$2,$3,$4,$5,$6)`
	_, err := n.db.Exec(query, newItem.Title, newItem.Description, newItem.Image, newItem.URL, newItem.CreatedAt, newItem.CreatedBy)
	return err
}
