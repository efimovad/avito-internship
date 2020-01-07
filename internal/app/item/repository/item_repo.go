package item_repo

import (
	"database/sql"
	"github.com/efimovad/avito-internship/internal/app/item"
	"github.com/efimovad/avito-internship/internal/model"
	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) item.Repository {
	return &Repository{db}
}

func (r *Repository) Create(item *model.Item) error {
	return r.db.QueryRow(
		`INSERT INTO items (title, description, date, price, images) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		item.Title,
		item.Description,
		item.Date,
		item.Price,
		pq.Array(item.Images),
	).Scan(&item.ID)
}

func (r *Repository) GetAll(id int64) (*model.Item, error) {
	var arr pq.StringArray

	myItem := &model.Item{}
	myItem.Images = make([]string, 0)
	if err := r.db.QueryRow(
		`SELECT id, title, price, description, images FROM items WHERE id = $1`,
		id,
	).Scan(
		&myItem.ID,
		&myItem.Title,
		&myItem.Price,
		&myItem.Description,
		&arr,
	); err != nil {
		return nil, err
	}
	for _, image := range arr {
		myItem.Images = append(myItem.Images, image)
	}
	return myItem, nil
}

func (r *Repository) Get(id int64) (*model.Item, error) {
	myItem := &model.Item{}

	if err := r.db.QueryRow(
		`SELECT id, title, price, images[1] FROM items WHERE id = $1`,
		id,
	).Scan(
		&myItem.ID,
		&myItem.Title,
		&myItem.Price,
		&myItem.MainImage,
	); err != nil {
		return nil, err
	}
	return myItem, nil
}

func (r *Repository) List(params model.Params) ([]model.Item, error) {
	var items []model.Item
	rows, err := r.db.Query(`
		 SELECT 
       			title, 
		        images[1] AS images, 
		        price 	
		 FROM 
		      items 
		 ORDER BY 
		          CASE WHEN ($1 AND $2) THEN date END DESC,
		          CASE WHEN ($1 AND $3) THEN price END DESC,
		          CASE WHEN (NOT $1 AND $2) THEN date END ASC, 
		          CASE WHEN (NOT $1 AND $3) THEN price END ASC
		 LIMIT 10 OFFSET CASE WHEN $4 > 0 THEN ($4 - 1) * 10 END;
	`, params.Desc, params.Date, params.Price, params.Page,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		i := model.Item{}
		err := rows.Scan(&i.Title, &i.MainImage, &i.Price)

		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	return items, nil
}