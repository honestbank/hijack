package generics_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/hijack/v2/internal/generics"
)

func TestFilter(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7}
	result := generics.Filter(items, func(item int) bool {
		return item%2 == 0
	})
	assert.Equal(t, []int{2, 4, 6}, result)
}
