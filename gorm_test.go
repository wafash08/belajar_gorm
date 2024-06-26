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
	"gorm.io/gorm/clause"
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

type UserResponse struct {
	ID        string
	FirstName string
	LastName  string
}

func TestQueryNonModel(t *testing.T) {
	var users []UserResponse
	// menyimpan hasil query model User ke data yang bertipe bukan model, dalam hal ini struct UserResponse
	err := db.Model(&User{}).Select("id", "first_name", "last_name").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 14, len(users))
	for _, user := range users {
		fmt.Println("user >> ", user)
	}
}

// Save mengubah secara keseluruhan
func TestUpdate(t *testing.T) {
	user := User{}
	err := db.Take(&user, "id = ?", "1").Error
	assert.Nil(t, err)

	// melakukan update
	user.Name.FirstName = "Budi"
	user.Name.MiddleName = ""
	user.Name.LastName = "Nugraha"
	user.Password = "rahasia123"

	// menyimpan hasil update
	// method Save akan mengubah semua kolom
	err = db.Save(&user).Error
	assert.Nil(t, err)
}

// Update/Updates mengubah secara parsial
func TestUpdateSelectedColumns(t *testing.T) {
	// Updates => mengubah beberapa kolom
	// jika menggunakan map maka "" (string kosong) akan dianggap sebagai perubahan juga
	err := db.Model(&User{}).Where("id = ?", "1").Updates(map[string]interface{}{
		"middle_name": "",
		"last_name":   "Morro",
	}).Error
	assert.Nil(t, err)

	// Update => mengubah satu kolom
	err = db.Model(&User{}).Where("id = ?", "1").Update("password", "diubahlagi").Error
	assert.Nil(t, err)

	// jika menggunakan struct maka "" (string kosong) tidak dianggap sebagai perubahan
	err = db.Where("id = ?", "1").Updates(User{
		Name: Name{
			FirstName: "Eko",
			LastName:  "Khannedy",
		},
	}).Error
	assert.Nil(t, err)
}

func TestAutoIncrement(t *testing.T) {
	for i := 0; i < 10; i++ {
		userLog := UserLog{
			UserId: "1",
			Action: "Test Action",
		}

		err := db.Create(&userLog).Error
		assert.Nil(t, err)

		assert.NotEqual(t, 0, userLog.ID)
		fmt.Println(userLog.ID)
	}
}

func TestSaveOrUpdate(t *testing.T) {
	userLog := UserLog{
		UserId: "1",
		Action: "Test Action",
	}

	// Save dapat berfungsi update dan create
	// berfungsi sebagai create jika data yang dikirim tidak memiliki ID
	// berfungsi sebagai update jika memiliki ID
	// Save sangat cocok untuk ID yang auto increment
	err := db.Save(&userLog).Error // insert or create
	assert.Nil(t, err)

	userLog.UserId = "2"
	err = db.Save(&userLog).Error // update
	assert.Nil(t, err)
}

func TestSaveOrUpdateNonAutoIncrement(t *testing.T) {
	user := User{
		ID: "99", // belum ada user dengan ID '99'
		Name: Name{
			FirstName: "User 99",
		},
	}

	err := db.Save(&user).Error // insert or create
	assert.Nil(t, err)

	user.Name.FirstName = "User 99 Updated"
	err = db.Save(&user).Error // update
	assert.Nil(t, err)
}

func TestConflict(t *testing.T) {
	user := User{
		ID: "88",
		Name: Name{
			FirstName: "User 88",
		},
	}

	// Clause digunakan untuk mengubah pengaturan konflik
	err := db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&user).Error // insert
	assert.Nil(t, err)
}

func TestDelete(t *testing.T) {
	var user User
	// ambil data terlebih dahulu
	err := db.Take(&user, "id = ?", "88").Error
	assert.Nil(t, err)
	// lalu hapus
	err = db.Delete(&user).Error
	assert.Nil(t, err)

	// hapus menggunakan inline condition
	err = db.Delete(&User{}, "id = ?", "99").Error
	assert.Nil(t, err)

	err = db.Where("id = ?", "77").Delete(&User{}).Error
	assert.Nil(t, err)
}

func TestSoftDelete(t *testing.T) {
	todo := Todo{
		UserId:      "1",
		Title:       "Todo 1",
		Description: "Description 1",
	}
	err := db.Create(&todo).Error
	assert.Nil(t, err)

	// UPDATE "todos" SET "deleted_at"='2024-05-18 14:15:32.083' WHERE "todos"."id" = 1 AND "todos"."deleted_at" IS NULL
	// tidak melakukan query DELETE
	// tapi hanya mengubah kolom deleted_at
	err = db.Delete(&todo).Error
	assert.Nil(t, err)
	assert.NotNil(t, todo.DeletedAt)

	// SELECT * FROM "todos" WHERE "todos"."deleted_at" IS NULL
	var todos []Todo
	err = db.Find(&todos).Error
	assert.Nil(t, err)
	assert.Equal(t, 0, len(todos))
}

