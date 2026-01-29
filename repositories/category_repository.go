package repositories

import (
	"category-api/models"
	"database/sql"
	"errors"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// get all categories
// GET
func (repo *CategoryRepository) GetAll() ([]models.Category, error) {
	query := "SELECT id, name, description FROM categories"
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.Name, &c.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// create new category
// POST
func (repo *CategoryRepository) Create(category *models.Category) error {
	query := "INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id"
	err := repo.db.QueryRow(query, category.Name, category.Description).Scan(&category.ID)
	return err
}

// get category by id
// GET
func (repo *CategoryRepository) GetByID(id int) (*models.Category, error) {
	query := "SELECT id, name, description FROM categories WHERE id = $1"
	var c models.Category
	err := repo.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.Description)
	if err == sql.ErrNoRows {
		return nil, errors.New("category not found")
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// update category
// PUT
func (repo *CategoryRepository) Update(category *models.Category) error {
	query := "UPDATE categories SET name = $1, description = $2 WHERE id = $3"
	result, err := repo.db.Exec(query, category.Name, category.Description, category.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("category not found")
	}
	return nil
}

// delete category
// DELETE
func (repo *CategoryRepository) Delete(id int) error {
	query := "DELETE FROM categories WHERE id = $1"

	res, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("category not found")
	}
	return err
}

// get category by id
// GET
func (repo *CategoryRepository) GetByIDWithProduct(id int) (*models.Category, error) {
	query := `
		SELECT 
			cs.id,
			cs.name,
			cs.description,
			pr.id,
			pr.name,
			pr.price,
			pr.stock
		FROM categories cs
		LEFT JOIN products pr ON pr.category_id = cs.id
		WHERE cs.id = $1
	`

	rows, err := repo.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var category *models.Category

	for rows.Next() {
		var (
			c            models.Category
			productID    sql.NullInt64
			productName  sql.NullString
			productPrice sql.NullInt64
			productStock sql.NullInt64
		)

		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Description,
			&productID,
			&productName,
			&productPrice,
			&productStock,
		)
		if err != nil {
			return nil, err
		}

		// init category sekali
		if category == nil {
			category = &models.Category{
				ID:          c.ID,
				Name:        c.Name,
				Description: c.Description,
				Products:    []models.Product{},
			}
		}

		// append product hanya jika ada
		if productID.Valid {
			category.Products = append(category.Products, models.Product{
				ID:    int(productID.Int64),
				Name:  productName.String,
				Price: int(productPrice.Int64),
				Stock: int(productStock.Int64),
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if category == nil {
		return nil, errors.New("category not found")
	}

	return category, nil
}
