package belajargorm

import (
	"context"
	"log"
	"strconv"
	"testing"

	faker "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

func OpenConnection() *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=belajar_gorm port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
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

func TestRawSQL(t *testing.T) {
	err := db.Exec("INSERT INTO public.sample(id, name) values (?, ?)", "1", "Nathan").Error
	assert.Nil(t, err)

	err = db.Exec("INSERT INTO public.sample(id, name) values (?, ?)", "2", "Jono").Error
	assert.Nil(t, err)
}

type Sample struct {
	ID   string
	Name string
}

func TestQuerySQL(t *testing.T) {
	var sample Sample
	err := db.Raw("SELECT id, name FROM sample WHERE id = ?", "1").Scan(&sample).Error

	assert.Nil(t, err)
	assert.Equal(t, "Nathan", sample.Name)

	var samples []Sample
	err = db.Raw("SELECT id, name FROM sample").Scan(&samples).Error
	assert.Nil(t, err)
}

func TestSQLRow(t *testing.T) {
	rows, err := db.Raw("SELECT id, name FROM sample").Rows()
	assert.Nil(t, err)
	defer rows.Close()

	var samples []Sample
	for rows.Next() {
		var id, name string
		err := rows.Scan(&id, &name)
		assert.Nil(t, err)

		samples = append(samples, Sample{
			ID:   id,
			Name: name,
		})
	}
	assert.NotNil(t, samples)
}

func TestScanRows(t *testing.T) {
	rows, err := db.Raw("SELECT id, name FROM sample").Rows()
	assert.Nil(t, err)
	defer rows.Close()

	var samples []Sample
	for rows.Next() {
		err = db.ScanRows(rows, &samples)
		assert.Nil(t, err)
	}
}

func TestInsertData(t *testing.T) {
	user := User{
		ID:       faker.UUID(),
		Password: "password",
		Name: Name{
			FirstName:  "Yonathan",
			MiddleName: "",
			LastName:   "Setiadi",
		},
		Information: "info",
	}

	resp := db.Create(&user)
	assert.Nil(t, resp.Error)
	assert.Equal(t, int64(1), resp.RowsAffected)
}

func TestBatchInsertData(t *testing.T) {
	var users []User
	for i := 2; i < 11; i++ {
		user := User{
			ID:       strconv.Itoa(i),
			Password: "password",
			Name: Name{
				FirstName: "User " + strconv.Itoa(i),
			},
		}
		users = append(users, user)
	}

	resp := db.Create(&users)
	assert.Nil(t, resp.Error)
	assert.Equal(t, int64(9), resp.RowsAffected)
}

func TestTransactionSuccess(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&User{
			ID:       "11",
			Password: "password",
			Name: Name{
				FirstName: "User 11",
			},
		}).Error
		if err != nil {
			return err
		}

		err = tx.Create(&User{
			ID:       "12",
			Password: "password",
			Name: Name{
				FirstName: "User 12",
			},
		}).Error
		if err != nil {
			return err
		}

		return nil
	})
	assert.Nil(t, err)
}

func TestTransactionError(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&User{
			ID:       "13",
			Password: "password",
			Name: Name{
				FirstName: "User 13",
			},
		}).Error
		if err != nil {
			return err
		}

		err = tx.Create(&User{
			ID:       "12",
			Password: "password",
			Name: Name{
				FirstName: "User 12",
			},
		}).Error
		if err != nil {
			return err
		}

		return nil
	})
	assert.NotNil(t, err)
}

func TestManualTransactionSuccess(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	err := tx.Create(&User{
		ID:       "13",
		Password: "password",
		Name: Name{
			FirstName: "User 13",
		},
	}).Error
	assert.Nil(t, err)

	err = tx.Create(&User{
		ID:       "14",
		Password: "password",
		Name: Name{
			FirstName: "User 14",
		},
	}).Error
	assert.Nil(t, err)

	if err == nil {
		tx.Commit()
	}
}

func TestManualTransactionFailed(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	err := tx.Create(&User{
		ID:       "15",
		Password: "password",
		Name: Name{
			FirstName: "User 15",
		},
	}).Error
	assert.Nil(t, err)

	err = tx.Create(&User{
		ID:       "14",
		Password: "password",
		Name: Name{
			FirstName: "User 14",
		},
	}).Error
	assert.NotNil(t, err)

	if err == nil {
		tx.Commit()
	}
}

