package pgtalk

import (
	"context"
	"fmt"
	"io"
)

const (
	MutationDelete = iota
	MutationInsert
	MutationUpdate
)

type MutationSet[T any] struct {
	tableInfo     TableInfo
	selectors     []ColumnAccessor
	condition     SQLExpression
	returning     []ColumnAccessor
	operationType int
}

func MakeMutationSet[T any](tableInfo TableInfo, selectors []ColumnAccessor, operationType int) MutationSet[T] {
	return MutationSet[T]{
		tableInfo:     tableInfo,
		selectors:     selectors,
		condition:     EmptyCondition,
		operationType: operationType}
}

// SQLOn returns the full SQL mutation query
func (m MutationSet[T]) SQLOn(w writeContext) {
	if m.operationType == MutationInsert {
		fmt.Fprint(w, "INSERT INTO ")
		fmt.Fprintf(w, "%s.%s", m.tableInfo.Schema, m.tableInfo.Name)
		fmt.Fprint(w, " (")
		m.columnsSectionOn(m.selectors, w)
		fmt.Fprint(w, ")\nVALUES (")
		m.valuesSectionOn(w)
		fmt.Fprint(w, ")")
		if len(m.returning) > 0 {
			fmt.Fprint(w, "\nRETURNING ")
			m.columnsSectionOn(m.returning, w)
		}
		return
	}
	if m.operationType == MutationDelete {
		fmt.Fprint(w, "DELETE FROM\n")
		m.tableInfo.SQLOn(w)
		if m.condition != EmptyCondition {
			fmt.Fprint(w, "\nWHERE ")
			m.condition.SQLOn(w)
		}
		if len(m.returning) > 0 {
			fmt.Fprint(w, "\nRETURNING ")
			m.columnsSectionOn(m.returning, w)
		}
		return
	}
	if m.operationType == MutationUpdate {
		fmt.Fprint(w, "UPDATE ")
		m.tableInfo.SQLOn(w)
		fmt.Fprint(w, "\nSET ")
		m.setSectionOn(w)
		if m.condition != EmptyCondition {
			fmt.Fprint(w, "\nWHERE ")
		}
		m.condition.SQLOn(w)
		if len(m.returning) > 0 {
			fmt.Fprint(w, "\nRETURNING ")
			m.columnsSectionOn(m.returning, w)
		}
		return
	}
}

func (m MutationSet[T]) Where(condition SQLExpression) MutationSet[T] {
	m.condition = condition
	return m
}

func (m MutationSet[T]) Returning(columns ...ColumnAccessor) MutationSet[T] {
	m.returning = columns
	return m
}

// todo
func (m MutationSet[T]) On() MutationSet[T] {
	return m
}

// Pre: must be run inside transaction
func (m MutationSet[T]) Exec(ctx context.Context, conn querier) *resultIterator[T] {
	query := SQL(m)
	rows, err := conn.Query(ctx, query, m.ValuesToInsert()...)
	if err == nil && !m.canProduceResults() {
		rows.Close()
	}
	return &resultIterator[T]{queryError: err, rows: rows, selectors: m.returning}
}

// ValuesToInsert returns the parameters values for the mutation query.
func (m MutationSet[T]) ValuesToInsert() []any {
	args := make([]any, len(m.selectors))
	for i, each := range m.selectors {
		args[i] = each.ValueToInsert()
	}
	return args
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
