package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Alvesafk/shor/internal/app"
)

type Connection struct {
	db *app.FireDB
}

func NewConnection(db *app.FireDB) *Connection {
	return &Connection{db}
}

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Content any    `json:"content"`
}

func (r Response) WriteJSON(w http.ResponseWriter, h int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(h)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(r)
}

func (c Connection) GetHelloWorld(w http.ResponseWriter, r *http.Request) {
	Response{
		Message: "This is a Hello, World!",
		Status:  "ok",
		Content: "Hello, World!",
	}.WriteJSON(w, http.StatusOK)
}