func TestSingleData(t *testing.T) {
	user := User{}
	err := db.First(&user, "password = ?", "rahasia").Error
	assert.Nil(t, err)
	assert.Equal(t, "6", user.ID)

	user = User{}
	err = db.Last(&user).Error
	assert.Nil(t, err)
	assert.Equal(t, "9", user.ID)

	user = User{}
	err = db.Take(&user, "password = ?", "rahasia").Error
	assert.Nil(t, err)
	assert.Equal(t, "6", user.ID)
}

func TestQueryAll(t *testing.T) {
	var users []User
	err := db.Find(&users, "id in ?", []string{"1", "3", "5"}).Error
	assert.Nil(t, err)
	assert.Equal(t, 3, len(users))
}

func TestQueryCondition(t *testing.T) {
	var users []User
	res := db.Where("first_name LIKE ?", "%User%").
		Where("password = ?", "password").
		Find(&users)
	assert.Nil(t, res.Error)
	assert.Equal(t, 12, res.RowsAffected)
}

func TestOrOperator(t *testing.T) {
	var users []User
	res := db.Where("first_name LIKE ?", "%User%").
		Or("password = ?", "password").
		Find(&users)
	assert.Nil(t, res.Error)
	assert.Equal(t, 14, len(users))
}

func TestNotOperator(t *testing.T) {
	var users []User
	res := db.Not("first_name LIKE ?", "%User%").
		Where("password = ?", "password").
		Find(&users)
	assert.Nil(t, res.Error)
	assert.Equal(t, 1, len(users))
}

func TestSelectOperator(t *testing.T) {
	var users []User
	err := db.Select("id", "first_name").Find(&users).Error
	assert.Nil(t, err)

	for _, u := range users {
		assert.NotNil(t, u.ID)
		assert.NotEqual(t, "", u.Name.FirstName)
	}

	assert.Equal(t, 14, len(users))
}

func TestStructCondition(t *testing.T) {
	userCondition := User{
		Name: Name{
			FirstName:  "User 5",
			MiddleName: "", // tidak bisa digunakan karena dianggap default value
		},
	}

	var users []User
	err := db.Where(userCondition).Find(&users).Error
	assert.Nil(t, err)
}

func TestMapCondition(t *testing.T) {
	mapCondition := map[string]interface{}{
		"last_name": "",
	}

	var users []User
	err := db.Where(mapCondition).Find(&users).Error
	assert.Nil(t, err)
}

func TestOrderLimitOffset(t *testing.T) {
	var users []User

	err := db.Order("id asc, first_name asc").Limit(5).Offset(5).Find(&users).Error
	assert.Nil(t, err)
}

type UserResponse struct {
	ID        string
	FirstName string
	LastName  string
}

func TestQueryNonModel(t *testing.T) {
	var users []UserResponse
	err := db.Model(&User{}).Select("id", "first_name", "last_name").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 14, len(users))
}

func TestUpdate(t *testing.T) {
	var user User
	err := db.Take(&user, "id = ?", "2").Error
	assert.Nil(t, err)

	user.Name.FirstName = "Ridwan"
	user.Name.MiddleName = "Mahendra"
	user.Name.LastName = "Hanif"
	user.Password = "password1234"

	err = db.Save(&user).Error
	assert.Nil(t, err)
}

func TestSelectedColumn(t *testing.T) {
	err := db.Model(&User{}).Where("id = ?", "3").Updates(map[string]interface{}{
		"first_name": "Merah",
		"last_name":  "Putih",
	}).Error
	assert.Nil(t, err)

	err = db.Model(&User{}).Where("id = ?", "4").Update("middle_name", "Ujang").Error
	assert.Nil(t, err)

	err = db.Where("id = ?", "5").Updates(User{
		Name: Name{
			FirstName: "Maria",
			LastName:  "Bellen",
		},
	}).Error
	assert.Nil(t, err)

}

func TestAutoIncrement(t *testing.T) {
	for i := 0; i < 10; i++ {
		userLog := UserLog{
			UserID: "1",
			Action: "Test Action",
		}

		err := db.Create(&userLog).Error
		assert.Nil(t, err)

		assert.NotEqual(t, 0, userLog.ID)
		t.Log(userLog.ID)
	}
}

