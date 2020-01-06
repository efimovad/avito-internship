package item

import "github.com/efimovad/avito-internship/internal/model"

type Repository interface {
	Create(item *model.Item) error
	Get(id int64) (*model.Item, error)
	GetAll(id int64) (*model.Item, error)
	List(params model.Params) ([]model.Item, error)
}
