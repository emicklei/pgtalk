package main

import (
	"testing"
)

func TestTableFilter_IncludeAll(t *testing.T) {
	f := NewTableFilter(".*", "")
	if got, want := f.Includes("some"), true; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTableFilter_Mix(t *testing.T) {
	f := NewTableFilter("a.*,br.*,help,xxxxde.*", "ar.*")
	if got, want := f.Includes("alfred"), true; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := f.Includes("bruno"), true; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := f.Includes("arvind"), false; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := f.Includes("helpme"), false; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := f.Includes("help"), true; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := !f.Includes("deno"), true; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
