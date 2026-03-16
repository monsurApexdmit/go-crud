package database

import (
    "log"
    "go-crud/models"
)

func Migrate() {
    err := DB.AutoMigrate(
        &models.User{},
        &models.Role{},
        &models.Permission{},
        &models.RolePermission{},
        &models.Author{},
        &models.Book{},
        &models.Category{},
        &models.Attribute{},
        &models.Product{},
        &models.ProductVariant{},
        &models.ProductImage{},
        &models.VariantInventory{},
        &models.Location{},
        &models.Customer{},
        &models.Vendor{},
        &models.Staff{},
        &models.StaffRole{},
        &models.SalaryPayment{},
        &models.Sell{},
        &models.OrderItem{},
        &models.ShippingAddress{},
        &models.OrderShipment{},
        &models.ShipmentTrackingHistory{},
        &models.CustomerReturn{},
        &models.CustomerReturnItem{},
        &models.VendorReturn{},
        &models.VendorReturnItem{},
        &models.StockTransfer{},
        &models.Coupon{},
        &models.CouponUsage{},
        &models.Settings{},
        // SaaS models
        &models.Company{},
        &models.SaasUser{},
        &models.SubscriptionPlan{},
        &models.Subscription{},
        &models.Payment{},
        &models.BillingContact{},
        &models.CompanySettings{},
        &models.Invitation{},
        &models.PasswordReset{},
    )

    if err != nil {
        log.Fatal("Migration failed:", err)
    }

    log.Println("✅ Database migrated successfully with all models including SaaS")
}
