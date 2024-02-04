package async

import (
	"testing"

	"github.com/reugn/async/internal/assert"
)

func TestGoroutineID(t *testing.T) {
	gid, err := GoroutineID()

	assert.IsNil(t, err)
	t.Log(gid)
}

func BenchmarkGetGroutineID3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GoroutineID()
		if err != nil {
			b.Error("failed to get gid")
		}
	}
}
