package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetAllExpense(t *testing.T) {
	seedExpense(t)

	var ex []Expense
	res := request(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&ex)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Greater(t, len(ex), 0)

}

func TestCreateExpense(t *testing.T) {
	body := bytes.NewBufferString(`{
		"title": "Dell Precision M3800",
		"amount": 1,
		"note": "discount 1000 bath for use code 'dis1000'", 
		"tags": ["Notebook", "Dell"]
	}`)

	var e Expense
	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, e.ID)
	assert.Equal(t, "Dell Precision M3800", e.Title)
	assert.Equal(t, 1, e.Amount)
	assert.Equal(t, "discount 1000 bath for use code 'dis1000'", e.Note)
	assert.Equal(t, pq.Array([]string{"Notebook", "Dell"}), pq.Array(e.Tags))
}

func TestGetExpenseByID(t *testing.T) {
	e := seedExpense(t)

	var latest Expense
	res := request(http.MethodGet, uri("expenses", strconv.Itoa(e.ID)), nil)
	err := res.Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, e.ID, latest.ID)
	assert.NotEmpty(t, latest.Title)
	assert.NotEmpty(t, latest.Amount)
	assert.NotEmpty(t, latest.Note)
	assert.NotEmpty(t, latest.Tags)
}

func seedExpense(t *testing.T) Expense {
	var e Expense
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`)
	err := request(http.MethodPost, uri("expenses"), body).Decode(&e)
	if err != nil {
		t.Fatal("can't create expense: ", err)
	}

	return e
}

func uri(paths ...string) string {
	host := "http://localhost" + os.Getenv("PORT")
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}
