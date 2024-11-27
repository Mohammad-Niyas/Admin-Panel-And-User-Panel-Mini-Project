package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"main.go/model"
)

var DB *gorm.DB

func DBconnect(){
	dsn:= "host=localhost user=postgres password=466450 dbname=project port=5432"
	db,err:=gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err!=nil{
		log.Fatal("Failed to Connect Databse")
	}
	fmt.Println("Database Connected Successfully")
	DB = db
	DB.AutoMigrate(&model.UserModel{},&model.AdminModel{})
}
