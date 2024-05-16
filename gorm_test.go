package belajargorm

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenConnection() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	url := os.Getenv("DB")
	dialect := postgres.Open(url)
	db, err := gorm.Open(dialect, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	return db
}

var db = OpenConnection()

func TestOpenConnection(t *testing.T) {
	assert.NotNil(t, db)
}

// Exec digunakan untuk memanipulasi data
func TestExecuteSQL(t *testing.T) {
	err := db.Exec("INSERT INTO sample (id, name) VALUES (?, ?)", "1", "Eko").Error
	assert.Nil(t, err)

	err = db.Exec("INSERT INTO sample (id, name) VALUES (?, ?)", "2", "Budi").Error
	assert.Nil(t, err)
}

type Sample struct {
	Id   string
	Name string
}

// Raw digunakan untuk melakukan query
func TestRawSQL(t *testing.T) {
	var sample Sample
	err := db.Raw("select id, name from sample where id = ?", "1").Scan(&sample).Error
	assert.Nil(t, err)
	assert.Equal(t, "Eko", sample.Name)

	var samples []Sample
	err = db.Raw("select id, name from sample").Scan(&samples).Error
	assert.Nil(t, err)
	assert.Equal(t, 2, len(samples))
}

func TestSQLRow(t *testing.T) {
	// Rows digunakan untuk mendapatkan hasil sebagai *sql.Rows
	rows, err := db.Raw("select id, name from sample").Rows()
	assert.Nil(t, err)
	defer rows.Close()

	var samples []Sample
	for rows.Next() {
		var id string
		var name string

		err := rows.Scan(&id, &name)
		assert.Nil(t, err)

		fmt.Println("id > ", id)
		fmt.Println("name > ", name)

		samples = append(samples, Sample{
			Id:   id,
			Name: name,
		})
	}
	assert.Equal(t, 2, len(samples))
}

func TestScanRow(t *testing.T) {
	// Rows digunakan untuk mendapatkan hasil sebagai *sql.Rows
	rows, err := db.Raw("select id, name from sample").Rows()
	assert.Nil(t, err)
	defer rows.Close()

	var samples []Sample
	for rows.Next() {
		// ScanRows to scan a row into a struct
		err := db.ScanRows(rows, &samples)
		assert.Nil(t, err)
	}
	fmt.Println("samples >> ", samples)
	assert.Equal(t, 2, len(samples))
}

// Create => memasukkan/membuat data ke database satu data satu query
func TestCreateUser(t *testing.T) {
	user := User{
		ID:       "1",
		Password: "rahasia",
		Name: Name{
			FirstName:  "Eko",
			MiddleName: "Kurniawan",
			LastName:   "Khannedy",
		},
		Information: "Ini akan diignore",
	}

	// parameter harus berupa pointer
	result := db.Create(&user)
	assert.Nil(t, result.Error)
	assert.Equal(t, int64(1), result.RowsAffected)
}

// Create(slices) => memasukkan banyak data
// CreateInBatches(slices, sizes) => memasukkan banyak data sekaligus artinya banyak data satu query
func TestBatchInsert(t *testing.T) {
	var users []User
	// for i := 2; i < 10; i++ {
	// 	users = append(users, User{
	// 		ID:       strconv.Itoa(i),
	// 		Password: "rahasia",
	// 		Name: Name{
	// 			FirstName: "User " + strconv.Itoa(i),
	// 		},
	// 	})
	// }

	for i := 10; i < 20; i++ {
		users = append(users, User{
			ID:       strconv.Itoa(i),
			Password: "rahasia",
			Name: Name{
				FirstName: "User " + strconv.Itoa(i),
			},
		})
	}

	// result := db.Create(&users)
	result := db.CreateInBatches(&users, 2)
	assert.Nil(t, result.Error)
	// assert.Equal(t, 8, int(result.RowsAffected))
	assert.Equal(t, 10, int(result.RowsAffected))
}

// transaction
// hanya bisa terjadi kalau menggunakan koneksi database yang sama
// bisa digunakan menggunakan method Transaction
func TestTransactionSuccess(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&User{ID: "10", Password: "rahasia", Name: Name{FirstName: "User 10"}}).Error
		if err != nil {
			return err
		}

		err = tx.Create(&User{ID: "11", Password: "rahasia", Name: Name{FirstName: "User 11"}}).Error
		if err != nil {
			return err
		}

		err = tx.Create(&User{ID: "12", Password: "rahasia", Name: Name{FirstName: "User 12"}}).Error
		if err != nil {
			return err
		}

		return nil
	})

	assert.Nil(t, err)
}

func TestTransactionRollback(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&User{ID: "13", Password: "rahasia", Name: Name{FirstName: "User 13"}}).Error
		if err != nil {
			return err
		}

		// ini error dan akan menyebabkan rollback
		err = tx.Create(&User{ID: "11", Password: "rahasia", Name: Name{FirstName: "User 11"}}).Error
		if err != nil {
			return err
		}

		return nil
	})

	assert.NotNil(t, err)
}

