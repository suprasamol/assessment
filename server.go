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
	"context"
	"database/sql"
	"fmt"

	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Expense struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount int      `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
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

func GetExpensesHandler(c echo.Context) error {
	stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query all expense statement"})
	}

	rows, err := stmt.Query()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't query all expenses"})
	}

	expenses := []Expense{}
	for rows.Next() {
		var e Expense
		err = rows.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expenses"})
		}
		expenses = append(expenses, e)

	}
	return c.JSON(http.StatusOK, expenses)
}

func GetExpenseHandler(c echo.Context) error {
	id := c.Param("id")
	stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query expense statment:" + err.Error()})
	}

	row := stmt.QueryRow(id)
	e := Expense{}
	err = row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "expense not found"})
	case nil:
		return c.JSON(http.StatusOK, e)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expense:" + err.Error()})
	}
}

func UpdateExpenseHandler(c echo.Context) error {
	var e Expense
	err := c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	_, err = db.Exec("UPDATE expenses SET title=$1, amount=$2, note=$3, tags=$4 WHERE id=$5", e.Title, e.Amount, e.Note, pq.Array(e.Tags), c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, e)
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
	e.Logger.SetLevel(log.INFO)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/expenses", CreateExpenseHandler)
	e.GET("/expenses", GetExpensesHandler)
	e.GET("/expenses/:id", GetExpenseHandler)
	e.PUT("/expenses/:id", UpdateExpenseHandler)

	go func() {
		if err := e.Start(os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	fmt.Println("bye bye")

}
