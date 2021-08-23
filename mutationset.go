package pgtalk

import (
	"context"
	"fmt"
	"io"

	"github.com/jackc/pgx/v4"
)

const (
	MutationDelete = iota
	MutationInsert
	MutationUpdate
)

type MutationSet[T any] struct {
	tableAccess   TableAccessor
	selectors     []ColumnAccessor
	condition     SQLExpression
	returning     []ColumnAccessor
	operationType int
}

func MakeMutationSet[T any](tableAccess TableAccessor, selectors []ColumnAccessor, operationType int) MutationSet[T] {
	if assertEnabled {
		assertEachAccessorHasTableInfo(selectors, tableAccess.TableInfo)
	}
	return MutationSet[T]{
		tableAccess:   tableAccess,
		selectors:     selectors,
		condition:     EmptyCondition,
		operationType: operationType}
}

// SQL returns the full SQL mutation query
func (m MutationSet[T]) SQLOn(w io.Writer) {
	if m.operationType == MutationInsert {
		fmt.Fprint(w, "INSERT INTO ")
		fmt.Fprintf(w, "%s.%s", m.tableAccess.Schema, m.tableAccess.Name)
		fmt.Fprint(w, " (")
		m.columnsSectionOn(m.selectors, w)
		fmt.Fprint(w, ") VALUES (")
		m.valuesSectionOn(w)
		fmt.Fprint(w, ")")
		if len(m.returning) > 0 {
			fmt.Fprint(w, " RETURNING ")
			m.columnsSectionOn(m.returning, w)
		}
		return
	}
	if m.operationType == MutationDelete {
		fmt.Fprint(w, "DELETE FROM ")
		m.tableAccess.SQLOn(w)
		fmt.Fprint(w, " WHERE ")
		m.condition.SQLOn(w)
		if len(m.returning) > 0 {
			fmt.Fprint(w, " RETURNING ")
			m.columnsSectionOn(m.returning, w)
		}
		return
	}
	if m.operationType == MutationUpdate {
		fmt.Fprint(w, "UPDATE ")
		m.tableAccess.SQLOn(w)
		fmt.Fprint(w, " SET ")
		m.setSectionOn(w)
		fmt.Fprint(w, " WHERE ")
		m.condition.SQLOn(w)
		if len(m.returning) > 0 {
			fmt.Fprint(w, " RETURNING ")
			m.columnsSectionOn(m.returning, w)
		}
		return
	}
}

func (m MutationSet[T]) Where(condition SQLExpression) MutationSet[T] {
	if assertEnabled {
		access := condition.Collect([]ColumnAccessor{})
		assertEachAccessorIn(access, m.tableAccess.AllColumns)
	}
	m.condition = condition
	return m
}

func (m MutationSet[T]) Returning(columns ...ColumnAccessor) MutationSet[T] {
	if assertEnabled {
		assertEachAccessorIn(columns, m.tableAccess.AllColumns)
	}
	m.returning = columns
	return m
}

// todo
func (m MutationSet[T]) On() MutationSet[T] {
	return m
}

// Pre: must be run inside transaction
func (m MutationSet[T]) Exec(ctx context.Context, conn *pgx.Conn) *ResultIterator {
	args := []interface{}{}
	for _, each := range m.selectors {
		args = append(args, each.ValueToInsert())
	}
	rows, err := conn.Query(ctx, SQL(m), args...)
	if err == nil && !m.canProduceResults() {
		rows.Close()
	}
	return &ResultIterator{queryError: err, rows: rows, selectors: m.returning}
}

func (m MutationSet[T]) canProduceResults() bool {
	return len(m.returning) > 0
}

func (m MutationSet[T]) columnsSectionOn(which []ColumnAccessor, buf io.Writer) {
	for i, each := range which {
		if i > 0 {
			io.WriteString(buf, ",")
		}
		io.WriteString(buf, each.Name())
	}
}

func (m MutationSet[T]) valuesSectionOn(buf io.Writer) {
	for i := range m.selectors {
		if i > 0 {
			io.WriteString(buf, ",")
		}
		fmt.Fprintf(buf, "$%d", i+1)
	}
}

func (m MutationSet[T]) setSectionOn(w io.Writer) {
	for i, each := range m.selectors {
		if i > 0 {
			io.WriteString(w, ",")
		}
		fmt.Fprintf(w, "%s = $%d", each.Name(), i+1)
	}
}
