package main

import (
	"github.com/efimovad/avito-internship/internal/app"
)

func main() {
	err := app.Start()
	println(err)
}
