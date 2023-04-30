package main

import "testing"

func TestNewCard(tests *testing.T) {
	type test struct {
		Input  string
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
	for _, testCase := range cases {
		got = newCard(testCase.Input)
		if got != testCase.Output {
			tests.Errorf("\nnewCard() Input:\n%v\n\nOutput\n%v\n\nCorrect output\n%v", testCase.Input, got, testCase.Output)
		}
	}
}