func TestSaveorUpdate(t *testing.T) {
	userLog := UserLog{
		UserID: "1",
		Action: "Test Action",
	}

	/**
	*	Skenario pada method save untuk yang id auto increment
	*	- melakukan insert terlebih dahulu jika tidak ada parameter id yang dimasukkan
	*	- kalau ada parameter ID yang dimasukkan baru melakukan update
	 */

	result := db.Save(&userLog)
	assert.Nil(t, result.Error)

	userLog.UserID = "2"
	result = db.Save(&userLog)
	assert.Nil(t, result.Error)
}

func TestSaveOrUpdateNonAutoIncrement(t *testing.T) {
	user := User{
		ID:       "100",
		Password: "ahoep;83",
		Name: Name{
			FirstName: "Bina",
		},
	}

	/**
	*	Skenario pada method save untuk yang id non auto increment
	*	- melakukan update terlebih dahulu, lalu cek ada row yang kena efek atau engga
	*	- kalau ada row yang keupdate, maka berhenti di query update
	*	- kalau tidak ada baru melakukan insert
	*
	*	[rows:0] UPDATE "users" SET "password"='ahoep;83',"first_name"='Bina',
	*	"middle_name"='',"last_name"='',"updated_at"='2024-06-05 13:09:03.771' WHERE "id" = '100'
	*
	*	[rows:1] INSERT INTO "users" ("id","password","first_name","middle_name","last_name","created_at","updated_at")
	*	VALUES ('100','ahoep;83','Bina','','','2024-06-05 13:09:03.772','2024-06-05 13:09:03.771')
	*	ON CONFLICT ("id") DO UPDATE SET "updated_at"='2024-06-05 13:09:03.772',"password"="excluded"."password",
	*	"first_name"="excluded"."first_name","middle_name"="excluded"."middle_name","last_name"="excluded"."last_name"
	*
	*	[rows:1] UPDATE "users" SET "password"='ahoep;83',"first_name"='Bina',"middle_name"='Ujang',
	*	"last_name"='',"updated_at"='2024-06-05 13:09:03.779' WHERE "id" = '100'
	 */

	result := db.Save(&user)
	assert.Nil(t, result.Error)

	user.Name.MiddleName = "Ujang"

	result = db.Save(&user)
	assert.Nil(t, result.Error)
}

func TestConflict(t *testing.T) {
	user := User{
		ID: "88",
		Name: Name{
			FirstName: "User 88",
		},
	}

	res := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&user)
	assert.Nil(t, res.Error)
}

func TestDelete(t *testing.T) {
	var user User

	res := db.Take(&user, "id = ?", "11")
	assert.Nil(t, res.Error)

	res = db.Delete(&user)
	assert.Nil(t, res.Error)

	res = db.Delete(&User{}, "id = ?", "99")
	assert.Nil(t, res.Error)

	res = db.Where("id = ?", "40").Delete(&User{})
	assert.Nil(t, res.Error)
}

func TestSoftDelete(t *testing.T) {
	todo := Todo{
		UserID: "1",
		Task:   "Test Soft Delete",
	}

	res := db.Create(&todo)
	assert.Nil(t, res.Error)

	res = db.Delete(&todo)
	assert.Nil(t, res.Error)
	assert.NotNil(t, todo.DeletedAt)

	var todos []Todo

	res = db.Or("deleted_at = ?", nil).Find(&todos)
	assert.Nil(t, res.Error)
	assert.Equal(t, 0, len(todos))
}

func TestUnscope(t *testing.T) {
	var todo TodoGorm

	create := TodoGorm{
		UserID: "1",
		Task:   "Test Soft Delete",
	}
	res := db.Create(&create)
	assert.Nil(t, res.Error)

	err := db.Unscoped().First(&todo, "id = ?", "1").Error
	assert.Nil(t, err)

	err = db.Unscoped().Delete(&todo).Error
	assert.Nil(t, err)

	var todos []Todo
	err = db.Unscoped().Find(&todos).Error
	assert.Nil(t, err)
}

