package main

import (
	"github.com/kfelter/jumpsuit"
)

type Order struct {
	ID       int64
	ItemID   int64
	Quantity int64
	Total    int64
}

type Item struct {
	ID    int64
	Name  string
	Price int64
	Color string
}

func main() {
	fstore := jumpsuit.NewFileStore("data")
	fstore.Put("orders", &Order{
		ID:       -1,
		ItemID:   0,
		Quantity: 1,
		Total:    1099 * 1,
	})

	fstore.Put("items", &Item{
		ID:    -1,
		Name:  "shorts",
		Price: 1099,
		Color: "blue",
	})
	server := jumpsuit.New(&jumpsuit.ServerOpts{Storage: fstore})
	server.NewAPI(Order{}, jumpsuit.APIOptions{})
	server.NewAPI(Item{}, jumpsuit.APIOptions{})
	if err := server.Start(":8080"); err != nil {
		panic(err)
	}
}
