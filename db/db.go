package db

import (
	"errors"
	"fmt"
	"management/model"
	"os"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
)

var txdbInitialized bool
var mutex sync.Mutex

func getDatabaseConnectString() (string, error) {
	host := os.Getenv("DB_HOST")
	if "" == host {
		return "", errors.New("DB_HOST : $DB_HOST is not set")
	}

	user := os.Getenv("DB_USER")
	if "" == user {
		return "", errors.New("DB_USER : $DB_USER is not set")
	}

	password := os.Getenv("DB_PASSWORD")
	if "" == password {
		return "", errors.New("DB_PASSWORD : $DB_PASSWORD is not set")
	}

	name := os.Getenv("DB_NAME")
	if "" == name {
		return "", errors.New("DB_NAME : $DB_NAME is not set")
	}

	port := os.Getenv("DB_PORT")
	if "" == port {
		return "", errors.New("DB_PORT : $DB_PORT is not set")
	}

	options := "charset=utf8mb4&parseTime=True&loc=Local"

	return fmt.Sprintf("%s:%s@(%s:%s)/%s?%s", user, password, host, port, name, options), nil
}

func CreateDatabaseConnection() (*gorm.DB, error) {
	connectString, err := getDatabaseConnectString()
	if nil != err {
		return nil, err
	}

	var db *gorm.DB

	for i := 0; i < 10; i++ {
		db, err = gorm.Open("mysql", connectString)
		if nil == err {
			break
		}

		time.Sleep(1 * time.Second)
	}

	if nil != err {
		return nil, err
	}

	db.DB().SetMaxIdleConns(10)
	db.LogMode(true)

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.User{},
	).Error
	if err != nil {
		return err
	}
	return nil
}
