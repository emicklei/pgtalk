package main

import "testing"

func Test_abbreviate(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"single char", args{"a"}, "a"},
		{"two chars", args{"ab"}, "a"},
		{"split chars", args{"a_b"}, "ab"},
		{"split words", args{"ab_cd_ef"}, "ace"},
		{"dot split words", args{"ab_c.d_ef"}, "acde"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := abbreviate(tt.args.s); got != tt.want {
				t.Errorf("abbreviate() = %v, want %v", got, tt.want)
			}
		})
	}
}
