package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDistinct(t *testing.T) {
	t.Run("Integer Slice", func(t *testing.T) {
		input := []int{1, 2, 5, 2, 3, 5, 4, 2, 4, 6}
		expected := []int{1, 2, 5, 3, 4, 6}
		result := Distinct(input)
		assert.Equal(t, expected, result, "Distinct for integers should remove duplicates")
	})

	t.Run("String Slice", func(t *testing.T) {
		input := []string{"abc", "def", "qwe", "abc"}
		expected := []string{"abc", "def", "qwe"}
		result := Distinct(input)
		assert.Equal(t, expected, result, "Distinct for strings should remove duplicates")
	})

	t.Run("Empty Slice", func(t *testing.T) {
		var input []int
		result := Distinct(input)
		assert.Empty(t, result, "Distinct for an empty slice should return an empty slice")
	})

	t.Run("Slice with no duplicates", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		result := Distinct(input)
		assert.Equal(t, input, result, "Distinct for a slice with no duplicates should return the original slice")
	})
}

func TestFirstOrDefault(t *testing.T) {
	t.Run("Slice with elements", func(t *testing.T) {
		input := []string{"a", "b", "c"}
		val, ok := FirstOrDefault(input)
		assert.True(t, ok, "FirstOrDefault should return true for a non-empty slice")
		assert.Equal(t, "a", val, "FirstOrDefault should return the first element")
	})

	t.Run("Empty slice", func(t *testing.T) {
		var input []int
		val, ok := FirstOrDefault(input)
		assert.False(t, ok, "FirstOrDefault should return false for an empty slice")
		assert.Zero(t, val, "FirstOrDefault should return the zero value for an empty slice")
	})
}

func TestLastOrDefault(t *testing.T) {
	t.Run("Slice with elements", func(t *testing.T) {
		input := []string{"a", "b", "c"}
		val, ok := LastOrDefault(input)
		assert.True(t, ok, "LastOrDefault should return true for a non-empty slice")
		assert.Equal(t, "c", val, "LastOrDefault should return the last element")
	})

	t.Run("Empty slice", func(t *testing.T) {
		var input []string
		val, ok := LastOrDefault(input)
		assert.False(t, ok, "LastOrDefault should return false for an empty slice")
		assert.Equal(t, "", val, "LastOrDefault should return the zero value for an empty slice")
	})
}