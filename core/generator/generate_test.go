package generator

import (
	"testing"
)

type TestCase struct {
	input    Dim
	expected []ClueGroup
}

func TestFind(t *testing.T) {

	tc := []TestCase{
		TestCase{
			input: Dim{1, 2, 3, 4, 5, 6, 7, 8, 9},
			expected: []ClueGroup{
				ClueGroup{Start: 0, Length: 9, Offset: 1},
			},
		},
		TestCase{
			input: Dim{1, 2, 3, 0, 5, 6, 7, 8, 9},
			expected: []ClueGroup{
				ClueGroup{Start: 0, Length: 3, Offset: 1},
				ClueGroup{Start: 4, Length: 5, Offset: 1},
			},
		},
		TestCase{
			input: Dim{0, 2, 3, 0, 5, 6, 7, 8, 9},
			expected: []ClueGroup{
				ClueGroup{Start: 1, Length: 2, Offset: 1},
				ClueGroup{Start: 4, Length: 5, Offset: 1},
			},
		},
		TestCase{
			input: Dim{0, 2, 3, 0, 5, 6, 7, 8, 0},
			expected: []ClueGroup{
				ClueGroup{Start: 1, Length: 2, Offset: 1},
				ClueGroup{Start: 4, Length: 4, Offset: 1},
			},
		},
		TestCase{
			input: Dim{0, 2, 3, 0, 5, 6, 7, 0, 0},
			expected: []ClueGroup{
				ClueGroup{Start: 1, Length: 2, Offset: 1},
				ClueGroup{Start: 4, Length: 3, Offset: 1},
			},
		},
		TestCase{
			input: Dim{0, 2, 3, 4, 0, 0, 0, 0, 0},
			expected: []ClueGroup{
				ClueGroup{Start: 1, Length: 3, Offset: 1},
			},
		},
		TestCase{
			input: Dim{0, 2, 3, 0, 5, 0, 0, 8, 0},
			expected: []ClueGroup{
				ClueGroup{Start: 1, Length: 2, Offset: 1},
				ClueGroup{Start: 4, Length: 1, Offset: 1},
				ClueGroup{Start: 7, Length: 1, Offset: 1},
			},
		},
	}

	for _, v := range tc {
		actual := scan_rows(v.input[:])

		if len(actual) != len(v.expected) {
			t.Errorf("incorrect length: expect: %v, actual: %v", len(v.expected), len(actual))
		}

		for _, ec := range v.expected {

			found := false
			for _, ac := range actual {
				if IsEqual(&ac, &ec) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected object not found: expected: %v, actual: %v", ec, actual)
			}
		}

	}

}

func TestColGroups(t *testing.T) {

	/* test_cases := []struct {
		tc: Grid
		expected: []ClueGroup
	}{
		input: Grid{0, 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0} ,
		expected: []ClueGroup{},

	}
	*/
	test_cases := []struct {
		input    Grid
		expected []ClueGroup
	}{
		{
			input: Grid{
				0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0},
			expected: []ClueGroup{},
		},
		{
			input: Grid{
				1, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0},
			expected: []ClueGroup{
				ClueGroup{Start: 0, Length: 1},
			},
		},
		{
			input: Grid{
				0, 1, 0, 0, 0, 0, 0, 1, 0,
				0, 2, 0, 1, 0, 1, 0, 2, 0,
				0, 3, 0, 0, 0, 2, 0, 3, 0,
				0, 0, 0, 2, 0, 0, 0, 4, 0,
				0, 0, 0, 0, 0, 0, 0, 5, 0,
				0, 4, 0, 3, 0, 0, 0, 6, 0,
				0, 5, 0, 0, 0, 0, 0, 7, 0,
				0, 0, 0, 4, 0, 8, 0, 8, 0,
				0, 0, 0, 0, 0, 9, 0, 9, 1},
			expected: []ClueGroup{
				ClueGroup{Start: 1, Length: 3},
				ClueGroup{Start: 46, Length: 2},
				ClueGroup{Start: 12, Length: 1},
				ClueGroup{Start: 30, Length: 1},
				ClueGroup{Start: 48, Length: 1},
				ClueGroup{Start: 66, Length: 1},
				ClueGroup{Start: 14, Length: 2},
				ClueGroup{Start: 68, Length: 2},
				ClueGroup{Start: 7, Length: 9},
				ClueGroup{Start: 80, Length: 1},
			},
		},
	}

	for _, v := range test_cases {
		actual := scan_columns(v.input[:])

		if len(actual) != len(v.expected) {
			t.Errorf("incorrect length: expect: %v, actual: %v", v.expected, actual)
		}
	}
}
