package util

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestContains(t *testing.T) {
	s := []string{
		"a",
		"b",
		"c",
	}

	v1 := Contains(s, "c")
	assert.Truef(t, v1, "test contain should contain c ")

	v2 := Contains(s, "d")
	assert.Falsef(t, v2, "test contain should not contain d")
}

func TestDistinctInt(t *testing.T) {
	s := []int{
		1,
		2,
		5,
		2,
		3,
		5,
		4,
		2,
		4,
		6,
	}

	ss := Distinct(s)
	assert.Equal(t, 6, len(ss), "count not equal 6")
	sort.Ints(s)

	assert.Equal(t, 1, ss[0], "the first is not 1")
}

func TestDistinctString(t *testing.T) {
	s := []string{
		"abc",
		"def",
		"qwe",
		"abc",
	}

	ss := Distinct(s)
	assert.Equal(t, 3, len(ss), "count not equal 3")
	sort.Strings(s)
	assert.Equal(t, "abc", ss[0], "the first is not abc")
}
