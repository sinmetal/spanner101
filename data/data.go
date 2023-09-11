package data

import (
	"math/rand"
)

var allAuthor = []string{"gold", "silver", "dia", "ruby", "sapphire"}
var allItem = []*Item{
	{
		ItemID: "pen",
		Price:  100,
	},
	{
		ItemID: "ball",
		Price:  300,
	},
	{
		ItemID: "note",
		Price:  150,
	},
}

type Item struct {
	ItemID string
	Price  int64
}

// RandomUserID is randomに1人返す
func RandomUserID() string {
	return allAuthor[rand.Intn(len(allAuthor))]
}

// RandomItem is randomに1つ返す
func RandomItem() *Item {
	return allItem[rand.Intn(len(allItem))]
}
