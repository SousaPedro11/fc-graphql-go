package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Category struct {
	db          *sql.DB
	ID          string
	Name        string
	Description string
}

func NewCategory(db *sql.DB) *Category {
	return &Category{db: db}
}

func (c *Category) Create(name string, description string) (Category, error) {
	if name == "" {
		return Category{}, errors.New("name is required")
	}
	id := uuid.New().String()

	stmt, err := c.db.Prepare("INSERT INTO categories (id, name, description) VALUES ($1, $2, $3)")

	if err != nil {
		return Category{}, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, name, description)
	if err != nil {
		return Category{}, fmt.Errorf("error executing statement: %w", err)
	}
	return Category{ID: id, Name: name, Description: description}, nil
}

func (c *Category) List() ([]Category, error) {
	rows, err := c.db.Query("SELECT id, name, description FROM categories")
	if err != nil {
		return nil, fmt.Errorf("error querying categories: %w", err)
	}
	defer rows.Close()

	categories := make([]Category, 0)
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description)
		if err != nil {
			return nil, fmt.Errorf("error scanning category: %w", err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading categories: %w", err)
	}

	return categories, nil
}
