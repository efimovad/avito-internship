package model

import "time"

type Item struct {
	ID int64 `json:"-"`
	Title string `json:"title" validate:"required,min=1,max=200"`
	Description string `json:"description,omitempty" validate:"required,min=1,max=1000"`
	Date time.Time `json:"-"`
	Price float32 `json:"price"`
	MainImage string `json:"mainImage,omitempty"`
	Images []string `json:"images,omitempty" validate:"min=0,max=3"`
}

type Params struct {
	Date bool
	Price bool
	Desc bool
	Page int64
}