package item

import "github.com/efimovad/avito-internship/internal/model"

type Usecase interface {
	Create(item *model.Item) error
	Get(id int64, allInfo bool) (*model.Item, error)
	List(params model.Params) ([]model.Item, error)
}
