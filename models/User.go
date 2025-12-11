package models
import "time"
type User struct {
    ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username"`
	Email    string `json:"email"`
    Password  string    `json:"-" gorm:"column:password"` // json:"-" prevents password from being returned
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}