func TestUnscoped(t *testing.T) {
	var todo Todo
	//  SELECT * FROM "todos" WHERE id = 1 ORDER BY "todos"."id" LIMIT 1
	// method Unscoped menjadikan query yang dilakukan
	// tidak menyertakan kondisi kolom deleted_at
	err := db.Unscoped().First(&todo, "id = ?", 1).Error
	assert.Nil(t, err)
	fmt.Println(todo)

	// DELETE FROM "todos" WHERE "todos"."id" = 1
	err = db.Unscoped().Delete(&todo).Error
	assert.Nil(t, err)

	// SELECT * FROM "todos"
	var todos []Todo
	err = db.Unscoped().Find(&todos).Error
	assert.Nil(t, err)
}

func TestLock(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		var user User
		// locking for UPDATE
		// SELECT * FROM "users" WHERE id = '1' LIMIT 1 FOR UPDATE
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&user, "id = ?", "1").Error
		if err != nil {
			return err
		}

		user.Name.FirstName = "Joko"
		user.Name.LastName = "Morro"
		err = tx.Save(&user).Error
		return err
	})
	assert.Nil(t, err)
}

func TestCreateWallet(t *testing.T) {
	wallet := Wallet{
		ID:      "1",
		UserId:  "1",
		Balance: 1000000,
	}

	err := db.Create(&wallet).Error
	assert.Nil(t, err)
}

func TestRetrieveRelation(t *testing.T) {
	var user User
	// method Preload => melakukan query terhadap relation
	// Preload akan menjalan dua query, pertama ke tabel users lalu kedua ke tabel wallets
	// SELECT * FROM "wallets" WHERE "wallets"."user_id" = '1'
	// SELECT * FROM "users" WHERE id = '1' LIMIT 1
	// dan ini sebenarnya cocok untuk relasi one to many dan relasi many to many
	err := db.Model(&User{}).Preload("Wallet").Take(&user, "id = ?", "1").Error
	assert.Nil(t, err)

	fmt.Println("user > ", user)
	assert.Equal(t, "1", user.ID)
	assert.Equal(t, "1", user.Wallet.ID)
}

func TestRetrieveRelationJoin(t *testing.T) {
	var user User
	// method Joins cocok untuk relasi one to one
	// Join hanya melakukan satu query
	//  SELECT "users"."id","users"."password","users"."first_name","users"."middle_name","users"."last_name","users"."created_at","users"."updated_at","Wallet"."id" AS "Wallet__id","Wallet"."user_id" AS "Wallet__user_id","Wallet"."balance" AS "Wallet__balance","Wallet"."created_at" AS "Wallet__created_at","Wallet"."updated_at" AS "Wallet__updated_at" FROM "users" LEFT JOIN "wallets" "Wallet" ON "users"."id" = "Wallet"."user_id" WHERE users.id = '1' LIMIT 1
	err := db.Model(&User{}).Joins("Wallet").Take(&user, "users.id = ?", "1").Error
	assert.Nil(t, err)

	fmt.Println("user > ", user)
	assert.Equal(t, "1", user.ID)
	assert.Equal(t, "1", user.Wallet.ID)
}

func TestAutoCreateUpdate(t *testing.T) {
	user := User{
		ID:       "20",
		Password: "rahasia",
		Name: Name{
			FirstName: "User 20",
		},
		Wallet: Wallet{
			ID:      "20",
			UserId:  "20",
			Balance: 1000000,
		},
	}

	// query yg dihasilkan
	// INSERT INTO "wallets" ("id","user_id","balance","created_at","updated_at") VALUES ('20','20',1000000,'2024-05-18 23:22:14.407','2024-05-18 23:22:14.407') ON CONFLICT ("id") DO UPDATE SET "user_id"="excluded"."user_id"
	// INSERT INTO "users" ("id","password","first_name","middle_name","last_name","created_at","updated_at") VALUES ('20','rahasia','User 20','','','2024-05-18 23:22:14.396','2024-05-18 23:22:14.396')
	err := db.Create(&user).Error
	assert.Nil(t, err)
}

func TestSkipAutoCreateUpdate(t *testing.T) {
	user := User{
		ID:       "21",
		Password: "rahasia",
		Name: Name{
			FirstName: "User 21",
		},
		Wallet: Wallet{
			ID:      "21",
			UserId:  "21",
			Balance: 1000000,
		},
	}

	// digunakan untuk mengabaikan pembuatan kolom relasi
	// INSERT INTO "users" ("id","password","first_name","middle_name","last_name","created_at","updated_at") VALUES ('21','rahasia','User 21','','','2024-05-19 08:16:31.115','2024-05-19 08:16:31.115')
	err := db.Omit(clause.Associations).Create(&user).Error
	assert.Nil(t, err)
}

