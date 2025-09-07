package data

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	InitDBConnection(db)
	return db, mock
}

func TestCanPurchaseWithMock(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	// Expect query with quantity=3, id=1 → true
	mock.ExpectQuery("SELECT \\(quantity >= \\?\\) FROM album WHERE id = \\?").
		WithArgs(int64(3), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"enough"}).AddRow(true))

	ok, err := CanPurchase(1, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Errorf("expected true, got false")
	}

	// Ensure all expectations met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
