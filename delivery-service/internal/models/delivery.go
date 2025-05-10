package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Custom type for scanning JSON/text array from DB to []float64
type Float64Slice []float64

func (f *Float64Slice) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return json.Unmarshal(src, f)
	case string:
		return json.Unmarshal([]byte(src), f)
	}
	return fmt.Errorf("unsupported type for Float64Slice: %T", src)
}

func (f Float64Slice) Value() (driver.Value, error) {
	return json.Marshal(f)
}

type CreateDeliveryRequest struct {
	OrderID           int64  `json:"orderId"`
	RestaurantAddress string `json:"restaurantAddress"`
	ShippingAddress   string `json:"shippingAddress"`
}

type DeliveryResponse struct {
	DeliveryID   int64           `json:"deliveryId"`
	OrderID      int64           `json:"orderId"`
	Distance     float64         `json:"distance"`
	Duration     float64         `json:"duration"`
	Fee          int64           `json:"fee"`
	FromCoords   Float64Slice    `json:"fromCoords"`
	ToCoords     Float64Slice    `json:"toCoords"`
	GeometryLine string          `json:"geometryLine"`
	Status       string          `json:"status"`
	Shipper      ShipperResponse `json:"shipper"`
}

type Delivery struct {
	DeliveryID        int64        `db:"delivery_id" json:"deliveryId"`
	OrderID           int64        `db:"order_id" json:"orderId"`
	ShipperID         int64        `db:"shipper_id" json:"shipperId"`
	RestaurantAddress string       `db:"restaurant_address" json:"restaurantAddress"`
	ShippingAddress   string       `db:"shipping_address" json:"shippingAddress"`
	Distance          float64      `db:"distance" json:"distance"`
	Duration          float64      `db:"duration" json:"duration"`
	Fee               int64        `db:"fee" json:"fee"`
	FromCoords        Float64Slice `db:"from_coords" json:"fromCoords"`
	ToCoords          Float64Slice `db:"to_coords" json:"toCoords"`
	GeometryLine      string       `db:"geometry_line" json:"geometryLine"`
	Status            string       `db:"status" json:"status"`
	CreatedAt         time.Time    `db:"created_at" json:"createdAt"`
	UpdatedAt         time.Time    `db:"updated_at" json:"updatedAt"`
}

type LocationRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type DistanceResponse struct {
	Distance     float64   `json:"distance"`
	Duration     float64   `json:"duration"`
	Fee          int64     `json:"fee"`
	FromCoords   []float64 `json:"fromCoords"`
	ToCoords     []float64 `json:"toCoords"`
	GeometryLine string    `json:"geometryLine,omitempty"`
}

type DeliveryGetByShipperId struct {
	DeliveryID   int64        `db:"delivery_id" json:"deliveryId"`
	OrderID      int64        `db:"order_id" json:"orderId"`
	Distance     float64      `db:"distance" json:"distance"`
	Duration     float64      `db:"duration" json:"duration"`
	Fee          int64        `db:"fee" json:"fee"`
	FromCoords   Float64Slice `db:"from_coords" json:"fromCoords"`
	ToCoords     Float64Slice `db:"to_coords" json:"toCoords"`
	GeometryLine string       `db:"geometry_line" json:"geometryLine"`
	Status       string       `db:"status" json:"status"`
}
