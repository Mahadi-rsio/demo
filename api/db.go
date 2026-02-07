package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB(url string) *gorm.DB {

	db, err := gorm.Open(sqlite.Open(url), &gorm.Config{})

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	fmt.Println("Database connected")

	return db

}

func Migrate(db *gorm.DB, entity any) {
	err := db.AutoMigrate(entity)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Database Migrated")
}
