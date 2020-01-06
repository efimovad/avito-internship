package item_ucase

import (
	"github.com/efimovad/avito-internship/internal/app/item"
	"github.com/efimovad/avito-internship/internal/model"
	"github.com/pkg/errors"
	"time"
)

type Usecase struct {
	repository        item.Repository
}

func NewItemUsecase(itemRep item.Repository) item.Usecase {
	return &Usecase{
		repository:        itemRep,
	}
}

func (u *Usecase) Create(item *model.Item) error {
	item.Date = time.Now().UTC()

	if err := u.repository.Create(item); err != nil {
		return errors.Wrap(err, "itemRepository.Create()")
	}
	return nil
}

func (u *Usecase) Get(id int64, allInfo bool) (*model.Item, error) {
	var myItem *model.Item
	var err error

	if allInfo {
		myItem, err = u.repository.GetAll(id)
	} else {
		myItem, err = u.repository.Get(id)
	}

	if err != nil {
		return nil, errors.Wrap(err, "itemRepository.Get()")
	}

	return myItem, nil
}

func (u *Usecase) List(params model.Params) ([]model.Item, error) {
	items, err := u.repository.List(params)
	if err != nil {
		return nil, errors.Wrap(err, "itemRepository.List()")
	}
	return items, nil
}