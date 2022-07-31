package store

import "errors"

var (
	ErrorNotFound = errors.New("not found")
	ErrorExist    = errors.New("exist")
)

type Book struct {
	Id      string   `json:"id"`      // 图示ISBN ID
	Name    string   `json:"name"`    // 图书名称
	Authors []string `json:"authors"` // 图书作者
	Press   string   `json:"press"`   // 出版社
}

type Store interface {
	Create(*Book) error               // 创建图书
	Update(*Book) error               // 更新图书
	GetBookByID(string) (Book, error) // 根据Id获取图书
	GetAllBooks() ([]Book, error)     // 获取所有图书
	DeleteBookByID(string) error      // 根据Id删除图书
}
