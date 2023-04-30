package main

import "testing"

func TestNewCard(tests *testing.T) {
	type test struct {
		Card        string
		CurrentTime int
		Output      string
	}

	cases := []test{
		{"", 1890795600, ""},
		{"foo", -5, "foo"},
		{"	", 1617249600, "		0	2.5	0	1617249600"},
		{"foo	", 6, "foo		0	2.5	0	6"},
		{"foo	", 0, "foo		0	2.5	0	0"},
		{"	bar", 1890795600, "	bar	0	2.5	0	1890795600"},
		{"		", 1890795600, "		"},
		{"		baz", 1890795600, "		baz"},
		{"			", 1890795600, "			"},
		{"foo	bar", 1617249600, "foo	bar	0	2.5	0	1617249600"},
		{"foo	bar", 1890795600, "foo	bar	0	2.5	0	1890795600"},
		{"foo	bar", 1682898064, "foo	bar	0	2.5	0	1682898064"},
		{"foo	bar", -14222, "foo	bar	0	2.5	0	-14222"},
		{"foo	bar	baz", 1890795600, "foo	bar	baz"},
		{"foo	bar	0	2.5	0	1617249600", 1890795600, "foo	bar	0	2.5	0	1617249600"},
		{"foo	bar	1	2.62	3	1890795600", 1231727323, "foo	bar	1	2.62	3	1890795600"},
	}

	// Iterate over test cases
	var got string
	for _, testCase := range cases {
		got = newCard(testCase.Card, testCase.CurrentTime)
		if got != testCase.Output {
			tests.Errorf("\nnewCard() Card:\n%v\nTime:\n%v\n\nOutput:\n%v\n\nCorrect output:\n%v", testCase.Card, testCase.CurrentTime, got, testCase.Output)
		}
	}
}
