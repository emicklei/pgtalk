package pgtalk

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestResultIterator_Next(t *testing.T) {
	t.Run("happy flow", func(t *testing.T) {
		rows := &mockRows{
			nextResult:  true,
			scanError:   nil,
			closeCalled: false,
		}
		selectors := []ColumnAccessor{
			mockColumnAccessor{
				fieldValueToScan: func(entity any) any {
					return &entity.(*testEntity).ID
				},
			},
		}
		i := &resultIterator[testEntity]{
			rows:             rows,
			orderedSelectors: selectors,
		}
		if !i.HasNext() {
			t.Fatal("expected next")
		}
		entity, err := i.Next()
		if err != nil {
			t.Fatal(err)
		}
		if entity.ID != "test-id" {
			t.Errorf("expected test-id, got %s", entity.ID)
		}
		if i.HasNext() {
			t.Fatal("unexpected next")
		}
		if !rows.closeCalled {
			t.Error("expected close to be called")
		}
	})
}

func TestResultIterator_Err(t *testing.T) {
	t.Run("query error", func(t *testing.T) {
		i := &resultIterator[testEntity]{
			queryError: errors.New("query error"),
		}
		if err := i.Err(); err == nil || err.Error() != "query error" {
			t.Errorf("expected query error, got %v", err)
		}
	})
	t.Run("rows error", func(t *testing.T) {
		rows := &mockRows{
			err: errors.New("rows error"),
		}
		i := &resultIterator[testEntity]{
			rows: rows,
		}
		if err := i.Err(); err == nil || err.Error() != "rows error" {
			t.Errorf("expected rows error, got %v", err)
		}
	})
}

func TestResultIterator_GetParams(t *testing.T) {
	params := []any{"param1", 2}
	i := &resultIterator[testEntity]{
		params: params,
	}
	p := i.GetParams()
	if len(p) != 2 {
		t.Errorf("expected 2 params, got %d", len(p))
	}
	if p[1] != "param1" {
		t.Errorf("expected param1, got %v", p[1])
	}
	if p[2] != 2 {
		t.Errorf("expected 2, got %v", p[2])
	}
}

// mockColumnAccessor is a mock for the ColumnAccessor interface
type mockColumnAccessor struct {
	fieldValueToScan func(entity any) any
}

func (m mockColumnAccessor) SQLOn(w WriteContext)                {}
func (m mockColumnAccessor) Name() string                        { return "" }
func (m mockColumnAccessor) ValueToInsert() any                  { return nil }
func (m mockColumnAccessor) Column() ColumnInfo                  { return ColumnInfo{} }
func (m mockColumnAccessor) FieldValueToScan(entity any) any     { return m.fieldValueToScan(entity) }
func (m mockColumnAccessor) AppendScannable(list []any) []any    { return list }
func (m mockColumnAccessor) Get(values map[string]any) any       { return nil }
func (m mockColumnAccessor) SetSource(parameterIndex int) string { return "" }

// mockRows is a mock for the pgx.Rows interface
type mockRows struct {
	nextResult  bool
	scanError   error
	closeCalled bool
	err         error
}

func (m *mockRows) Close()                                       { m.closeCalled = true }
func (m *mockRows) Err() error                                   { return m.err }
func (m *mockRows) CommandTag() pgconn.CommandTag                { return pgconn.CommandTag{} }
func (m *mockRows) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (m *mockRows) Next() bool                                   { return m.nextResult }
func (m *mockRows) Scan(dest ...any) error {
	if m.scanError != nil {
		return m.scanError
	}
	// simulate scanning a value
	if len(dest) > 0 {
		if id, ok := dest[0].(*string); ok {
			*id = "test-id"
		}
	}
	m.nextResult = false // only one row
	return nil
}
func (m *mockRows) RawValues() [][]byte    { return nil }
func (m *mockRows) Conn() *pgx.Conn        { return nil }
func (m *mockRows) Values() ([]any, error) { return nil, nil }

// testEntity is a simple struct for testing
type testEntity struct {
	ID string
}
