package model

import "time"

type Item struct {
	ID int64 `json:"-"`
	Title string `json:"title"`
	Description string `json:"description"`
	Date time.Time `json:"-"`
	Price float32 `json:"price"`
	MainImage string `json:"mainImage"`
	Images [3]string `json:"images"`
}

type Params struct {
	Date bool
	Price bool
	Desc bool
	Page int64
}