package store

import (
	mystore "bookstore/store"
	"bookstore/store/factory"
	"sync"
)

type MemStore struct {
	sync.RWMutex
	books map[string]*mystore.Book
}

func init() {
	factory.Register("mem", &MemStore{
		books: make(map[string]*mystore.Book),
	})
}

// Create creates a new book in the store
func (ms *MemStore) Create(book *mystore.Book) error {
	ms.Lock()
	defer ms.Unlock()

	// check
	if _, ok := ms.books[book.Id]; !ok {
		return mystore.ErrorExist
	}

	// save
	nBook := *book
	ms.books[book.Id] = &nBook
	return nil
}

// Update updates the existed Book in the store
func (ms *MemStore) Update(book *mystore.Book) error {
	ms.Lock()
	defer ms.Unlock()

	oldBook, ok := ms.books[book.Id]
	// check book id exist
	if !ok {
		return mystore.ErrorNotFound
	}

	// check
	nBook := *oldBook
	if book.Name != "" {
		nBook.Name = book.Name
	}
	if book.Authors != nil {
		nBook.Authors = book.Authors
	}
	if book.Press != "" {
		nBook.Press = book.Press
	}
	// update
	ms.books[book.Id] = &nBook

	return nil
}

// GetBookByID retrieves a book from the store, by id. If no such id exists. an error is required.
func (ms *MemStore) GetBookByID(id string) (mystore.Book, error) {
	ms.RLock()
	defer ms.RUnlock()

	book, ok := ms.books[id]
	if !ok {
		// not found id, return err
		return mystore.Book{}, mystore.ErrorNotFound
	}
	return *book, nil
}

// GetAllBooks returns all the books in the store, in arbitrary order.
func (ms *MemStore) GetAllBooks() ([]mystore.Book, error) {
	ms.RLock()
	defer ms.RUnlock()

	// get all books in the store
	allBooks := make([]mystore.Book, 0, len(ms.books))
	for _, book := range ms.books {
		allBooks = append(allBooks, *book)
	}
	return allBooks, nil
}

// DeleteBookByID deletes the book with the given id. If no such id exist. an error is returned.
func (ms *MemStore) DeleteBookByID(id string) error {
	ms.Lock()
	defer ms.Unlock()

	// check id exists in the store
	if _, ok := ms.books[id]; !ok {
		return mystore.ErrorNotFound
	}
	// delete book in the store by id
	delete(ms.books, id)
	return nil
}
