package model

type MenuItem struct {
	Id          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	Price       int    `db:"price" json:"price"`
	IsAvailable bool   `db:"is_available" json:"isAvailable"`
	ImageUrl    string `db:"image_url" json:"imageUrl"`
}