func TestLock(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		var user User
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&user, "id = ?", "1").Error
		if err != nil {
			return err
		}

		user.Name.FirstName = "Nathan"

		return tx.Save(user).Error
	})

	assert.Nil(t, err)
}

func TestCreateWallet(t *testing.T) {
	wallet := Wallet{
		UserID:  "1",
		Balance: 1000000,
	}

	err := db.Create(&wallet).Error
	assert.Nil(t, err)

}

func TestEagerLoad(t *testing.T) {
	var user User
	err := db.Model(&User{}).Preload("Wallet").Take(&user, "id = ?", "1").Error
	assert.Nil(t, err)

	assert.Equal(t, "1", user.ID)
	assert.Equal(t, uint(1), user.Wallet.ID)
}

func TestRetriveJoin(t *testing.T) {
	var users []User
	err := db.Model(&User{}).Joins("Wallet").Find(&users).Error
	assert.Nil(t, err)

	assert.Equal(t, 12, len(users))

}

func TestAutoUpsert(t *testing.T) {
	var user = User{
		ID:       "39",
		Password: "rahasia",
		Name: Name{
			FirstName: "Ujang",
			LastName:  "Pasaran",
		},
		Wallet: Wallet{
			UserID:  "39",
			Balance: 9307390,
		},
	}

	err := db.Create(&user).Error
	assert.Nil(t, err)
}

func TestSkipUpsert(t *testing.T) {
	var user = User{
		ID:       "40",
		Password: "rahasia",
		Name: Name{
			FirstName: "Pren",
			LastName:  "Pasaran",
		},
		Wallet: Wallet{
			UserID:  "40",
			Balance: 9307390,
		},
	}

	err := db.Omit(clause.Associations).Create(&user).Error
	assert.Nil(t, err)
}

func TestUserAddress(t *testing.T) {
	user := User{
		ID:       "50",
		Password: "rahasia",
		Name: Name{
			FirstName: "User 50",
		},
		Wallet: Wallet{
			UserID:  "50",
			Balance: 1000000,
		},
		Addresses: []Address{
			{
				UserID:  "50",
				Address: "Jalan jalan kemana pun",
			},
			{
				UserID:  "50",
				Address: "Coba tulis",
			},
		},
	}

	err := db.Create(&user).Error
	assert.Nil(t, err)
}

func TestUserPreload(t *testing.T) {
	var usersPreload []User
	err := db.Model(&User{}).Preload("Addresses").Joins("Wallet").Find(&usersPreload).Error
	assert.Nil(t, err)
}

func TestBelongsTo(t *testing.T) {
	var addresses []Address
	err := db.Model(&Address{}).Preload("User").Find(&addresses).Error
	assert.Nil(t, err)
}

var productID = int64(8174854164025333465)

func TestCreateManyToMany(t *testing.T) {
	product := Product{
		ID:    productID,
		Name:  faker.Name(),
		Price: 100000,
	}
	err := db.Create(&product).Error
	assert.Nil(t, err)

	err = db.Table("user_like_product").Create(map[string]any{
		"user_id":    "1",
		"product_id": productID,
	}).Error
	assert.Nil(t, err)

	err = db.Table("user_like_product").Create(map[string]any{
		"user_id":    "2",
		"product_id": productID,
	}).Error
	assert.Nil(t, err)
}

func TestPreloadManyToMany(t *testing.T) {
	var product Product
	err := db.Preload("LikedByUsers").First(&product, "id = ?", productID).Error
	assert.Nil(t, err)
	assert.Equal(t, 2, len(product.LikedByUsers))
}

func TestPreloadManyToManyProduct(t *testing.T) {
	var user User
	err := db.Preload("LikeProducts").First(&user, "id = ?", "1").Error
	assert.Nil(t, err)
	assert.Equal(t, 1, len(user.LikeProducts))
}

func TestAssociationFind(t *testing.T) {
	var product Product
	err := db.Take(&product, "id = ?", productID).Error
	assert.Nil(t, err)

	var users []User
	err = db.Model(&product).Where("users.first_name LIKE ?", "User%").Association("LikedByUsers").Find(&users)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
}

func TestAssociationAppend(t *testing.T) {
	var user User
	err := db.Take(&user, "id = ?", "3").Error
	assert.Nil(t, err)

	var product Product
	err = db.Take(&product, "id = ?", int64(8174854164025333465)).Error
	assert.Nil(t, err)

	err = db.Model(&product).Association("LikedByUsers").Append(&user)
	assert.Nil(t, err)
}

