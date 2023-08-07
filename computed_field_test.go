package pgtalk

import (
	"testing"
)

type hasExpressionResult struct {
	m map[string]any
}

func newHasExpressionResult() *hasExpressionResult {
	return &hasExpressionResult{make(map[string]any)}
}
func (r *hasExpressionResult) AddExpressionResult(key string, value any) {
	r.m[key] = value
}

func (r *hasExpressionResult) GetExpressionResult(key string) any {
	pv := r.m[key].(*any)
	return *pv
}

func TestFieldValueToScan(t *testing.T) {
	cf := new(computedField)
	cf.ResultName = "one"

	ent1 := newHasExpressionResult()
	ent2 := newHasExpressionResult()

	{
		scanAddr := cf.FieldValueToScan(ent1)
		// simulate scan
		var newValue any = "help1"
		pscanAddr := scanAddr.(*any)
		*pscanAddr = newValue
	}

	{
		scanAddr := cf.FieldValueToScan(ent2)
		// simulate scan
		var newValue any = "help2"
		pscanAddr := scanAddr.(*any)
		*pscanAddr = newValue
	}

	if got, want := ent1.GetExpressionResult("one"), "help1"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	if got, want := ent2.GetExpressionResult("one"), "help2"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}
