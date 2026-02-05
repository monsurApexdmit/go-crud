
migrate create -ext sql -dir database/migrations -seq modify_users_email

migrate \
  -database "mysql://admin:admin123@tcp(localhost:3306)/go" \
  -path database/migrations \
  up

migrate \
  -database "mysql://admin:admin123@tcp(localhost:3306)/go" \
  -path database/migrations \
  version

migrate \
  -database "mysql://admin:admin123@tcp(central_mysql:3306)/go" \
  -path database/migrations \
  down 1

in local mysql
migrate -database 'mysql://gocrud:StrongPass123!@tcp(localhost:3306)/go' \
  -path database/migrations up



Staff Roles

GET /staff-roles/ — no params

POST /staff-roles/


{
  "name": "string",
  "permissions": [
    { "name": "Dashboard", "read": true, "write": false, "delete": false },
    { "name": "Products", "read": true, "write": true, "delete": false },
    { "name": "Categories", "read": true, "write": true, "delete": true },
    { "name": "Attributes", "read": false, "write": false, "delete": false },
    { "name": "Coupons", "read": false, "write": false, "delete": false },
    { "name": "Customers", "read": false, "write": false, "delete": false },
    { "name": "Orders", "read": false, "write": false, "delete": false },
    { "name": "POS", "read": false, "write": false, "delete": false },
    { "name": "Sells", "read": false, "write": false, "delete": false },
    { "name": "Staff", "read": false, "write": false, "delete": false },
    { "name": "Settings", "read": false, "write": false, "delete": false },
    { "name": "International", "read": false, "write": false, "delete": false },
    { "name": "Store", "read": false, "write": false, "delete": false },
    { "name": "Pages", "read": false, "write": false, "delete": false }
  ]
}
GET /staff-roles/:id — no body

PUT /staff-roles/:id — same fields as POST (both name and permissions are optional)

DELETE /staff-roles/:id — no body

Salary Payments

GET /salary-payments/ — query params: ?staff_id=1&month=2026-01&status=Pending|Partial|Paid|all&page=1&limit=10

POST /salary-payments/


{
  "staffId": 1,
  "month": "2026-01",
  "amount": 5000.00,
  "paidAmount": 0,
  "paymentDate": "2026-01-15",
  "paymentMethod": "string",
  "notes": "string"
}
(Status auto-calculates: Pending / Partial / Paid)

GET /salary-payments/:id — no body

PUT /salary-payments/:id


{
  "staffId": 1,
  "month": "2026-01",
  "amount": 5000.00,
  "paidAmount": 3000.00,
  "paymentDate": "2026-01-20",
  "paymentMethod": "Bank Transfer",
  "notes": "Partial payment"
}
(Status recalculates automatically based on amount vs paidAmount)

DELETE /salary-payments/:id — no body