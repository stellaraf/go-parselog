package types_test

import (
	"testing"
	"time"

	"github.com/stellaraf/go-parselog/types"
	"github.com/stretchr/testify/assert"
)

func Test_Log(t *testing.T) {
	t.Run("isis is", func(t *testing.T) {
		t.Parallel()
		log1 := &types.ISISLog{Base: types.Base{Type: types.ISIS}}
		log2 := &types.ISISLog{Base: types.Base{Type: types.ISIS}}
		assert.True(t, log1.Is(log2))
	})
	t.Run("bgp is", func(t *testing.T) {
		t.Parallel()
		log1 := &types.BGPLog{Base: types.Base{Type: types.BGP}}
		log2 := &types.BGPLog{Base: types.Base{Type: types.BGP}}
		assert.True(t, log1.Is(log2))
	})
	t.Run("isis up", func(t *testing.T) {
		t.Parallel()
		log := &types.ISISLog{Base: types.Base{Type: types.ISIS}, State: types.UP}
		assert.True(t, log.Up())
	})
	t.Run("isis down", func(t *testing.T) {
		t.Parallel()
		log := &types.ISISLog{Base: types.Base{Type: types.ISIS}, State: types.DOWN}
		assert.True(t, log.Down())
	})
	t.Run("bgp up", func(t *testing.T) {
		t.Parallel()
		log := &types.BGPLog{Base: types.Base{Type: types.BGP}, State: types.UP}
		assert.True(t, log.Up())
	})
	t.Run("bgp down", func(t *testing.T) {
		t.Parallel()
		log := &types.BGPLog{Base: types.Base{Type: types.BGP}, State: types.DOWN}
		assert.True(t, log.Down())
	})
	t.Run("isis id", func(t *testing.T) {
		t.Parallel()
		log := &types.ISISLog{Base: types.Base{Type: types.ISIS}, Local: "local", Remote: "remote", Interface: "interface"}
		assert.NotEmpty(t, log.ID())
	})
	t.Run("bgp id", func(t *testing.T) {
		t.Parallel()
		log := &types.BGPLog{Base: types.Base{Type: types.BGP}, Local: "local", Remote: "remote", RemoteAS: "remote_as", Table: "table"}
		assert.NotEmpty(t, log.ID())
	})
	t.Run("isis attrs", func(t *testing.T) {
		t.Parallel()
		log := &types.ISISLog{Base: types.Base{Type: types.ISIS, Extra: nil, Original: "original"},
			Local:     "local",
			Remote:    "remote",
			Interface: "interface",
			Timestamp: time.Now(),
			State:     types.UP,
		}
		attrs := log.Attrs()
		assert.Equal(t, log.Local, attrs["local"])
		assert.Equal(t, log.Remote, attrs["remote"])
		assert.Equal(t, log.Interface, attrs["interface"])
		assert.Equal(t, log.Timestamp, attrs["timestamp"])
		assert.Equal(t, log.State, attrs["state"])
		assert.Equal(t, log.Type, attrs["type"])
		assert.Equal(t, log.Extra, attrs["extra"])
	})
	t.Run("bgp attrs", func(t *testing.T) {
		t.Parallel()
		log := &types.BGPLog{Base: types.Base{Type: types.ISIS, Extra: nil, Original: "original"},
			Local:     "local",
			Remote:    "remote",
			RemoteAS:  "remote_as",
			Timestamp: time.Now(),
			State:     types.UP,
			Table:     "table",
		}
		attrs := log.Attrs()
		assert.Equal(t, log.Local, attrs["local"])
		assert.Equal(t, log.Remote, attrs["remote"])
		assert.Equal(t, log.RemoteAS, attrs["remote_as"])
		assert.Equal(t, log.Timestamp, attrs["timestamp"])
		assert.Equal(t, log.State, attrs["state"])
		assert.Equal(t, log.Type, attrs["type"])
		assert.Equal(t, log.Extra, attrs["extra"])
		assert.Equal(t, log.Table, attrs["table"])
	})
}
