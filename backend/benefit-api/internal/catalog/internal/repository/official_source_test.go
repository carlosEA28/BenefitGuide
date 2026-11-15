package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/carlosEA28/benefitGuide/internal/catalog/internal/models"
	"github.com/carlosEA28/benefitGuide/internal/shared/exceptions"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type sourceConnector struct {
	query func(string, []driver.NamedValue) (driver.Rows, error)
	exec  func(string, []driver.NamedValue) (driver.Result, error)
}

func (c sourceConnector) Connect(context.Context) (driver.Conn, error) { return sourceConn(c), nil }
func (sourceConnector) Driver() driver.Driver                          { return sourceDriver{} }

type sourceDriver struct{}

func (sourceDriver) Open(string) (driver.Conn, error) { return nil, errors.New("use connector") }

type sourceConn sourceConnector

func (sourceConn) Prepare(string) (driver.Stmt, error) { return nil, errors.New("unexpected prepare") }
func (sourceConn) Close() error                        { return nil }
func (sourceConn) Begin() (driver.Tx, error)           { return nil, errors.New("unexpected transaction") }
func (c sourceConn) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return c.query(query, args)
}
func (c sourceConn) ExecContext(_ context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	return c.exec(query, args)
}

type sourceRows struct {
	columns []string
	values  [][]driver.Value
	index   int
}

func (r *sourceRows) Columns() []string { return r.columns }
func (r *sourceRows) Close() error      { return nil }
func (r *sourceRows) Next(dest []driver.Value) error {
	if r.index == len(r.values) {
		return io.EOF
	}
	copy(dest, r.values[r.index])
	r.index++
	return nil
}

func newSourceRepository(t *testing.T, connector sourceConnector) OfficialSourceRepository {
	t.Helper()
	db := sql.OpenDB(connector)
	t.Cleanup(func() { _ = db.Close() })
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db, PreferSimpleProtocol: true}), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		t.Fatal(err)
	}
	return NewOfficialSourceRepository(gormDB)
}

func TestOfficialSourceUpdateChangesOnlyEditableFields(t *testing.T) {
	id, publicServiceID := uuid.New(), uuid.New()
	called := false
	repo := newSourceRepository(t, sourceConnector{exec: func(query string, args []driver.NamedValue) (driver.Result, error) {
		called = true
		for _, column := range []string{`"title"`, `"url"`, `"agency"`, `"last_verified_at"`} {
			if !strings.Contains(query, column+"=") {
				t.Fatalf("missing editable column %s in %s", column, query)
			}
		}
		for _, column := range []string{`"public_service_id"=`, `"created_at"=`, `"id"=`} {
			if strings.Contains(query, column) {
				t.Fatalf("immutable column %s updated in %s", column, query)
			}
		}
		if !strings.Contains(query, `WHERE "id" =`) {
			t.Fatalf("update has no ID predicate: %s", query)
		}
		if len(args) != 5 || args[3].Value != nil || args[4].Value != id.String() {
			t.Fatalf("unexpected update args: %#v", args)
		}
		return driver.RowsAffected(1), nil
	}})
	source := &models.OfficialSource{ID: id, PublicServiceID: publicServiceID, Title: "Nova fonte", URL: "https://gov.br", Agency: "Órgão", CreatedAt: time.Now()}
	got, err := repo.Update(context.Background(), source)
	if err != nil || got != source || !called {
		t.Fatalf("Update = (%#v, %v), called = %v", got, err, called)
	}
}

func TestOfficialSourceCreateReturnsDatabaseID(t *testing.T) {
	id, publicServiceID := uuid.New(), uuid.New()
	repo := newSourceRepository(t, sourceConnector{query: func(query string, args []driver.NamedValue) (driver.Rows, error) {
		if !strings.Contains(query, `INSERT INTO "official_sources"`) || strings.Contains(query, `"id"`) && !strings.Contains(query, `RETURNING "id"`) {
			t.Fatalf("unexpected create query: %s", query)
		}
		return &sourceRows{columns: []string{"id"}, values: [][]driver.Value{{id.String()}}}, nil
	}})
	source := &models.OfficialSource{PublicServiceID: publicServiceID, Title: "Fonte", URL: "https://gov.br", Agency: "Órgão"}
	got, err := repo.Create(context.Background(), source)
	if err != nil || got != source || got.ID != id {
		t.Fatalf("Create = (%#v, %v)", got, err)
	}
}

