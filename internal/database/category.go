package database

import (
	"database/sql"
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
	id := uuid.New().String()
	_, err := c.db.Exec("INSERT INTO categories (id, name, description) VALUES ($1, $2, $3)",
		id, name, description)
	if err != nil {
		return Category{}, err
	}
	return Category{ID: id, Name: name, Description: description}, nil
}

func (c *Category) FindAll() ([]Category, error) {
	// Seleciona id e name from categories
	rows, err := c.db.Query("SELECT id, name, description FROM categories")
	// Se tiver algum erro nessa busca ele vai retornar
	if err != nil {
		return nil, err
	}
	// Fecha a conexão da linha que está em aberto
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	// Cria uma categoria em branco
	var categories []Category
	// Da um for nessa categoria criando de fato o objeto
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description)
		if err != nil {
			return nil, err
		}
		//Para cada objeto criado ele adiciona no array de categorias
		categories = append(categories, category)
	}
	return categories, nil
}

func (c *Category) Find(id string) (Category, error) {
	var name, description string
	err := c.db.QueryRow("SELECT name, description FROM categories WHERE id = $1", id).
		Scan(&name, &description)
	if err != nil {
		return Category{}, err
	}
	return Category{ID: id, Name: name, Description: description}, nil
}

func (c *Category) FindByCourseId(courseId string) (Category, error) {
	query := "SELECT c.id, c.name, c.description FROM categories c JOIN courses co ON c.id = co.category_id WHERE co.id = $1"
	var id, name, description string
	err := c.db.QueryRow(query, courseId).Scan(&id, &name, &description)
	if err != nil {
		return Category{}, err
	}
	return Category{ID: id, Name: name, Description: description}, nil
}
