package async

import (
	"testing"

	"github.com/reugn/async/internal"
)

func TestGoroutineID(t *testing.T) {
	gid, err := GoroutineID()

	internal.AssertEqual(t, nil, err)
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