func TestOfficialSourceFindByID(t *testing.T) {
	id, publicServiceID := uuid.New(), uuid.New()
	repo := newSourceRepository(t, sourceConnector{query: func(query string, args []driver.NamedValue) (driver.Rows, error) {
		if !strings.Contains(query, `FROM "official_sources"`) || len(args) != 2 || args[0].Value != id.String() {
			t.Fatalf("unexpected find query %s with args %#v", query, args)
		}
		return &sourceRows{
			columns: []string{"id", "public_service_id", "title", "url", "agency"},
			values:  [][]driver.Value{{id.String(), publicServiceID.String(), "Fonte", "https://gov.br", "Órgão"}},
		}, nil
	}})
	source, err := repo.FindByID(context.Background(), id)
	if err != nil || source == nil || source.ID != id || source.PublicServiceID != publicServiceID {
		t.Fatalf("FindByID = (%#v, %v)", source, err)
	}
}

func TestOfficialSourceMissingIDs(t *testing.T) {
	repo := newSourceRepository(t, sourceConnector{exec: func(string, []driver.NamedValue) (driver.Result, error) {
		return driver.RowsAffected(0), nil
	}})
	if _, err := repo.FindByID(context.Background(), uuid.Nil); !errors.Is(err, exceptions.ErrOfficialSourceNotFound) {
		t.Fatalf("FindByID nil ID = %v", err)
	}
	if _, err := repo.Update(context.Background(), &models.OfficialSource{}); !errors.Is(err, exceptions.ErrOfficialSourceNotFound) {
		t.Fatalf("Update nil ID = %v", err)
	}
	if err := repo.Delete(context.Background(), uuid.Nil); !errors.Is(err, exceptions.ErrOfficialSourceNotFound) {
		t.Fatalf("Delete nil ID = %v", err)
	}
	if err := repo.Delete(context.Background(), uuid.New()); !errors.Is(err, exceptions.ErrOfficialSourceNotFound) {
		t.Fatalf("Delete missing ID = %v", err)
	}
}

func TestOfficialSourceDeleteExisting(t *testing.T) {
	id := uuid.New()
	repo := newSourceRepository(t, sourceConnector{exec: func(query string, args []driver.NamedValue) (driver.Result, error) {
		if !strings.Contains(query, `DELETE FROM "official_sources"`) || !strings.Contains(query, `WHERE id =`) || len(args) != 1 || args[0].Value != id.String() {
			t.Fatalf("unexpected delete query %s with args %#v", query, args)
		}
		return driver.RowsAffected(1), nil
	}})
	if err := repo.Delete(context.Background(), id); err != nil {
		t.Fatalf("Delete = %v", err)
	}
}

func TestOfficialSourceListUsesStableOrderAndEmptySlice(t *testing.T) {
	publicServiceID := uuid.New()
	repo := newSourceRepository(t, sourceConnector{query: func(query string, args []driver.NamedValue) (driver.Rows, error) {
		if !strings.Contains(query, `"official_sources"`) || !strings.Contains(query, `ORDER BY created_at ASC, id ASC`) {
			t.Fatalf("unexpected list query: %s", query)
		}
		if len(args) != 1 || args[0].Value != publicServiceID.String() {
			t.Fatalf("unexpected list args: %#v", args)
		}
		return &sourceRows{columns: []string{"id"}}, nil
	}})
	sources, err := repo.FindByPublicServiceID(context.Background(), publicServiceID)
	if err != nil || sources == nil || len(sources) != 0 {
		t.Fatalf("FindByPublicServiceID = (%#v, %v)", sources, err)
	}
}
