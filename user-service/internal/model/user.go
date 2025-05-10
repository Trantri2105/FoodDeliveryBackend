package model

const (
	ADMIN    = "admin"
	CUSTOMER = "customer"
	SHIPPER  = "shipper"
)

type User struct {
	UserId   int    `db:"user_id" json:"userId,omitempty"`
	Email    string `db:"email" json:"email,omitempty"`
	Password string `db:"password" json:"password,omitempty"`
	Name     string `db:"name" json:"name,omitempty"`
	Gender   string `db:"gender" json:"gender,omitempty"`
	Phone    string `db:"phone" json:"phone,omitempty"`
	Role     string `db:"role" json:"role,omitempty"`
}
