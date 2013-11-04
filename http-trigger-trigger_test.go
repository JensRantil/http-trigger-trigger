package main

import (
	"testing"
)

type escapeSplitTest struct {
	s string
	expected []string
}

func TestSplitEscapedString(t *testing.T) {
	tests := []escapeSplitTest{
		escapeSplitTest{
			"hej",
			[]string{"hej"},
		},
		escapeSplitTest{
			"hej bah",
			[]string{"hej", "bah"},
		},
		escapeSplitTest{
			`hej\ bah`,
			[]string{"hej bah"},
		},
		escapeSplitTest{
			`hej\ bah tjing`,
			[]string{"hej bah", "tjing"},
		},
	}
	for _, example := range tests {
		res := splitEscapedString(example.s)
		if len(res) != len(example.expected) {
			t.Error("Not same size.")
			t.Error("String:", example.s)
			t.Error("Was:", res)
			t.Error("Expected:", example.expected)
			continue
		}
		for i, piece := range res {
			if piece != example.expected[i] {
				t.Error("String:", example.s)
				t.Error("Was:", piece, len(piece))
				ts := example.expected[i]
				t.Error("Expected:", ts, len(ts))
			}
		}
	}
}
