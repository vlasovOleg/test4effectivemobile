package sqlstore_test

import (
	"test4effectivemobile/internal/model"
	"test4effectivemobile/internal/storage"
	"test4effectivemobile/internal/storage/sqlstore"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestSave(t *testing.T) {
	db, truncate := sqlstore.TestDB(t)
	defer truncate("person")
	store := sqlstore.New(db)

	person := model.Person{
		Name:       "Test",
		Surname:    "Test",
		Patronymic: "Test",
	}

	id, err := store.Person().Save(person)
	require.NoError(t, err)
	require.NotEmpty(t, id)
}

func TestGetByID(t *testing.T) {
	db, truncate := sqlstore.TestDB(t)
	defer truncate("person")
	store := sqlstore.New(db)

	person := model.Person{
		Name:       "TestName",
		Surname:    "TestSurname",
		Patronymic: "TestPatronymic",
	}

	id, err := store.Person().Save(person)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	person2, err := store.Person().GetByID(id)
	require.NoError(t, err)
	person.ID = id
	require.Equal(t, person, person2)

	_, err = store.Person().GetByID(10)
	require.ErrorIs(t, storage.ErrNotFound, err)
}

func TestDellete(t *testing.T) {
	db, truncate := sqlstore.TestDB(t)
	defer truncate("person")
	store := sqlstore.New(db)

	person := model.Person{
		Name:       "Test",
		Surname:    "Test",
		Patronymic: "Test",
	}

	id, err := store.Person().Save(person)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	err = store.Person().Delete(id)
	require.NoError(t, err)

	person2, err := store.Person().GetByID(id)
	// assert.Error(t, err)
	require.ErrorIs(t, err, storage.ErrNotFound)
	require.Empty(t, person2)
}

func TestUpdate(t *testing.T) {
	db, truncate := sqlstore.TestDB(t)
	defer truncate("person")
	store := sqlstore.New(db)

	person := model.Person{
		ID:         0,
		Name:       "Name",
		Surname:    "Surname",
		Patronymic: "Patronymic",
		Age:        30,
		Gender:     "Gender",
		Country:    "Country",
	}

	var err error
	person.ID, err = store.Person().Save(person)
	require.NoError(t, err)
	require.NotEmpty(t, person.ID)

	person2 := model.Person{
		ID:         person.ID,
		Name:       "Name2",
		Surname:    "Surname2",
		Patronymic: "Patronymic2",
		Age:        32,
		Gender:     "Gender2",
		Country:    "Country2",
	}

	err = store.Person().Update(person2)
	require.NoError(t, err)

	person3, err := store.Person().GetByID(person.ID)
	require.NoError(t, err)
	require.Equal(t, person2, person3)
}

func TestGetByFilters(t *testing.T) {
	db, truncate := sqlstore.TestDB(t)
	defer truncate("person")
	store := sqlstore.New(db)

	persons := []model.Person{
		{Age: 20},
		{Age: 100},
		{Gender: "Gender"},
		{Country: "RU"},
	}

	table := []struct {
		filter model.Filter
		count  int
		name   string
	}{
		{
			filter: model.Filter{Limit: 2},
			name:   "Limit",
			count:  2,
		},
		{
			filter: model.Filter{Offset: 1, Limit: 100},
			count:  3,
			name:   "Offset",
		},
		{
			filter: model.Filter{AgeMin: 99, Limit: 100},
			count:  1,
			name:   "AgeMin",
		},
		{
			filter: model.Filter{AgeMax: 21, Limit: 100},
			count:  3,
			name:   "AgeMax",
		},
		{
			filter: model.Filter{Gender: "Gender", Limit: 100},
			count:  1,
			name:   "Gender",
		},
		{
			filter: model.Filter{Country: "RU", Limit: 100},
			count:  1,
			name:   "Country",
		},
	}

	// Prepare DB
	var err error
	for _, p := range persons {
		_, err = store.Person().Save(p)
		require.NoError(t, err)
	}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			persons, err = store.Person().Get(tc.filter)
			require.NoError(t, err)
			require.Len(t, persons, tc.count)
		})
	}
}
