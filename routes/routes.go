package routes

import (
	"go-crud/config"
	"go-crud/controllers"
	"go-crud/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes() *gin.Engine {
	r := gin.Default()

	// ========================================
	// CORS CONFIGURATION
	// ========================================
	// Apply CORS middleware globally
	r.Use(config.CORSConfig())

	// Disable automatic redirect from /path to /path/
	// This prevents issues with POST requests containing form data
	// r.RedirectTrailingSlash = false

	// ========================================
	// PUBLIC ENDPOINTS (No Auth Required)
	// ========================================
	// Health & Info endpoints
	r.GET("/health", controllers.HealthCheck)
	r.GET("/cors-test", controllers.CORSTest)
	r.GET("/api/info", controllers.APIInfo)

	// Authentication
	r.POST("/login", controllers.Login)
	r.POST("/logout", middlewares.AuthMiddleware(), controllers.Logout)

	r.Group("/books")
	{
		r.GET("/books/", controllers.ListBooks)
		r.POST("/books/", controllers.CreateBook)
		r.GET("/books/:id", controllers.GetBook)
		r.PUT("/books/:id", controllers.UpdateBook)
		r.DELETE("/books/:id", controllers.DeleteBook)
	}

	// Public tracking endpoint (no auth required)
	r.GET("/track/:trackingNumber", controllers.TrackShipment)

	r.Group("/authors")
	{
		r.GET("/authors/", controllers.ListAuthors)
		r.POST("/authors/", controllers.CreateAuthor)
		r.GET("/authors/:id", controllers.GetAuthor)
		r.PUT("/authors/:id", controllers.UpdateAuthor)
		r.DELETE("/authors/:id", controllers.DeleteAuthor)
	}

	// Public endpoint for accepting invitations
	r.POST("/company/team/users/accept-invitation", controllers.AcceptInvitation)

	// Public SaaS Auth Endpoints
	r.POST("/auth/signup", controllers.SaasSignup)
	r.POST("/auth/login", controllers.SaasLogin)
	r.POST("/auth/forgot-password", controllers.ForgotPassword)
	r.POST("/auth/reset-password", controllers.ResetPassword)

	// ========================================
	// PROTECTED ROUTES (Auth Required)
	// ========================================
	protected := r.Group("/", middlewares.AuthMiddleware())
	{
		users := protected.Group("/users")
		{
			users.GET("/", controllers.ListUsers)
			users.POST("/", controllers.CreateUser)
			users.GET("/:id", controllers.GetUser)
			users.PUT("/:id", controllers.UpdateUser)
			users.DELETE("/:id", controllers.DeleteUser)
		}

		categories := protected.Group("/categories")
		{
			categories.GET("/", controllers.ListCategories)                      // List with pagination, search, filters
			categories.GET("/simple", controllers.GetAllCategoriesSimple)        // Simple list for dropdowns
			categories.GET("/stats", controllers.GetCategoryStats)               // Statistics
			categories.GET("/:id", controllers.GetCategory)                      // Get single category
			categories.POST("/", controllers.CreateCategory)                     // Create category
			categories.PUT("/:id", controllers.UpdateCategory)                   // Update category
			categories.PATCH("/:id/toggle-status", controllers.ToggleCategoryStatus) // Toggle status
			categories.DELETE("/:id", controllers.DeleteCategory)                // Delete single
			categories.POST("/bulk-delete", controllers.BulkDeleteCategories)    // Bulk delete
		}

		attributes := protected.Group("/attributes")
		{
			attributes.GET("/", controllers.ListAttributes)                       // List with pagination, search, filters
			attributes.GET("/simple", controllers.GetAllAttributesSimple)         // Simple list for dropdowns
			attributes.GET("/stats", controllers.GetAttributeStats)               // Statistics
			attributes.GET("/:id", controllers.GetAttribute)                      // Get single attribute
			attributes.POST("/", controllers.CreateAttribute)                     // Create attribute
			attributes.PUT("/:id", controllers.UpdateAttribute)                   // Update attribute
			attributes.PATCH("/:id/toggle-status", controllers.ToggleAttributeStatus) // Toggle status
			attributes.DELETE("/:id", controllers.DeleteAttribute)                // Delete attribute
			attributes.POST("/bulk-delete", controllers.BulkDeleteAttributes)     // Bulk delete
		}

		coupons := protected.Group("/coupons")
		{
			coupons.GET("/", controllers.ListCoupons)                           // List all coupons
			coupons.POST("/", controllers.CreateCoupon)                         // Create coupon (JSON)
			coupons.POST("/with-image", controllers.CreateCouponWithImage)      // Create coupon (form-data with image)
			coupons.POST("/validate", controllers.ValidateCoupon)               // Validate coupon code
			coupons.GET("/by-code/:code", controllers.GetCouponByCode)          // Get coupon by code (public)
			coupons.GET("/:id", controllers.GetCoupon)                          // Get single coupon by ID
			coupons.GET("/:id/usage", controllers.GetCouponUsageStats)          // Get usage analytics
			coupons.PUT("/:id", controllers.UpdateCoupon)                       // Update coupon (JSON)
			coupons.PUT("/:id/with-image", controllers.UpdateCouponWithImage)   // Update coupon (form-data with image)
			coupons.DELETE("/:id", controllers.DeleteCoupon)                    // Delete coupon
		}

		locations := protected.Group("/locations")
		{
			locations.GET("/", controllers.ListLocations)
			locations.POST("/", controllers.CreateLocation)
			locations.GET("/:id", controllers.GetLocation)
			locations.PUT("/:id", controllers.UpdateLocation)
			locations.DELETE("/:id", controllers.DeleteLocation)
		}

		products := protected.Group("/products")
		{
			products.GET("/", controllers.ListProducts)
			products.POST("/", controllers.CreateProduct)
			products.GET("/:id", controllers.GetProduct)
			products.PUT("/:id", controllers.UpdateProduct)
			products.PATCH("/:id/status", controllers.UpdateProductStatus)
			products.DELETE("/:id", controllers.DeleteProduct)
		}

		customers := protected.Group("/customers")
		{
			customers.GET("/", controllers.ListCustomers)
			customers.POST("/", controllers.CreateCustomer)
			customers.GET("/:id", controllers.GetCustomer)
			customers.PUT("/:id", controllers.UpdateCustomer)
			customers.DELETE("/:id", controllers.DeleteCustomer)
		}

		vendors := protected.Group("/vendors")
		{
			vendors.GET("/", controllers.ListVendors)
			vendors.POST("/", controllers.CreateVendor)
			vendors.GET("/:id", controllers.GetVendor)
			vendors.PUT("/:id", controllers.UpdateVendor)
			vendors.DELETE("/:id", controllers.DeleteVendor)
		}

		staffRoutes := protected.Group("/staff")
		{
			staffRoutes.GET("/", controllers.ListStaff)
			staffRoutes.POST("/", controllers.CreateStaff)
			staffRoutes.GET("/:id", controllers.GetStaff)
			staffRoutes.PUT("/:id", controllers.UpdateStaff)
			staffRoutes.DELETE("/:id", controllers.DeleteStaff)
		}

		staffRoles := protected.Group("/staff-roles")
		{
			staffRoles.GET("/", controllers.ListStaffRoles)
			staffRoles.POST("/", controllers.CreateStaffRole)
			staffRoles.GET("/:id", controllers.GetStaffRole)
			staffRoles.PUT("/:id", controllers.UpdateStaffRole)
			staffRoles.DELETE("/:id", controllers.DeleteStaffRole)
		}

		salaryPayments := protected.Group("/salary-payments")
		{
			salaryPayments.GET("/", controllers.ListSalaryPayments)
			salaryPayments.POST("/", controllers.CreateSalaryPayment)
			salaryPayments.GET("/:id", controllers.GetSalaryPayment)
			salaryPayments.PUT("/:id", controllers.UpdateSalaryPayment)
			salaryPayments.DELETE("/:id", controllers.DeleteSalaryPayment)
		}

		inventory := protected.Group("/inventory")
		{
			inventory.GET("/", controllers.ListInventory)
		}

		transfers := protected.Group("/transfers")
		{
			transfers.GET("/", controllers.ListTransfers)
			transfers.GET("/products-by-location/:location_id", controllers.GetProductsByLocation)
			transfers.POST("/", controllers.CreateTransfer)
			transfers.PUT("/:id/cancel", controllers.CancelTransfer)
		}

		sells := protected.Group("/sells")
		{
			sells.GET("/", controllers.ListSells)
			sells.GET("/stats", controllers.GetSellsStats)
			sells.GET("/invoice/:invoiceNo", controllers.GetSellByInvoice)
			sells.POST("/", controllers.CreateSell)
			sells.GET("/:id", controllers.GetSell)
			sells.PUT("/:id", controllers.UpdateSell)
			sells.PATCH("/:id/status", controllers.UpdateSellStatus)
			sells.DELETE("/:id", controllers.DeleteSell)
		}

		customerReturns := protected.Group("/customer-returns")
		{
			customerReturns.GET("/", controllers.ListCustomerReturns)
			customerReturns.GET("/stats", controllers.GetCustomerReturnStats)
			customerReturns.POST("/", controllers.CreateCustomerReturn)
			customerReturns.GET("/customer/:customerId", controllers.GetCustomerReturnsByCustomer)
			customerReturns.GET("/:id", controllers.GetCustomerReturn)
			customerReturns.PUT("/:id", controllers.UpdateCustomerReturn)
			customerReturns.PATCH("/:id/approve", controllers.ApproveCustomerReturn)
			customerReturns.PATCH("/:id/reject", controllers.RejectCustomerReturn)
			customerReturns.DELETE("/:id", controllers.DeleteCustomerReturn)
		}

		vendorReturns := protected.Group("/vendor-returns")
		{
			vendorReturns.GET("/", controllers.ListVendorReturns)
			vendorReturns.GET("/stats", controllers.GetVendorReturnStats)
			vendorReturns.POST("/", controllers.CreateVendorReturn)
			vendorReturns.GET("/vendor/:vendorId", controllers.GetVendorReturnsByVendor)
			vendorReturns.GET("/:id", controllers.GetVendorReturn)
			vendorReturns.PUT("/:id", controllers.UpdateVendorReturn)
			vendorReturns.PATCH("/:id/status", controllers.UpdateVendorReturnStatus)
			vendorReturns.DELETE("/:id", controllers.DeleteVendorReturn)
		}

		shippingAddresses := protected.Group("/shipping-addresses")
		{
			shippingAddresses.GET("/", controllers.ListShippingAddresses)
			shippingAddresses.POST("/", controllers.CreateShippingAddress)
			shippingAddresses.GET("/:id", controllers.GetShippingAddress)
			shippingAddresses.PUT("/:id", controllers.UpdateShippingAddress)
			shippingAddresses.DELETE("/:id", controllers.DeleteShippingAddress)
			shippingAddresses.PATCH("/:id/set-default", controllers.SetDefaultShippingAddress)
		}

		shipments := protected.Group("/shipments")
		{
			shipments.GET("/", controllers.ListOrderShipments)
			shipments.GET("/stats", controllers.GetShipmentStats)
			shipments.POST("/", controllers.CreateOrderShipment)
			shipments.GET("/:id", controllers.GetOrderShipment)
			shipments.PATCH("/:id/status", controllers.UpdateShipmentStatus)
			shipments.POST("/:id/tracking", controllers.AddTrackingEvent)
		}

		settings := protected.Group("/settings")
		{
			settings.GET("/", controllers.GetSettings)
			settings.PUT("/general", controllers.UpdateGeneralSettings)
			settings.PATCH("/tax", controllers.UpdateTaxSettings)
			settings.PATCH("/shipping", controllers.UpdateShippingSettings)
			settings.PATCH("/payment", controllers.UpdatePaymentSettings)
			settings.PATCH("/business", controllers.UpdateBusinessSettings)
			settings.PATCH("/regional", controllers.UpdateRegionalSettings)
			settings.PATCH("/notifications", controllers.UpdateNotificationSettings)
			settings.PATCH("/store-hours", controllers.UpdateStoreHours)
			settings.POST("/upload-logo", controllers.UploadLogo)
			settings.POST("/upload-banner", controllers.UploadBanner)
		}

		// Protected SaaS Routes (require authentication)
		// Auth endpoints
		protected.GET("/auth/me", controllers.GetCurrentUser)
		protected.POST("/auth/logout", controllers.SaasLogout)

		// Company endpoints
		protected.GET("/company/profile", controllers.GetCompanyProfile)
		protected.PATCH("/company/profile", controllers.UpdateCompanyProfile)
		protected.GET("/company/status", controllers.GetCompanyStatus)
		protected.GET("/company/settings", controllers.GetCompanySettings)
		protected.PATCH("/company/settings", controllers.UpdateCompanySettings)

		// Billing Contact endpoints
		protected.GET("/company/billing-contact", controllers.GetBillingContact)
		protected.PATCH("/company/billing-contact", controllers.UpdateBillingContact)

		// Team/Users endpoints
		protected.GET("/company/team/users", controllers.GetTeamUsers)
		protected.POST("/company/team/users/invite", controllers.InviteUser)
		protected.PATCH("/company/team/users/:userId/role", controllers.UpdateUserRole)
		protected.DELETE("/company/team/users/:userId", controllers.RemoveUser)
		protected.POST("/company/team/users/resend-invitation/:userId", controllers.ResendInvitation)

		// Billing/Subscription endpoints
		protected.GET("/billing/plans", controllers.GetPlans)
		protected.GET("/billing/subscription/current", controllers.GetCurrentSubscription)
		protected.GET("/billing/payments/history", controllers.GetPaymentHistory)
		protected.POST("/billing/subscription/renew", controllers.RenewSubscription)
		protected.POST("/billing/subscription/cancel", controllers.CancelSubscription)
		protected.POST("/billing/subscription/upgrade", controllers.UpgradeSubscription)
		protected.POST("/billing/subscription/create", controllers.CreateSubscriptionForCompany) // For seeding/testing
	}

	return r
}
