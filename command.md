
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
