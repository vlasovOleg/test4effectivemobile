package sqlstore

import (
	"database/sql"
	"errors"
	"strings"
	"test4effectivemobile/internal/model"
	"test4effectivemobile/internal/storage"
)

type Person struct {
	db *sql.DB
}

func (p Person) Save(pe model.Person) (int64, error) {
	stmt, err := p.db.Prepare(
		`INSERT INTO person (
			Name, Surname, Patronymic,
			Age, Gender, Country
		)
		VALUES($1, $2, $3, $4, $5, $6)
		RETURNING ID;`,
	)
	if err != nil {
		return 0, errors.Join(storage.ErrInternal, err)
	}
	defer func() { _ = stmt.Close() }()

	row := stmt.QueryRow(
		pe.Name, pe.Surname, pe.Patronymic,
		pe.Age, pe.Gender, pe.Country,
	)
	if row.Err() != nil {
		return 0, errors.Join(storage.ErrInternal, err)
	}

	var id int64
	err = row.Scan(&id)
	if err != nil {
		return 0, errors.Join(storage.ErrInternal, err)
	}

	return id, nil
}

func (p Person) GetByID(id int64) (model.Person, error) {
	stmt, err := p.db.Prepare(
		`SELECT * FROM person WHERE ID = $1;`)
	if err != nil {
		return model.Person{}, errors.Join(storage.ErrInternal, err)
	}
	defer func() { _ = stmt.Close() }()

	row := stmt.QueryRow(id)
	if row.Err() != nil {
		return model.Person{}, errors.Join(storage.ErrInternal, err)
	}

	person := model.Person{}
	err = row.Scan(
		&person.ID,
		&person.Name,
		&person.Surname,
		&person.Patronymic,
		&person.Age,
		&person.Gender,
		&person.Country,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Person{}, storage.ErrNotFound
	}
	if err != nil {
		return model.Person{}, errors.Join(storage.ErrInternal, err)
	}

	return person, nil
}

func (p Person) Delete(id int64) error {
	stmt, err := p.db.Prepare(`DELETE FROM person WHERE ID = $1 RETURNING id;`)
	defer func() { _ = stmt.Close() }()
	if err != nil {
		return errors.Join(storage.ErrInternal, err)
	}

	row := stmt.QueryRow(id)
	if row.Err() != nil {
		return errors.Join(storage.ErrInternal, err)
	}

	if errors.Is(row.Scan(), sql.ErrNoRows) {
		return storage.ErrNotFound
	}

	return nil
}

func (p Person) Update(pe model.Person) error {
	stmt, err := p.db.Prepare(
		`UPDATE person SET 
			Name = $2, Surname = $3, Patronymic = $4,
			Age = $5, Gender = $6, Country = $7
		WHERE ID = $1 RETURNING id;`,
	)
	if err != nil {
		return errors.Join(storage.ErrInternal, err)
	}
	defer func() { _ = stmt.Close() }()

	row := stmt.QueryRow(
		pe.ID, pe.Name, pe.Surname, pe.Patronymic,
		pe.Age, pe.Gender, pe.Country,
	)
	if row.Err() != nil {
		return errors.Join(storage.ErrInternal, err)
	}

	err = row.Scan()
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrNotFound
	}

	return nil
}

func (p Person) Get(f model.Filter) ([]model.Person, error) {
	const stmtLen = 120
	/*
		"SELECT * FROM person WHERE 1=1
		AND Age > $1
		AND Age < $2
		AND Gender == $3
		AND Country == $4
		LIMIT $5
		OFFSET $6
		;
	*/

	paramsNum := []string{"$1", "$2", "$3", "$4", "$5", "$6"}
	params := make([]any, 0, len(paramsNum))
	request := strings.Builder{}
	request.Grow(stmtLen)

	request.WriteString("SELECT * FROM person WHERE 1=1 ")

	if f.AgeMax != 0 {
		params = append(params, &f.AgeMax)
		request.WriteString(" AND age < ")
		request.WriteString(paramsNum[len(params)-1])
	}
	if f.AgeMin != 0 {
		params = append(params, &f.AgeMin)
		request.WriteString(" AND age > ")
		request.WriteString(paramsNum[len(params)-1])
	}
	if len(f.Gender) != 0 {
		params = append(params, &f.Gender)
		request.WriteString(" AND gender = ")
		request.WriteString(paramsNum[len(params)-1])
	}
	if len(f.Country) != 0 {
		params = append(params, &f.Country)
		request.WriteString(" AND country = ")
		request.WriteString(paramsNum[len(params)-1])
	}

	params = append(params, &f.Limit)
	request.WriteString(" LIMIT ")
	request.WriteString(paramsNum[len(params)-1])

	params = append(params, &f.Offset)
	request.WriteString(" OFFSET ")
	request.WriteString(paramsNum[len(params)-1])
	request.WriteString(" ;")

	stmt, err := p.db.Prepare(request.String())
	defer func() { _ = stmt.Close() }()
	if err != nil {
		return []model.Person{}, errors.Join(storage.ErrInternal, err, errors.New(request.String()))
	}
	defer func() { _ = stmt.Close() }()

	rows, err := stmt.Query(params...)
	defer func() { _ = rows.Close() }()
	if err != nil || rows.Err() != nil {
		return []model.Person{}, errors.Join(storage.ErrInternal, err)
	}

	persons := []model.Person{}
	for rows.Next() {
		person := model.Person{}
		err = rows.Scan(
			&person.ID,
			&person.Name,
			&person.Surname,
			&person.Patronymic,
			&person.Age,
			&person.Gender,
			&person.Country,
		)
		if err != nil {
			return []model.Person{}, errors.Join(storage.ErrInternal, err)
		}
		persons = append(persons, person)
	}

	return persons, nil
}
