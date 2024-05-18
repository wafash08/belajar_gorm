package belajargorm

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	// berisi
	// ID        uint `gorm:"primarykey"`
	// CreatedAt time.Time
	// UpdatedAt time.Time
	// DeletedAt DeletedAt `gorm:"index"`
	// dan cocok digunakan jika field struct sesuai dengan model convention GORM
	UserId      string `gorm:"column:user_id"`
	Title       string `gorm:"column:title"`
	Description string `gorm:"column:description"`
}

func (t *Todo) TableName() string {
	return "todos"
}
