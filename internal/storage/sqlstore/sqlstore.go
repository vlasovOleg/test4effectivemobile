package sqlstore

import (
	"database/sql"
	"test4effectivemobile/internal/storage"
)

type Storage struct {
	db               *sql.DB
	personRepository storage.Person
}

func New(db *sql.DB) *Storage {
	return &Storage{
		db: db,
		personRepository: Person{
			db: db,
		},
	}
}

func (s Storage) Close() error {
	return s.db.Close()
}

func (s Storage) Person() storage.Person {
	return s.personRepository
}
