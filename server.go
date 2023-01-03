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
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Expense struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount int      `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:tags`
}

type Err struct {
	Message string `json:"message"`
}

var db *sql.DB

func CreateExpenseHandler(c echo.Context) error {
	var e Expense
	err := c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	row := db.QueryRow("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4)  RETURNING id", e.Title, e.Amount, e.Note, pq.Array(e.Tags))

	err = row.Scan(&e.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, e)
}

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

	e.POST("/expenses", CreateExpenseHandler)

	log.Printf("Server started at %v\n", os.Getenv("PORT"))
	log.Fatal(e.Start(os.Getenv("PORT")))
	log.Println("bye bye!")
}
