package server

import (
	"bookstore/server/middleware"
	"bookstore/store"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type BookStoreServer struct {
	store store.Store
	srv   *http.Server
}

func NewBookStoreServer(addr string, s store.Store) *BookStoreServer {
	srv := &BookStoreServer{
		store: s,
		srv: &http.Server{
			Addr: addr,
		},
	}

	router := mux.NewRouter()
	router.HandleFunc("/book", srv.createBookHandler).Methods("POST")
	router.HandleFunc("/book/{id}", srv.updateBookHandler).Methods("POST")
	router.HandleFunc("/book/{id}", srv.getBookByIdHandler).Methods("GET")
	router.HandleFunc("/book", srv.getAllBooksHandler).Methods("GET")
	router.HandleFunc("/book/{id}", srv.deleteBookByIdHandler).Methods("DELETE")

	srv.srv.Handler = middleware.Logging(middleware.Validating(router))
	return srv
}

func (bs *BookStoreServer) ListenAndServe() (<-chan error, error) {
	var err error
	errChan := make(chan error)
	go func() {
		err = bs.srv.ListenAndServe()
		errChan <- err
	}()

	select {
	case err = <-errChan:
		return nil, err
	case <-time.After(time.Second):
		return errChan, nil
	}
}

func (bs *BookStoreServer) Shutdown(ctx context.Context) error {
	return bs.srv.Shutdown(ctx)
}

func (bs *BookStoreServer) createBookHandler(writer http.ResponseWriter, request *http.Request) {
	// Book 解码器
	var decoder = json.NewDecoder(request.Body)
	var book store.Book

	if err := decoder.Decode(&book); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := bs.store.Create(&book); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
}

func (bs *BookStoreServer) updateBookHandler(writer http.ResponseWriter, request *http.Request) {
	id, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, "no id found in request", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(request.Body)
	var book store.Book
	if err := decoder.Decode(&book); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	book.Id = id
	if err := bs.store.Update(&book); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
}

func (bs *BookStoreServer) getBookByIdHandler(writer http.ResponseWriter, request *http.Request) {
	id, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, "no id found in request", http.StatusBadRequest)
		return
	}
	bookByID, err := bs.store.GetBookByID(id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	response(writer, bookByID)
}

func (bs *BookStoreServer) getAllBooksHandler(writer http.ResponseWriter, _ *http.Request) {
	books, err := bs.store.GetAllBooks()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	response(writer, books)
}

func (bs *BookStoreServer) deleteBookByIdHandler(writer http.ResponseWriter, request *http.Request) {
	id, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, "no id found in request", http.StatusBadRequest)
		return
	}

	err := bs.store.DeleteBookByID(id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
}

func response(writer http.ResponseWriter, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(data)
}
