package pgtalk

import "testing"

func TestMutationSet_Insert(t *testing.T) {
	m := MakeMutationSet[poly](polyAccess, polyColumns, MutationInsert)
	if got, want := SQL(m), "INSERT INTO public.polies (ftime,ffloat) VALUES ($1,$2)"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	{
		m := MakeMutationSet[poly](polyAccess, polyColumns, MutationDelete)
		if got, want := SQL(m), "DELETE FROM public.polies p1"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
	{
		m := MakeMutationSet[poly](polyAccess, polyColumns, MutationDelete)
		m = m.Where(polyFFloat.Equals(42.0))
		if got, want := SQL(m), "DELETE FROM public.polies p1 WHERE (p1.ffloat = 42)"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
	{
		toUpdate := []ColumnAccessor{}
		toUpdate = append(toUpdate, polyFFloat.Set(42.0))
		m := MakeMutationSet[poly](polyAccess, toUpdate, MutationUpdate)
		if got, want := SQL(m), "UPDATE public.polies p1 SET ffloat = $1"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
	{
		toUpdate := []ColumnAccessor{}
		toUpdate = append(toUpdate, polyFFloat.Set(42.0))
		m := MakeMutationSet[poly](polyAccess, toUpdate, MutationUpdate)
		m = m.Where(polyFFloat.Equals(12.0))
		if got, want := SQL(m), "UPDATE public.polies p1 SET ffloat = $1 WHERE (p1.ffloat = 12)"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
}
