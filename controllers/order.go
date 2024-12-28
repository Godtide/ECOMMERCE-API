package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ecommerce-api/models"
)

func PlaceOrder(c *gin.Context) {
	var orderRequest struct {
		Products []struct {
			ProductID uint `json:"product_id"`
			Quantity  int  `json:"quantity"`
		} `json:"products"`
	}

	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var totalAmount float64
	var orderProducts []models.OrderProduct

	for _, p := range orderRequest.Products {
		var product models.Product
		if err := db.First(&product, p.ProductID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		if product.Stock < p.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
			return
		}

		product.Stock -= p.Quantity
		db.Save(&product)

		totalAmount += float64(p.Quantity) * product.Price
		orderProducts = append(orderProducts, models.OrderProduct{
			ProductID: p.ProductID,
			Quantity:  p.Quantity,
		})
	}

	userID := c.GetUint("user_id") // Extract from JWT claims
	order := models.Order{
		UserID:      userID,
		Products:    orderProducts,
		TotalAmount: totalAmount,
	}

	if err := db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to place order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": order})
}

func ListUserOrders(c *gin.Context) {
	userID := c.GetUint("user_id") // Extract from JWT claims
	var orders []models.Order

	if err := db.Preload("Products.Product").Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orders})
}

func GetOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	if err := db.Preload("Products.Product").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

func CancelOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	if err := db.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if order.Status != "Pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only pending orders can be cancelled"})
		return
	}

	order.Status = "Cancelled"
	if err := db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}

func UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	if err := db.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	var input struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Status != "Completed" && input.Status != "Cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	order.Status = input.Status
	if err := db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}
