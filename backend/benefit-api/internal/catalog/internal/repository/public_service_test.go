package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/carlosEA28/benefitGuide/internal/catalog/internal/models"
	"github.com/carlosEA28/benefitGuide/internal/shared/exceptions"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type publicServiceConnector struct {
	query func(string, []driver.NamedValue) (driver.Rows, error)
}

func (c publicServiceConnector) Connect(context.Context) (driver.Conn, error) {
	return publicServiceConn{query: c.query}, nil
}
func (c publicServiceConnector) Driver() driver.Driver { return publicServiceDriver{} }

type publicServiceDriver struct{}

func (publicServiceDriver) Open(string) (driver.Conn, error) { return nil, errors.New("use connector") }

type publicServiceConn struct {
	query func(string, []driver.NamedValue) (driver.Rows, error)
}

func (publicServiceConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("unexpected prepared statement")
}
func (publicServiceConn) Close() error              { return nil }
func (publicServiceConn) Begin() (driver.Tx, error) { return nil, errors.New("unexpected transaction") }
func (c publicServiceConn) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return c.query(query, args)
}

type publicServiceRows struct {
	values [][]driver.Value
	index  int
}

func (r *publicServiceRows) Columns() []string { return []string{"id", "status"} }
func (r *publicServiceRows) Close() error      { return nil }
func (r *publicServiceRows) Next(dest []driver.Value) error {
	if r.index == len(r.values) {
		return io.EOF
	}
	copy(dest, r.values[r.index])
	r.index++
	return nil
}

func testPublicServiceRepository(t *testing.T, query func(string, []driver.NamedValue) (driver.Rows, error)) *repository {
	t.Helper()
	db := sql.OpenDB(publicServiceConnector{query: query})
	t.Cleanup(func() { _ = db.Close() })
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return New(gormDB)
}

func TestFindPublicServiceByID(t *testing.T) {
	id := uuid.New()
	queryCount := 0
	repo := testPublicServiceRepository(t, func(query string, args []driver.NamedValue) (driver.Rows, error) {
		queryCount++
		if !strings.Contains(query, `"public_services"`) || strings.Contains(query, `"service_requirements"`) {
			t.Fatalf("unexpected query: %s", query)
		}
		if len(args) != 2 || args[0].Value != id.String() || args[1].Value != int64(1) {
			t.Fatalf("unexpected args: %#v", args)
		}
		return &publicServiceRows{values: [][]driver.Value{{id.String(), string(models.ServiceStatusActive)}}}, nil
	})

	service, err := repo.FindPublicServiceByID(context.Background(), id)
	if err != nil || service == nil || service.ID != id || service.Status != models.ServiceStatusActive {
		t.Fatalf("FindPublicServiceByID = (%#v, %v)", service, err)
	}
	if queryCount != 1 {
		t.Fatalf("got %d queries, want 1", queryCount)
	}
}

func TestFindPublicServiceByIDMissingAndDatabaseError(t *testing.T) {
	id := uuid.New()
	repo := testPublicServiceRepository(t, func(string, []driver.NamedValue) (driver.Rows, error) {
		return &publicServiceRows{}, nil
	})
	if _, err := repo.FindPublicServiceByID(context.Background(), id); !errors.Is(err, exceptions.ErrPublicServiceNotFound) {
		t.Fatalf("missing service error = %v", err)
	}
	if _, err := repo.FindPublicServiceByID(context.Background(), uuid.Nil); !errors.Is(err, exceptions.ErrPublicServiceNotFound) {
		t.Fatalf("nil ID error = %v", err)
	}

	dbErr := errors.New("database unavailable")
	repo = testPublicServiceRepository(t, func(string, []driver.NamedValue) (driver.Rows, error) {
		return nil, dbErr
	})
	if _, err := repo.FindPublicServiceByID(context.Background(), id); !errors.Is(err, dbErr) {
		t.Fatalf("database error = %v", err)
	}
}
