package main

import (
	"github.com/kfelter/jumpsuit"
)

type Order struct {
	ID        int64
	ItemID    int64
	ItemName  string
	SalePrice int64
	Quantity  int64
	Color     string
}

func main() {
	fstore := jumpsuit.NewFileStore("data.json")
	fstore.Put(0, Order{
		ID:        0,
		ItemID:    1,
		ItemName:  "shorts",
		SalePrice: 1099,
		Quantity:  1,
	})
	server := jumpsuit.New(&jumpsuit.ServerOpts{Storage: fstore})
	server.NewAPI(Order{}, jumpsuit.APIOptions{})
	if err := server.Start(":8080"); err != nil {
		panic(err)
	}
}
