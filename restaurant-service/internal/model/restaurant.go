package model

type Restaurant struct {
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	Address     string `db:"address" json:"address"`
	PhoneNumber string `db:"phone_number" json:"phoneNumber"`
	IsActive    bool   `db:"is_active" json:"isActive"`
	OpenTime    string `db:"open_time" json:"openTime"`
	CloseTime   string `db:"close_time" json:"closeTime"`
}
