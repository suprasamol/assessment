package main

/*
Requirements

Expense tracking system
ให้สร้างระบบ REST API เพื่อจัดเก็บประวัติการใช้จ่าย (expense) ของลูกค้าธนาคาร ซึ่งความสามารถระบบมีดังนี้
ระบบสามารถจัดเก็บข้อมูล เรื่อง(title), ยอดค่าใช้จ่าย(amount), บันทึกย่อ(note) และ หมวดหมู่(tags)
ระบบสามารถเพิ่มประวัติการใช้จ่ายใหม่ได้
ระบบสามารถปรับเปลี่ยน/แก้ไข ข้อมูลของการใช้จ่ายได้
ระบบสามารถดึงข้อมูลการใช้จ่ายทั้งหมดออกมาแสดงได้
ระบบสามารถดึงข้อมูลการใช้จ่ายทีละรายการได้

*/

import (
	"database/sql"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	createTb := `
		CREATE TABLE IF NOT EXISTS expenses (id SERIAL PRIMARY KEY, title TEXT, amount FLOAT, note TEXT, tags TEXT[]);
	`
	_, err = db.Exec(createTb)
	if err != nil {
		log.Fatal("can't create table", err)
	}
}

func main() {
	InitDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	log.Printf("Server started at %v\n", os.Getenv("PORT"))
	log.Fatal(e.Start(os.Getenv("PORT")))
	log.Println("bye bye!")
}