func TestAssociationDelete(t *testing.T) {
	var user User
	err := db.Take(&user, "id = ?", "3").Error
	assert.Nil(t, err)

	var product Product
	err = db.Take(&product, "id = ?", int64(8174854164025333465)).Error
	assert.Nil(t, err)

	err = db.Model(&product).Association("LikedByUsers").Delete(&user)
	assert.Nil(t, err)
}

func TestAssociationClear(t *testing.T) {
	var product Product
	err := db.Take(&product, "id = ?", int64(8174854164025333465)).Error
	assert.Nil(t, err)

	err = db.Model(&product).Association("LikedByUsers").Clear()
	assert.Nil(t, err)
}

func TestPreloadCondition(t *testing.T) {
	var user User
	err := db.Preload("Wallet", "balance > ?", 1000000).Take(&user, "id = ?", "1").Error
	assert.Nil(t, err)
}

func TestNestedPreload(t *testing.T) {
	var wallet Wallet
	err := db.Preload("User.Addresses").Take(&wallet, "id = ?", 1).Error
	assert.Nil(t, err)

	log.Println("wallet =====", wallet)
	log.Println("user =======", wallet.User)
	log.Println("address ====", wallet.User.Addresses)
}

func TestPreloadAll(t *testing.T) {
	var user User
	err := db.Preload(clause.Associations).Take(&user, "id = ?", "1").Error
	assert.Nil(t, err)
}

func TestJoinQuery(t *testing.T) {
	var users []User
	err := db.Joins("JOIN wallets ON wallets.user_id = users.id").Find(&users).Error
	assert.Nil(t, err)

	log.Println(len(users))

	users = []User{}
	err = db.Joins("Wallet").Find(&users).Error
	assert.Nil(t, err)

	log.Println(len(users))
}

func TestJoinWithCondition(t *testing.T) {
	// var users []User
	// err := db.Joins("JOIN wallets ON wallets.user_id = users.id AND wallets.balance > ?", 1000000).Find(&users).Error
	// assert.Nil(t, err)

	// log.Println(users)
	// log.Println(len(users))

	var users = []User{}
	err := db.Joins("Wallet").Where("Wallet.balance > ?", 1000000).Find(&users).Error
	assert.Nil(t, err)

	log.Println(users)
	log.Println(len(users))
}

func TestCount(t *testing.T) {
	var count int64
	err := db.Model(&User{}).Joins("Wallet").Where("\"Wallet\".balance > ?", 500000).Count(&count).Error
	assert.Nil(t, err)
	log.Println(count)
}

type AggregationResult struct {
	TotalBalance int64
	MinBalance   int64
	MaxBalance   int64
	AvgBalance   float64
}

func TestAggregation(t *testing.T) {
	var result AggregationResult
	err := db.Model(&Wallet{}).Select("sum(balance) AS total_balance", "min(balance) AS min_balance", "max(balance) AS max_balance",
		"avg(balance) AS avg_balance").Take(&result).Error

	assert.Nil(t, err)
	log.Println(result)
}

func TestGroupByHaving(t *testing.T) {
	var result []AggregationResult
	err := db.Model(&Wallet{}).Select("sum(balance) AS total_balance", "min(balance) AS min_balance", "max(balance) AS max_balance",
		"avg(balance) AS avg_balance").
		Joins("User").
		Group("\"User\".id").
		Having("sum(balance) > ?", 1000000).
		Find(&result).Error

	assert.Nil(t, err)
	log.Println(result)
}

func TestContext(t *testing.T) {
	ctx := context.Background()

	var users []User
	err := db.WithContext(ctx).Find(&users).Error

	assert.Nil(t, err)
}

func BrokeWallet(db *gorm.DB) *gorm.DB {
	return db.Where("balance = ?", 0)
}

func SultanWallet(db *gorm.DB) *gorm.DB {
	return db.Where("balance >= ?", 1000000)
}

func TestScope(t *testing.T) {
	var wallets []Wallet

	err := db.Scopes(BrokeWallet).Find(&wallets).Error
	assert.Nil(t, err)
}
