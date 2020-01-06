package item_repo

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/efimovad/avito-internship/internal/model"
	"github.com/lib/pq"
	"testing"
	"time"
)

func testItem(t *testing.T) *model.Item {
	t.Helper()
	return &model.Item{
		ID: 		 1,
		Title:       "My item",
		Description: "My item description",
		Date:        time.Now(),
		Price:       100,
		MainImage:   "http://image/1",
		Images:      [3]string{"http://image/1", "http://image/2", "http://image/3"},
	}
}

func compare(l model.Item, r model.Item, t *testing.T) bool {
	t.Helper()
	return l.Title == r.Title && l.Price == r.Price && l.Images[0] == r.Images[0]
}

func compareAll(l model.Item, r model.Item, t *testing.T) bool {
	t.Helper()
	return l.Title == r.Title &&
		l.Price == r.Price &&
		l.Images[0] == r.Images[0] &&
		l.Images[1] == r.Images[1] &&
		l.Images[2] == r.Images[2] &&
		l.Description == r.Description
}

func TestRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}

	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	repo := NewItemRepository(db)

	rows := sqlmock.NewRows([]string{"id"})

	var elemID int64 = 1
	expect := []*model.Item{
		{ID: elemID},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID)
	}

	i := testItem(t)

	//ok query
	mock.
		ExpectQuery(`INSERT INTO items`).
		WithArgs(i.Title, i.Description, i.Date, i.Price, pq.Array(i.Images)).
		WillReturnRows(rows)

	err = repo.Create(i)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error
	mock.
		ExpectQuery(`INSERT INTO items`).
		WithArgs(i.Title, i.Description, i.Date, i.Price, pq.Array(i.Images)).
		WillReturnError(fmt.Errorf("bad query"))

	err = repo.Create(i)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRepository_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}

	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	var elemID int64 = 1

	// good query
	rows := sqlmock.
		NewRows([]string{"id", "title", "price", "images"})

	expect := []*model.Item{
		testItem(t),
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Title, item.Price, item.Images[0])
	}

	mock.
		ExpectQuery("SELECT").
		WithArgs(elemID).
		WillReturnRows(rows)

	repo := NewItemRepository(db)


	item, err := repo.Get(elemID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !compare(*item, *expect[0], t) {
		t.Errorf("results not match, want %v, have %v", expect[0], item)
		return
	}

	// query error
	mock.
		ExpectQuery("SELECT").
		WithArgs(elemID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.Get(elemID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// row scan error
	expect = []*model.Item{
		testItem(t),
	}

	mock.
		ExpectQuery("SELECT").
		WithArgs(elemID).
		WillReturnRows(rows)

	_, err = repo.Get(elemID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}

	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	var elemID int64 = 1

	// good query
	rows := sqlmock.
		NewRows([]string{"id", "title", "price", "description", "images"})

	expect := []*model.Item{
		testItem(t),
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Title, item.Price, item.Description, pq.Array(item.Images))
	}

	mock.
		ExpectQuery("SELECT").
		WithArgs(elemID).
		WillReturnRows(rows)

	repo := NewItemRepository(db)


	item, err := repo.GetAll(elemID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !compare(*item, *expect[0], t) {
		t.Errorf("results not match, want %v, have %v", expect[0], item)
		return
	}

	// query error
	mock.
		ExpectQuery("SELECT").
		WithArgs(elemID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.Get(elemID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// row scan error
	expect = []*model.Item{
		testItem(t),
	}

	mock.
		ExpectQuery("SELECT").
		WithArgs(elemID).
		WillReturnRows(rows)

	_, err = repo.Get(elemID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestRepository_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	// good query
	rows := sqlmock.
		NewRows([]string{"title", "images", "price"})

	expect := []*model.Item{
		testItem(t),
		testItem(t),
		testItem(t),
	}

	for _, item := range expect {
		rows = rows.AddRow(item.Title, item.Images[0], item.Price)
	}

	params := model.Params{
		Date:  false,
		Price: false,
		Desc:  false,
		Page:  1,
	}

	mock.
		ExpectQuery("SELECT").
		WithArgs(params.Desc, params.Date, params.Price, params.Page).
		WillReturnRows(rows)

	repo := NewItemRepository(db)

	items, err := repo.List(params)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !compare(items[0], *expect[0], t) {
		t.Errorf("results not match, want %v, have %v", *expect[0], items[0])
		return
	}

	// query error
	mock.
		ExpectQuery("SELECT").
		WithArgs(params.Desc, params.Date, params.Price, params.Page).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.List(params)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}