package routes

import (
	"ecommerce-api/controllers"
	"ecommerce-api/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes sets up all the API routes
func RegisterRoutes() *gin.Engine {
	router := gin.Default()

	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// User routes
	userRoutes := router.Group("/users")
	{
		// @Summary Register a new user
		// @Description Create a new user account
		// @Tags User
		// @Accept json
		// @Produce json
		// @Param user body models.User true "User details"
		// @Success 200 {object} models.User
		// @Router /users/register [post]
		userRoutes.POST("/register", controllers.RegisterUser)

		// @Summary User login
		// @Description Authenticate a user and return a JWT token
		// @Tags User
		// @Accept json
		// @Produce json
		// @Param credentials body models.LoginRequest true "Login details"
		// @Success 200 {object} models.LoginResponse
		// @Router /users/login [post]
		userRoutes.POST("/login", controllers.LoginUser)
	}

	// Product routes (requires admin privileges)
	productRoutes := router.Group("/products")
	productRoutes.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// CreateProduct creates a new product in the system. Requires admin privileges.
		// @Summary Create a new product
		// @Description Allows an admin to create a new product by providing product details.
		// @Tags Products
		// @Accept json
		// @Produce json
		// @Param product body models.Product true "Product details"
		// @Success 201 {object} models.Product "Product successfully created"
		// @Failure 400 {object} gin.H "Bad Request: Invalid product data"
		// @Failure 401 {object} gin.H "Unauthorized: Authentication required"
		// @Failure 403 {object} gin.H "Forbidden: Admin privileges required"
		// @Router /products [post]
		productRoutes.POST("/", controllers.CreateProduct)

		// ListProducts retrieves all products in the system.
		// @Summary List all products
		// @Description Retrieves a list of all available products.
		// @Tags Products
		// @Accept json
		// @Produce json
		// @Success 200 {array} models.Product "List of products"
		// @Failure 500 {object} gin.H "Internal Server Error"
		// @Router /products [get]
		productRoutes.GET("/", controllers.ListProducts)

		// GetProduct retrieves a specific product by its ID.
		// @Summary Get a single product
		// @Description Retrieve detailed information about a specific product by its ID.
		// @Tags Products
		// @Accept json
		// @Produce json
		// @Param id path int true "Product ID"
		// @Success 200 {object} models.Product "Product details"
		// @Failure 400 {object} gin.H "Bad Request: Invalid product ID"
		// @Failure 404 {object} gin.H "Product not found"
		// @Router /products/{id} [get]
		productRoutes.GET("/:id", controllers.GetProduct)
		// UpdateProduct updates an existing product's details. Requires admin privileges.
		// @Summary Update a product's details
		// @Description Allows an admin to update the details of an existing product by providing the product's ID and new information.
		// @Tags Products
		// @Accept json
		// @Produce json
		// @Param id path int true "Product ID"
		// @Param product body models.Product true "Updated product details"
		// @Success 200 {object} models.Product "Product successfully updated"
		// @Failure 400 {object} gin.H "Bad Request: Invalid product data"
		// @Failure 401 {object} gin.H "Unauthorized: Authentication required"
		// @Failure 403 {object} gin.H "Forbidden: Admin privileges required"
		// @Failure 404 {object} gin.H "Product not found"
		// @Router /products/{id} [put]
		productRoutes.PUT("/:id", controllers.UpdateProduct)
		// DeleteProduct deletes a product by its ID. Requires admin privileges.
		// @Summary Delete a product
		// @Description Allows an admin to delete a product by its ID.
		// @Tags Products
		// @Accept json
		// @Produce json
		// @Param id path int true "Product ID"
		// @Success 200 {object} gin.H "Product successfully deleted"
		// @Failure 401 {object} gin.H "Unauthorized: Authentication required"
		// @Failure 403 {object} gin.H "Forbidden: Admin privileges required"
		// @Failure 404 {object} gin.H "Product not found"
		// @Router /products/{id} [delete]
		productRoutes.DELETE("/:id", controllers.DeleteProduct)
	}

	// Order routes (requires user authentication)
	orderRoutes := router.Group("/orders")
	orderRoutes.Use(middleware.AuthMiddleware())
	{

		// PlaceOrder places a new order
		// @Summary Place a new order
		// @Description Place a new order for authenticated users
		// @Tags Order
		// @Accept json
		// @Produce json
		// @Param Authorization header string true "Bearer token"
		// @Param order body []struct{product_id uint, quantity int} true "Order Details"
		// @Success 201 {object} gin.H "Success"
		// @Failure 400 {object} gin.H "Bad Request"
		// @Failure 401 {object} gin.H "Unauthorized"
		// @Failure 404 {object} gin.H "Product Not Found"
		// @Router /orders [post]
		// @Security BearerAuth

		orderRoutes.POST("/", controllers.PlaceOrder)
		// ListUserOrders retrieves all orders for a specific user.
		// @Summary List all orders for a user
		// @Description Fetch all orders placed by the authenticated user. Requires user authentication.
		// @Tags Orders
		// @Accept json
		// @Produce json
		// @Success 200 {array} models.Order "List of user orders"
		// @Failure 401 {object} gin.H "Unauthorized: Authentication required"
		// @Router /orders [get]
		orderRoutes.GET("/", controllers.ListUserOrders)
		// GetOrder retrieves a specific order by its ID.
		// @Summary Get a single order
		// @Description Retrieve detailed information about a specific order by its ID. Requires user authentication.
		// @Tags Orders
		// @Accept json
		// @Produce json
		// @Param id path int true "Order ID"
		// @Success 200 {object} models.Order "Order details"
		// @Failure 400 {object} gin.H "Bad Request"
		// @Failure 401 {object} gin.H "Unauthorized"
		// @Failure 404 {object} gin.H "Order not found"
		// @Router /orders/{id} [get]
		orderRoutes.GET("/:id", controllers.GetOrder)
		// CancelOrder cancels an order if it is in the Pending status.
		// @Summary Cancel an order
		// @Description Allows a user to cancel an order if its status is still "Pending". Requires user authentication.
		// @Tags Orders
		// @Accept json
		// @Produce json
		// @Param id path int true "Order ID"
		// @Success 200 {object} gin.H "Order successfully canceled"
		// @Failure 400 {object} gin.H "Bad Request: Order is not in a cancellable state"
		// @Failure 401 {object} gin.H "Unauthorized: Authentication required"
		// @Failure 403 {object} gin.H "Forbidden: User does not own this order"
		// @Failure 404 {object} gin.H "Order not found"
		// @Router /orders/{id}/cancel [put]
		orderRoutes.PUT("/:id/cancel", controllers.CancelOrder)
		// UpdateOrderStatus updates the status of an order. Requires admin privileges.
		// @Summary Update the status of an order
		// @Description Allows an admin to update the status of an order (e.g., from "Pending" to "Shipped").
		// @Tags Orders
		// @Accept json
		// @Produce json
		// @Param id path int true "Order ID"
		// @Param status body models.OrderStatus true "New Order Status"
		// @Success 200 {object} models.Order "Order status successfully updated"
		// @Failure 400 {object} gin.H "Bad Request: Invalid status value"
		// @Failure 401 {object} gin.H "Unauthorized: Authentication required"
		// @Failure 403 {object} gin.H "Forbidden: Admin privileges required"
		// @Failure 404 {object} gin.H "Order not found"
		// @Router /orders/{id}/status [put]
		orderRoutes.PUT("/:id/status", middleware.AdminMiddleware(), controllers.UpdateOrderStatus)
	}

	return router
}
