package main

import (
	"context"
	"fmt"

	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/suprasamol/assessment/expense"

	"github.com/labstack/gommon/log"
)

func main() {
	expense.InitDB()

	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "apidesign" && password == "45678" {
			return true, nil
		}
		return false, nil
	}))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/expenses", expense.CreateExpenseHandler)
	e.GET("/expenses", expense.GetExpensesHandler)
	e.GET("/expenses/:id", expense.GetExpenseHandler)
	e.PUT("/expenses/:id", expense.UpdateExpenseHandler)

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
