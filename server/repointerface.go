package server

import (
	"book_shelf/domain"
	"context"
)

type Repository interface {
	CreateUser(ctx context.Context, u domain.User) (int, error)
	CreateBook(ctx context.Context, b domain.Book) (int, error)
	GetUserInfo(ctx context.Context, key string) (*domain.User, error)
	GetAllBooks(ctx context.Context) ([]domain.Book, error)
	Check(c context.Context, key string) (string, error)
	EditBook(ctx context.Context, id int, status int) (*domain.Book, error)
	DeleteBook(ctx context.Context, id int) error
}
