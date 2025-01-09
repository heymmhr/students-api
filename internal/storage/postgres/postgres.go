package postgres

import (
	"database/sql"
	"fmt"

	"github.com/heymmhr/students-api/internal/config"
	"github.com/heymmhr/students-api/internal/types"
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

func (p *Postgres) GetStudentById(id int64) (types.Student, error) {

	query := "SELECT id, name, email, age FROM students WHERE id =$1"

	var student types.Student
	err := p.Db.QueryRow(query, id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)

	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %d", id)
		}
		return types.Student{}, fmt.Errorf("failed to get student: %w", err)
	}
	return student, nil
}