// manual transaction
// ini tidak direkomendasikan
func TestManualTransactionSuccess(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	err := tx.Create(&User{ID: "13", Password: "rahasia", Name: Name{FirstName: "User 13"}}).Error
	assert.Nil(t, err)

	err = tx.Create(&User{ID: "14", Password: "rahasia", Name: Name{FirstName: "User 14"}}).Error
	assert.Nil(t, err)

	if err == nil {
		tx.Commit()
	}
}

func TestManualTransactionRollback(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	err := tx.Create(&User{ID: "15", Password: "rahasia", Name: Name{FirstName: "User 15"}}).Error
	assert.Nil(t, err)

	err = tx.Create(&User{ID: "14", Password: "rahasia", Name: Name{FirstName: "User 14"}}).Error
	assert.Nil(t, err)

	if err == nil {
		tx.Commit()
	}
}

// query that returns single object
func TestQuerySingleObject(t *testing.T) {
	user := User{}
	// First => mereturn single dalam keadaan terurut berdasarkan id
	err := db.First(&user).Error
	assert.Nil(t, err)
	assert.Equal(t, "1", user.ID)

	user = User{}
	// Last => mereturn single dalam keadaan terurut berdasarkan id
	err = db.Last(&user).Error
	assert.Nil(t, err)
	assert.Equal(t, "9", user.ID)
}

// inline condition
// akan otomatis menjadi kondisi where di sql selectnya
func TestQuerySingleObjectInlineCondition(t *testing.T) {
	user := User{}
	// inline condition
	err := db.Take(&user, "id = ?", "5").Error
	assert.Nil(t, err)
	assert.Equal(t, "5", user.ID)
	assert.Equal(t, "User 5", user.Name.FirstName)
}

func TestQueryAllObjects(t *testing.T) {
	var users []User
	// inline parameter dapat berupa slice
	err := db.Find(&users, "id in ?", []string{"1", "2", "3", "4"}).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(users))
}

func TestQueryCondition(t *testing.T) {
	var users []User
	// Where digunakan sebelum Find
	// ketika menggunakan where, maka query akan dianggap menggunakan operator AND SQL
	err := db.Where("first_name like ?", "%User%").Where("password = ?", "rahasia").Find(&users).Error
	assert.Nil(t, err)
	for _, user := range users {
		fmt.Println("user >> ", user.Name.FirstName)
	}
	assert.Equal(t, 13, len(users))
}

func TestOrOperator(t *testing.T) {
	var users []User
	// operator OR SQL dari method Or
	err := db.Where("first_name like ?", "%User%").Or("password = ?", "rahasia").Find(&users).Error
	assert.Nil(t, err)
	for _, user := range users {
		fmt.Println("user >> ", user.Name.FirstName)
	}
	assert.Equal(t, 14, len(users))
}

func TestNotOperator(t *testing.T) {
	var users []User
	// operator NOT SQL dari method Not
	// SELECT * FROM "users" WHERE NOT first_name like '%User%' AND password = 'rahasia'
	err := db.Not("first_name like ?", "%User%").Where("password = ?", "rahasia").Find(&users).Error
	assert.Nil(t, err)
	for _, user := range users {
		fmt.Println("user >> ", user.Name.FirstName)
	}
	assert.Equal(t, 1, len(users))
}

func TestSelectFields(t *testing.T) {
	var users []User
	// method Select digunakan untuk menentukan kolom apa saja yang akan dibaca
	// SELECT "id","first_name" FROM "users"
	err := db.Select("id", "first_name").Find(&users).Error
	assert.Nil(t, err)

	for _, user := range users {
		assert.NotNil(t, user.ID)
		assert.NotEqual(t, "", user.Name.FirstName)
	}

	assert.Equal(t, 14, len(users))
}

func TestStructCondition(t *testing.T) {
	userCondition := User{
		Name: Name{
			FirstName: "User 5",
			LastName:  "", // tidak bisa, karena dianggap default value
		},
		Password: "rahasia",
	}

	var users []User
	// field atau key akan menjadi nama kolom
	// value struct akan menjadi value query
	// SELECT * FROM "users" WHERE "users"."password" = 'rahasia' AND "users"."first_name" = 'User 5'
	err := db.Where(userCondition).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
}

func TestMapCondition(t *testing.T) {
	mapCondition := map[string]interface{}{
		"middle_name": "", // meskipun berisi string kosong tetap dianggap sebagai nilai query
		"last_name":   "", // meskipun berisi string kosong tetap dianggap sebagai nilai query
	}

	var users []User
	//  SELECT * FROM "users" WHERE "last_name" = '' AND "middle_name" = ''
	err := db.Where(mapCondition).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 13, len(users))
}

func TestOrderLimitOffset(t *testing.T) {
	var users []User
	// Order => untuk melakukan sorting
	// Limit dan Offset => untuk melakukan paging
	// SELECT * FROM "users" ORDER BY id asc, first_name desc LIMIT 5 OFFSET 5
	err := db.Order("id asc, first_name desc").Limit(5).Offset(5).Find(&users).Error
	assert.Nil(t, err)
	for _, user := range users {
		fmt.Println("user >> ", user.Name.FirstName)
	}
	assert.Equal(t, 5, len(users))
}