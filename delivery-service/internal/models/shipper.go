package models

type Shipper struct {
	ID              int64  `db:"id"`
	Password        string `db:"password"`
	Email           string `db:"email"`
	Role            string `db:"role"`
	Name            string `db:"name"`
	Gender          string `db:"gender"`
	Phone           string `db:"phone"`
	VehicleType     string `db:"vehicle_type"`
	VehiclePlate    string `db:"vehicle_plate"`
	TotalDeliveries int    `db:"total_deliveries"`
	Status          string `db:"status"`
}

type ShipperRequest struct {
	ID           int64  `json:"userId"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	Name         string `json:"name" binding:"required"`
	Gender       string `json:"gender" binding:"required"`
	Phone        string `json:"phone" binding:"required"`
	VehicleType  string `json:"vehicleType" binding:"required"`
	VehiclePlate string `json:"vehiclePlate" binding:"required"`
}

type ShipperResponse struct {
	ID              int64  `json:"userId" db:"id"`
	Email           string `json:"email" db:"email"`
	Name            string `json:"name" db:"name"`
	Gender          string `json:"gender" db:"gender"`
	Phone           string `json:"phone" db:"phone"`
	Role            string `json:"role" db:"role"`
	VehicleType     string `json:"vehicleType" db:"vehicle_type"`
	VehiclePlate    string `json:"vehiclePlate" db:"vehicle_plate"`
	TotalDeliveries int    `json:"totalDeliveries" db:"total_deliveries"`
	Status          string `json:"status" db:"status"`
}
