package repository

import (
	"book_shelf/domain"
	"context"
	"database/sql"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) CreateUser(ctx context.Context, u domain.User) (int, error) {
	var (
		id int
	)
	query := `
	INSERT INTO users (name, key, secret) VALUES ($1, $2, $3) RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query, u.Name, u.Key, u.Secret).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repo) CreateBook(ctx context.Context, b domain.Book) (int, error) {
	var (
		id int
	)
	query := `
	INSERT INTO books (isbn, title, author, published, pages, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query, b.ISBN, b.Title, b.Author, b.Published, b.Pages, b.Status).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (r *Repo) GetUserInfo(ctx context.Context, key string) (*domain.User, error) {
	u := domain.User{}
	query := `
	SELECT * FROM users WHERE key = $1
	`

	if err := r.db.QueryRowContext(ctx, query, key).Scan(&u.ID, &u.Name, &u.Key, &u.Secret); err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *Repo) GetAllBooks(ctx context.Context) ([]domain.Book, error) {
	books := make([]domain.Book, 0)

	query := `
	SELECT * FROM books
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		b := domain.Book{}
		if err := rows.Scan(&b.ID, &b.ISBN, &b.Title, &b.Author, &b.Published, &b.Pages, &b.Status); err != nil {
			return nil, err
		}

		books = append(books, b)
	}

	return books, err
}

func (r *Repo) Check(c context.Context, key string) (string, error) {
	var secret string
	query := `
	SELECT secret FROM users WHERE key = $1
	`
	err := r.db.QueryRowContext(c, query, key).Scan(&secret)
	if err != nil {
		return secret, err
	}

	return secret, err
}

func (r *Repo) EditBook(ctx context.Context, id int, status int) (*domain.Book, error) {
	book := domain.Book{}
	query := `
	UPDATE books SET status = $1 WHERE id = $2
	`
	if _, err := r.db.ExecContext(ctx, query, status, id); err != nil {
		return nil, err
	}
	query = `
	SELECT * FROM books WHERE id = $1
	`

	if err := r.db.QueryRowContext(ctx, query, id).Scan(&book.ID, &book.ISBN, &book.Title, &book.Author, &book.Published, &book.Pages, &book.Status); err != nil {
		return nil, err
	}

	return &book, nil
}

func (r *Repo) DeleteBook(ctx context.Context, id int) error {
	query := `
	DELETE FROM books WHERE id = $1
	`

	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return err
	}

	return nil
}
