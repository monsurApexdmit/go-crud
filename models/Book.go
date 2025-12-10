package models
import "time"
type Book struct {
    ID     uint   `gorm:"primaryKey" json:"id"`
    Title  string `json:"title"`
    AuthorID  uint      `json:"author_id"`           // Foreign key column
    Author    Author    `json:"author" gorm:"foreignKey:AuthorID"` 
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
