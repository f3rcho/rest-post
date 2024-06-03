package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/f3rcho/rest-posts/models"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostGresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db}, nil
}

func (p *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	_, err := p.db.ExecContext(ctx, "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)", user.ID, user.Email, user.Password)
	return err
}

func (p *PostgresRepository) GetUserByID(ctx context.Context, ID string) (*models.User, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT id, email FROM users WHERE id = $1", ID)
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var user = models.User{}

	for rows.Next() {
		if err = rows.Scan(&user.ID, &user.Email); err == nil {
			return &user, nil
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &user, nil
}
func (p *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT id, email, password FROM users WHERE email = $1", email)
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var user = models.User{}

	for rows.Next() {
		if err = rows.Scan(&user.ID, &user.Email, &user.Password); err == nil {
			return &user, nil
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *PostgresRepository) InsertPost(ctx context.Context, post *models.Post) error {
	_, err := p.db.ExecContext(ctx, "INSERT INTO posts (id, post_content, created_at, user_id) VALUES ($1, $2, $3, $4)", post.Id, post.PostContent, post.CreatedAt, post.UserId)
	return err
}

func (p *PostgresRepository) Close() error {
	return p.db.Close()
}
