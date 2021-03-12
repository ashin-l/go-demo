package main

import (
	"fmt"
)

type Book struct {
	Name string
}

func (b Book) GetName() string {
	return b.Name
}

type BookShelf struct {
	Books []Book
	Last  int
}

func (b *BookShelf) GetBookAt(index int) Book {
	return b.Books[index]
}

func (b *BookShelf) AppendBook(book Book) {
	b.Books = append(b.Books, book)
	b.Last++
}

func (b *BookShelf) GetLength() int {
	return b.Last
}

func (b BookShelf) Iterator() Iterator {
	return &BookShelfIterator{bs: b}
}

type BookShelfIterator struct {
	bs    BookShelf
	index int
}

func (b *BookShelfIterator) HasNext() bool {
	return b.index < b.bs.GetLength()
}

func (b *BookShelfIterator) Next() interface{} {
	book := b.bs.GetBookAt(b.index)
	b.index++
	return book
}

type Iterator interface {
	HasNext() bool
	Next() interface{}
}

func main() {
	bs := BookShelf{}
	bs.AppendBook(Book{"a"})
	bs.AppendBook(Book{"b"})
	bs.AppendBook(Book{"c"})
	it := bs.Iterator()
	for it.HasNext() {
		b := it.Next().(Book)
		fmt.Println(b.GetName())
	}
}
