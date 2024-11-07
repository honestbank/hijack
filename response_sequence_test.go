package hijack_test

import (
	"testing"

	"github.com/honestbank/hijack/v2"

	"github.com/stretchr/testify/assert"
)

func TestSequence(t *testing.T) {
	t.Run("can return sequence of items", func(t *testing.T) {
		seq := hijack.Sequence(1, 2, 3, 4, 5)
		assert.Equal(t, 1, seq.GetNext())
		assert.Equal(t, 2, seq.GetNext())
		assert.Equal(t, 3, seq.GetNext())
		assert.Equal(t, 4, seq.GetNext())
		assert.Equal(t, 5, seq.GetNext())
	})
	t.Run("panics when it runs out of items in the sequence", func(t *testing.T) {
		seq := hijack.Sequence(1, 2, 3, 4, 5)
		assert.NotPanics(t, func() { seq.GetNext() })
		assert.NotPanics(t, func() { seq.GetNext() })
		assert.NotPanics(t, func() { seq.GetNext() })
		assert.NotPanics(t, func() { seq.GetNext() })
		assert.NotPanics(t, func() { seq.GetNext() })
		assert.Panics(t, func() {
			seq.GetNext()
		})
	})
}
