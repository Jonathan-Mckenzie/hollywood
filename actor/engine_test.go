package actor

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummy struct{}

func newDummy() Receiver {
	return &dummy{}
}
func (d *dummy) Receive(_ *Context) {}

func TestSpawn(t *testing.T) {
	e := NewEngine()
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			tag := strconv.Itoa(i)
			pid := e.Spawn(newDummy, "dummy", tag)
			e.Send(pid, 1)
			wg.Done()
		}(i)
	}

	wg.Wait()
	assert.Equal(t, e.registry.lookup.Count(), 1)
	addrs, ok := e.registry.lookup.Get("local")
	assert.True(t, ok)
	assert.Equal(t, 10, addrs.Count())
}

func TestSpawnPID(t *testing.T) {
	e := NewEngine()

	pid := e.Spawn(newDummy, "dummy", "1")
	assert.Equal(t, "local/dummy/1", pid.String())
	assert.Equal(t, e.registry.lookup.Count(), 1)
}

func TestPoison(t *testing.T) {
	e := NewEngine()

	for i := 0; i < 4; i++ {
		tag := strconv.Itoa(i)
		pid := e.Spawn(newDummy, "dummy", tag)
		e.Poison(pid)
	}

	assert.Equal(t, e.registry.lookup.Count(), 0)
}

func BenchmarkXxx(b *testing.B) {
	e := NewEngine()
	pid := e.Spawn(newDummy, "dummy")

	for i := 0; i < b.N; i++ {
		e.Send(pid, pid)
	}
}