func TestUserAndAddresses(t *testing.T) {
	user := User{
		ID:       "2",
		Password: "rahasia",
		Name: Name{
			FirstName: "User 50",
		},
		Wallet: Wallet{
			ID:      "2",
			UserId:  "2",
			Balance: 1000000,
		},
		Addresses: []Address{
			{
				UserId:  "2",
				Address: "Jalan A",
			},
			{
				UserId:  "2",
				Address: "Jalan B",
			},
		},
	}

	// INSERT INTO "wallets" ("id","user_id","balance","created_at","updated_at") VALUES ('2','2',1000000,'2024-05-19 08:39:09.697','2024-05-19 08:39:09.697') ON CONFLICT ("id") DO UPDATE SET "user_id"="excluded"."user_id"
	// INSERT INTO "addresses" ("user_id","address","created_at","updated_at") VALUES ('2','Jalan A','2024-05-19 08:39:09.711','2024-05-19 08:39:09.711'),('2','Jalan B','2024-05-19 08:39:09.711','2024-05-19 08:39:09.711') ON CONFLICT ("id") DO UPDATE SET "user_id"="excluded"."user_id" RETURNING "id"
	// UPDATE "users" SET "password"='rahasia',"first_name"='User 50',"middle_name"='',"last_name"='',"updated_at"='2024-05-19 08:39:09.682' WHERE "id" = '2'
	err := db.Save(&user).Error
	assert.Nil(t, err)
}

func TestPreloadJoinOneToMany(t *testing.T) {
	var users []User
	// SELECT * FROM "addresses" WHERE "addresses"."user_id" IN ('1','20','2','10','7','6','8','4','11','99','12','3','5','21','14','13','9')
	// SELECT "users"."id","users"."password","users"."first_name","users"."middle_name","users"."last_name","users"."created_at","users"."updated_at","Wallet"."id" AS "Wallet__id","Wallet"."user_id" AS "Wallet__user_id","Wallet"."balance" AS "Wallet__balance","Wallet"."created_at" AS "Wallet__created_at","Wallet"."updated_at" AS "Wallet__updated_at" FROM "users" LEFT JOIN "wallets" "Wallet" ON "users"."id" = "Wallet"."user_id"
	err := db.Model(&User{}).Preload("Addresses").Joins("Wallet").Find(&users).Error
	assert.Nil(t, err)
}

func TestTakePreloadJoinOneToMany(t *testing.T) {
	var user User
	// di sini Preload digunakan untuk relasi one to many
	// dan Joins digunakan untuk relasi one to one
	err := db.Model(&User{}).Preload("Addresses").Joins("Wallet").
		Take(&user, "users.id = ?", "2").Error
	assert.Nil(t, err)
}

func TestBelongsTo(t *testing.T) {
	fmt.Println("Preload")
	var addresses []Address
	// SELECT * FROM "users" WHERE "users"."id" = '2'
	err := db.Model(&Address{}).Preload("User").Find(&addresses).Error
	assert.Nil(t, err)
	// assert.Equal(t, 4, len(addresses))

	fmt.Println("Joins")
	addresses = []Address{}
	// SELECT "addresses"."id","addresses"."user_id","addresses"."address","addresses"."created_at","addresses"."updated_at","User"."id" AS "User__id","User"."password" AS "User__password","User"."first_name" AS "User__first_name","User"."middle_name" AS "User__middle_name","User"."last_name" AS "User__last_name","User"."created_at" AS "User__created_at","User"."updated_at" AS "User__updated_at" FROM "addresses" LEFT JOIN "users" "User" ON "addresses"."user_id" = "User"."id"
	err = db.Model(&Address{}).Joins("User").Find(&addresses).Error
	assert.Nil(t, err)
	// assert.Equal(t, 4, len(addresses))
}

func TestBelongsToWallet(t *testing.T) {
	fmt.Println("Preload")
	var wallets []Wallet
	// SELECT * FROM "users" WHERE "users"."id" IN ('1','20','2')
	// SELECT * FROM "wallets"
	err := db.Model(&Wallet{}).Preload("User").Find(&wallets).Error
	assert.Nil(t, err)

	fmt.Println("Joins")
	wallets = []Wallet{}
	// SELECT "wallets"."id","wallets"."user_id","wallets"."balance","wallets"."created_at","wallets"."updated_at","User"."id" AS "User__id","User"."password" AS "User__password","User"."first_name" AS "User__first_name","User"."middle_name" AS "User__middle_name","User"."last_name" AS "User__last_name","User"."created_at" AS "User__created_at","User"."updated_at" AS "User__updated_at" FROM "wallets" LEFT JOIN "users" "User" ON "wallets"."user_id" = "User"."id"
	err = db.Model(&Wallet{}).Joins("User").Find(&wallets).Error
	assert.Nil(t, err)
}
