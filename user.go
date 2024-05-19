package belajargorm

import "time"

type User struct {
	ID          string    `gorm:"primaryKey;column:id;<-:create"`
	Password    string    `gorm:"column:password"`
	Name        Name      `gorm:"embedded"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime;<-:create"`      // tidak perlu ditambahkan autoCreateTime
	UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"` // tidak perlu ditambahkan autoCreateTime dan autoUpdateTime
	Information string    `gorm:"-"`
	Wallet      Wallet    `gorm:"foreignKey:user_id;references:id"`
	// foreignKey merujuk kolom yang dijadikan sebagai foreign key yang ada di tabel wallet
	// references merujuk pada kolom id pada tabel saat ini yaitu tabel user
}

// secara default
// gorm akan memilih nama tabel dari nama struct menggunakan lower_case dan jamak
// dan nama kolom dari nama field menggunakan snake_case
// dan memilih field ID sebagai primary key
// akan tetapi disarankan untuk menentukan nama tabel, kolom, dan id secara menggunakan field tags gorm
// lihat https://gorm.io/docs/models.html#Fields-Tags

// mengubah nama table mapping
func (u *User) TableName() string {
	return "users"
}

// field permission
// <- = write permission, create and update
// <-:create = create only
// <-:update = update only
// -> = read permisson
// ->:false = no read permission
// - = ignoring the field, no write/read permission

// embedded struct
// dilakukan dengan menambahkan gorm:"embedded"
// contoh
type Name struct {
	FirstName  string `gorm:"column:first_name"`
	MiddleName string `gorm:"column:middle_name"`
	LastName   string `gorm:"column:last_name"`
}

// lalu sisipkan ke field

type UserLog struct {
	ID        int       `gorm:"primaryKey;column:id;autoIncrement"`
	UserId    string    `gorm:"column:user_id"`
	Action    string    `gorm:"column:action"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (l *UserLog) TableName() string {
	return "user_logs"
}
