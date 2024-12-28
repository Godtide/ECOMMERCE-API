package models

import "time"

type User struct {
	ID             uint   `gorm:"primaryKey"`
	Name           string `json:"name"`
	Email          string `json:"email" gorm:"unique"`
	HashedPassword string `json:"-"`
	Role           string `json:"role"` // admin or user
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Product struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Order struct {
	ID          uint           `gorm:"primaryKey"`
	UserID      uint           `json:"user_id"`
	User        User           `gorm:"foreignKey:UserID"`
	Products    []OrderProduct `gorm:"foreignKey:OrderID" json:"products"`
	Status      string         `json:"status" gorm:"default:'Pending'"` // Pending, Completed, Cancelled
	TotalAmount float64        `json:"total_amount"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type OrderProduct struct {
	ID        uint    `gorm:"primaryKey"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID"`
	Quantity  int     `json:"quantity"`
}
