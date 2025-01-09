package postgres

import (
	"database/sql"

	"github.com/heymmhr/students-api/internal/config"
	_ "github.com/lib/pq"
)

type Postgres struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Postgres, error) {

	db, err := sql.Open("postgres", "user=postgres password=password dbname=student_db sslmode=disable")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
	id SERIAL PRIMARY KEY,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Postgres{
		Db: db,
	}, nil

}

func (p *Postgres) CreateStudent(name string, email string, age int) (int64, error) {

	query := "INSERT INTO students (name, email, age) VALUES($1, $2, $3) RETURNING id"

	var id int64
	err := p.Db.QueryRow(query, name, email, age).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
