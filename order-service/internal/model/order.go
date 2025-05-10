package model

type Order struct {
	Id              int         `db:"id" json:"id"`
	UserId          int         `db:"user_id" json:"user_id"`
	ShippingAddress string      `db:"shipping_address" json:"shippingAddress"`
	PhoneNumber     string      `db:"phone_number" json:"phoneNumber"`
	Status          string      `db:"status" json:"status"`
	Subtotal        int         `db:"subtotal" json:"subtotal"`
	DeliveryFee     int         `db:"delivery_fee" json:"deliveryFee"`
	TotalAmount     int         `db:"total_amount" json:"totalAmount"`
	OrderItems      []OrderItem `json:"orderItems,omitempty"`
	Delivery        *Delivery   `json:"delivery,omitempty"`
}

type OrderItem struct {
	MenuItemId int `db:"menu_item_id" json:"menuItemId"`
	Quantity   int `db:"quantity" json:"quantity"`
	UnitPrice  int `db:"unit_price" json:"unitPrice"`
	TotalPrice int `db:"total_price" json:"totalPrice"`
	OrderId    int `db:"order_id" json:"-"`
}

type Delivery struct {
	Distance     float64   `json:"distance,omitempty"`
	Duration     float64   `json:"duration,omitempty"`
	Fee          int       `json:"fee,omitempty"`
	FromCoords   []float64 `json:"fromCoords,omitempty"`
	ToCoords     []float64 `json:"toCoords,omitempty"`
	GeometryLine string    `json:"geometryLine,omitempty"`
	Status       string    `json:"status,omitempty"`
	Shipper      Shipper   `json:"shipper,omitempty"`
}

type Shipper struct {
	Name         string `json:"name,omitempty"`
	Gender       string `json:"gender,omitempty"`
	Phone        string `json:"phone,omitempty"`
	VehicleType  string `json:"vehicleType,omitempty"`
	VehiclePlate string `json:"vehiclePlate,omitempty"`
}
