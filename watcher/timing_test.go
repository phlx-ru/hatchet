package watcher

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewEmptyTiming(t *testing.T) {
	timing := NewEmptyTiming()
	time.Sleep(20 * time.Millisecond)
	timing.Send("test.timings")

	history := timing.History()
	item := history[0]
	require.Equal(t, "test.timings", item.Bucket)
	require.GreaterOrEqual(t, item.Duration, 20*time.Millisecond)
	require.Less(t, item.Duration, 50*time.Millisecond)
}
