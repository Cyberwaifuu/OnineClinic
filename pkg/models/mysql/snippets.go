package mysql

import (
	"codes/pkg/models"
	"database/sql"
	"errors"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires, category string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires, category)
             VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY), ?)`
	result, err := m.DB.Exec(stmt, title, content, expires, category)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires, category FROM snippets
WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.Category)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires, category FROM snippets
             WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	snippets := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.Category)

		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
func (m *SnippetModel) LastCategory(category string) ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created,expires,  category FROM 
    snippets WHERE category = ? ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt, category)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*models.Snippet{}
	for rows.Next() {
		s := &models.Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.Category)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
