package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Book struct {
	Id     uint    `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float32 `db:"price"`
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True"

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	fmt.Printf("gorm open failed %v\n", err)
	// }
	// db.Create(&Book{Title: "《活着》", Author: "余华", Price: 34.6})
	// db.Create(&Book{Title: "《瓦登湖》", Author: "梭罗", Price: 46.7})
	// db.Create(&Book{Title: "《人间至味》", Author: "汪曾祺", Price: 55.3})
	// db.Create(&Book{Title: "《沙家浜》", Author: "汪曾祺", Price: 56.2})

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect failed %v\n", err)
		return
	}
	var books = []Book{}
	err = db.Select(&books, "Select * From books b Where b.price > ?", 50)
	if err != nil {
		fmt.Printf("db.Select failed %v\n", err)
	}
	fmt.Println(books)
}
