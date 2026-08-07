package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Alvesafk/shor/internal/app"
	"github.com/Alvesafk/shor/internal/models"
)

const (
	ok   = "ok"
	fail = "failed"
)

var (
	errEmptyString = errors.New("the passed string is empty")
	errNotFound    = errors.New("the content was not found")
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
		Status:  ok,
		Content: "Hello, World!",
	}.WriteJSON(w, http.StatusOK)
}

func (c Connection) PostURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Response{
			Message: "method not allowed",
			Status:  fail,
		}.WriteJSON(w, http.StatusMethodNotAllowed)

		return
	}

	var u struct {
		Url string `json:"url"`
	}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil || u.Url == "" {
		Response{
			Message: "invalid json request",
			Status:  fail,
		}.WriteJSON(w, http.StatusBadRequest)

		return
	}

	url, alreadyRegistered, err := c.db.URLAlreadyRegistered(u.Url, r.Context())
	if err != nil {
		Response{
			Message: "error on checking if url is already registered",
			Status:  fail,
		}.WriteJSON(w, http.StatusInternalServerError)

		return
	}

	if alreadyRegistered {
		Response{
			Message: "url already registered",
			Status:  fail,
			Content: url,
		}.WriteJSON(w, http.StatusConflict)

		return
	}

	shortCode, err := generateShortURL()
	if err != nil {
		Response{
			Message: "error on generating short code",
			Status:  fail,
		}.WriteJSON(w, http.StatusInternalServerError)

		return
	}

	for {
		exist, err := c.db.ShortURLExists(shortCode, r.Context())
		if err != nil {
			Response{
				Message: "error on checking short code",
				Status:  fail,
			}.WriteJSON(w, http.StatusInternalServerError)

			return
		}

		if !exist {
			break
		}

		shortCode, err = generateShortURL()
		if err != nil {
			Response{
				Message: "error on generating short code",
				Status:  fail,
			}.WriteJSON(w, http.StatusInternalServerError)

			return
		}
	}

	url = models.GenURLStruct(u.Url, shortCode)

	err = c.db.CreateURL(*url, r.Context())
	if err != nil {
		Response{
			Message: "error on creating db entry",
			Status:  fail,
		}.WriteJSON(w, http.StatusInternalServerError)

		return
	}

	Response{
		Message: "success",
		Status:  ok,
		Content: url,
	}.WriteJSON(w, http.StatusOK)
}

func (c Connection) GetURL(w http.ResponseWriter, r *http.Request) {
	shortUrlString := r.PathValue("shortUrl")
	res, code, err := ValidateShortURL(shortUrlString, c, r.Context())
	if err != nil {
		res.WriteJSON(w, code)

		return
	}

	url, err := c.db.GetURLByShortCode(shortUrlString, r.Context())
	if err != nil {
		Response{
			Message: "could not find url",
			Status:  fail,
		}.WriteJSON(w, http.StatusInternalServerError)

		return
	}

	Response{
		Message: "success",
		Status:  ok,
		Content: url,
	}.WriteJSON(w, http.StatusOK)
}

func (c Connection) UpdateURL(w http.ResponseWriter, r *http.Request) {
	shortUrlString := r.PathValue("shortUrl")
	res, code, err := ValidateShortURL(shortUrlString, c, r.Context())
	if err != nil {
		res.WriteJSON(w, code)

		return
	}

	var u struct {
		Url string `json:"url"`
	}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil || u.Url == "" {
		Response{
			Message: "invalid json request",
			Status:  fail,
		}.WriteJSON(w, http.StatusBadRequest)

		return
	}

	url, err := c.db.UpdateURL(u.Url, shortUrlString, r.Context())
	if err != nil {
		Response{
			Message: "could not update url",
			Status:  fail,
		}.WriteJSON(w, http.StatusInternalServerError)

		return
	}

	Response{
		Message: "success",
		Status:  ok,
		Content: url,
	}.WriteJSON(w, http.StatusOK)
}

func (c Connection) DeleteURL(w http.ResponseWriter, r *http.Request) {
	shortUrlString := r.PathValue("shortUrl")
	res, code, err := ValidateShortURL(shortUrlString, c, r.Context())
	if err != nil {
		res.WriteJSON(w, code)

		return
	}

	if err := c.db.DeleteURL(shortUrlString, r.Context()); err != nil {
		Response{
			Message: "could not delete the url",
			Status:  fail,
		}.WriteJSON(w, http.StatusInternalServerError)

		return
	}

	Response{
		Message: "success",
		Status:  ok,
	}.WriteJSON(w, http.StatusNoContent)
}

func ValidateShortURL(shortURL string, c Connection, ctx context.Context) (*Response, int, error) {
	if shortURL == "" {
		return &Response{Message: "invalid json request", Status: fail}, http.StatusBadRequest, errEmptyString
	}

	exist, err := c.db.ShortURLExists(shortURL, ctx)
	if err != nil {
		if errors.Is(err, app.UrlNotFound) {
			return &Response{Message: "url not found", Status: fail}, http.StatusNotFound, app.UrlNotFound
		}

		return &Response{Message: "could not find url", Status: fail}, http.StatusInternalServerError, errNotFound
	}

	if !exist {
		return &Response{Message: "short url does not exist", Status: fail}, http.StatusBadRequest, errNotFound
	}

	return nil, http.StatusOK, nil
}
