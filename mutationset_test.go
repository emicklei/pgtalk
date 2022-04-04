package pgtalk

import "testing"

func TestMutationSet_Insert(t *testing.T) {
	m := MakeMutationSet[poly](polyTable, polyColumns, MutationInsert)
	if got, want := oneliner(SQL(m)), "INSERT INTO public.polies (ftime,ffloat) VALUES ($1,$2)"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	{
		m := MakeMutationSet[poly](polyTable, polyColumns, MutationDelete)
		if got, want := oneliner(SQL(m)), "DELETE FROM public.polies p1"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
	{
		m := MakeMutationSet[poly](polyTable, polyColumns, MutationDelete)
		m = m.Where(polyFFloat.Equals(42.0))
		if got, want := oneliner(SQL(m)), "DELETE FROM public.polies p1 WHERE (p1.ffloat = 42)"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
	{
		toUpdate := []ColumnAccessor{}
		toUpdate = append(toUpdate, polyFFloat.Set(42.0))
		m := MakeMutationSet[poly](polyTable, toUpdate, MutationUpdate)
		if got, want := oneliner(SQL(m)), "UPDATE public.polies p1 SET ffloat = $1"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
	{
		toUpdate := []ColumnAccessor{}
		toUpdate = append(toUpdate, polyFFloat.Set(42.0))
		m := MakeMutationSet[poly](polyTable, toUpdate, MutationUpdate)
		m = m.Where(polyFFloat.Equals(12.0))
		if got, want := oneliner(SQL(m)), "UPDATE public.polies p1 SET ffloat = $1 WHERE (p1.ffloat = 12)"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
}

func TestInsertWithQueryArgument(t *testing.T) {
	toUpdate := []ColumnAccessor{}
	toUpdate = append(toUpdate, polyFFloat.Set(42.0))

	ps := NewParameterSet()
	v12 := ps.NewParameter(12.0)

	m := MakeMutationSet[poly](polyTable, toUpdate, MutationUpdate).Where(polyFFloat.Equals(v12))
	if got, want := oneliner(SQL(m)), "UPDATE public.polies p1 SET ffloat = $1 WHERE (p1.ffloat = ?)"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
