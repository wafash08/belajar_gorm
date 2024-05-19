package belajargorm

import "time"

type Wallet struct {
	// field UserId dijadikan sebagai foreign key yang merujuk pada kolom id di tabel users
	ID        string    `gorm:"primary_key;column:id"`
	UserId    string    `gorm:"column:user_id"`
	Balance   int64     `gorm:"column:balance"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	User      *User     `gorm:"foreignKey:user_id;references:id"`
	// gunakan pointer untuk menghindari cylic
}

func (w *Wallet) TableName() string {
	return "wallets"
}
