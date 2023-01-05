package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

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
