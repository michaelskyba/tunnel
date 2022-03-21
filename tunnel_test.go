package main

import (
	"testing"
)

// Camel case forced
func TestNewCard(tests *testing.T) {
	type test struct {
		Input string
		Output string
	}

	cases := []test{
		test{"", ""},
		test{"foo", "foo"},
		test{"	", "		0	2.5	0	1617249600"},
		test{"foo	", "foo		0	2.5	0	1617249600"},
		test{"	bar", "	bar	0	2.5	0	1617249600"},
		test{"		", "		"},
		test{"		baz", "		baz"},
		test{"			", "			"},
		test{"foo	bar", "foo	bar	0	2.5	0	1617249600"},
		test{"foo	bar	baz", "foo	bar	baz"},
		test{"foo	bar	0	2.5	0	1617249600", "foo	bar	0	2.5	0	1617249600"},
		test{"foo	bar	1	2.62	3	1890795600", "foo	bar	1	2.62	3	1890795600"}}

	// Iterate over test cases
	var got string
	for _, test_case := range cases {
		got = new_card(test_case.Input)
		if got != test_case.Output {
			tests.Errorf("\nnew_card() Input:\n%v\n\nOutput\n%v\n\nCorrect output\n%v", test_case.Input, got, test_case.Output)
		}
	}
}
