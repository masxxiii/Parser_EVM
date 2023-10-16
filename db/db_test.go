package db

import (
	"fmt"
	"github.com/ory/dockertest/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
	"time"
)

var database *gorm.DB

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	pool.MaxWait = time.Minute * 2
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "latest", []string{"MYSQL_ROOT_PASSWORD=root"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// wait for 30 seconds and then connect to Database
	time.Sleep(30 * time.Second)
	database, err = gorm.Open(mysql.Open(
		fmt.Sprintf("root:root@tcp(localhost:%s)/mysql?parseTime=True", resource.GetPort("3306/tcp"))), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	// Migrate database models
	migrateErr := database.AutoMigrate(&Block{}, &Wallet{}, &Transaction{})
	if migrateErr != nil {
		log.Fatal("[Migration Error]", migrateErr)
	}

	// Add mock data into the database
	createMock()

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func createMock() {
	BlockErr := database.Create(&Block{
		Height:  1000,
		Network: "ETH",
	}).Error
	if BlockErr != nil {
		log.Fatal(BlockErr)
	}

	TrxErr := database.Create(&Transaction{
		Hash:  "xxx",
		From:  "abc",
		To:    "def",
		Value: "100",
	}).Error
	if TrxErr != nil {
		log.Fatal(TrxErr)
	}

	WalletErr := database.Create(&Wallet{
		Address: "zzz",
	}).Error
	if WalletErr != nil {
		log.Fatal(WalletErr)
	}
}

func TestConnectDB(t *testing.T) {
	var block = Block{ID: 1}
	err1 := database.First(&block).Error
	if err1 != nil {
		t.Fatal(err1)
	}

	var trx = Transaction{}
	err2 := database.Where("hash = ?", "xxx").First(&trx).Error
	if err2 != nil {
		t.Fatal(err2)
	}

	var wallet = Wallet{}
	err3 := database.Where("address = ?", "zzz").First(&wallet).Error
	if err3 != nil {
		t.Fatal(err3)
	}
}
