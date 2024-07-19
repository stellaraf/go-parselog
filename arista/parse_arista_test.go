package arista_test

import (
	"testing"
	"time"

	"github.com/stellaraf/go-parselog/arista"
	"github.com/stellaraf/go-parselog/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseISIS(t *testing.T) {
	t.Run("up", func(t *testing.T) {
		t.Parallel()
		now := time.Now()
		msg := "L2 Neighbor State Change for SystemID 1004.2550.1100 on Et5 to UP"
		result, err := arista.ParseISIS(msg, "leaf0801", now, nil)
		require.NoError(t, err)
		attrs := result.Attrs()
		assert.Equal(t, types.UP, attrs["state"])
		assert.Equal(t, "1004.2550.1100", attrs["remote"])
		assert.Equal(t, "Et5", attrs["interface"])
		assert.Equal(t, msg, attrs["original"])
		assert.Empty(t, attrs["reason"])
		assert.False(t, result.Down())
		assert.True(t, result.Up())
	})
	t.Run("down", func(t *testing.T) {
		t.Parallel()
		now := time.Now()
		msg := "L2 Neighbor State Change for SystemID 1004.2550.1100 on Et5 to DOWN: interface went down or no IP address on interface"
		result, err := arista.ParseISIS(msg, "leaf0401", now, nil)
		require.NoError(t, err)
		attrs := result.Attrs()
		assert.Equal(t, types.DOWN, attrs["state"])
		assert.Equal(t, "1004.2550.1100", attrs["remote"])
		assert.Equal(t, "Et5", attrs["interface"])
		assert.Equal(t, "interface went down or no IP address on interface", attrs["reason"])
		assert.Equal(t, msg, attrs["original"])
		assert.False(t, result.Up())
		assert.True(t, result.Down())
	})
	t.Run("init", func(t *testing.T) {
		now := time.Now()
		msg := "L2 Neighbor State Change for SystemID 1004.2550.1100 on Et5 from UP to INIT"
		result, err := arista.ParseISIS(msg, "leaf0401", now, nil)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
	t.Run("missing fields", func(t *testing.T) {
		t.Parallel()
		_, err := arista.ParseISIS("L2 Neighbor State Change for SystemID 1004.2550.1100", "", time.Now(), nil)
		assert.ErrorIs(t, err, types.ErrIncompleteMatch)
	})
}

func Test_ParseBGP(t *testing.T) {
	t.Run("up", func(t *testing.T) {
		t.Parallel()
		msg := "peer 10.0.0.1 (VRF default AS 65000) old state OpenConfirm event Established new state Established"
		result, err := arista.ParseBGP(msg, "leaf0401", time.Now(), nil)
		require.NoError(t, err)
		attrs := result.Attrs()
		assert.Equal(t, types.UP, attrs["state"])
		assert.Equal(t, "10.0.0.1", attrs["remote"])
		assert.Equal(t, "65000", attrs["remote_as"])
		assert.Equal(t, "default", attrs["table"])
		assert.Equal(t, msg, attrs["original"])
		assert.False(t, result.Down())
		assert.True(t, result.Up())
	})
	t.Run("down", func(t *testing.T) {
		t.Parallel()
		msg := "peer 10.4.255.121 (VRF default AS 65004) old state Established event AdminShutdown new state Idle"
		result, err := arista.ParseBGP(msg, "exit0401", time.Now(), nil)
		require.NoError(t, err)
		attrs := result.Attrs()
		assert.Equal(t, types.DOWN, attrs["state"])
		assert.Equal(t, "10.4.255.121", attrs["remote"])
		assert.Equal(t, "65004", attrs["remote_as"])
		assert.Equal(t, "default", attrs["table"])
		assert.Equal(t, msg, attrs["original"])
		assert.False(t, result.Up())
		assert.True(t, result.Down())
	})
	t.Run("missing fields", func(t *testing.T) {
		t.Parallel()
		_, err := arista.ParseBGP("peer 10.4.255.121 (VRF default AS 65004) old state Established", "", time.Now(), nil)
		assert.ErrorIs(t, err, types.ErrIncompleteMatch)
	})
}

func Test_Parse(t *testing.T) {
	t.Run("isis", func(t *testing.T) {
		t.Parallel()
		req := &types.Request{Messages: []string{"L2 Neighbor State Change for SystemID 1004.2550.1100 on Et5 to UP"}}
		result, err := arista.Parse(req)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
	t.Run("bgp", func(t *testing.T) {
		t.Parallel()
		req := &types.Request{Messages: []string{"peer 10.0.0.1 (VRF default AS 65000) old state OpenConfirm event Established new state Established"}}
		result, err := arista.Parse(req)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
	t.Run("no match", func(t *testing.T) {
		t.Parallel()
		req := &types.Request{Messages: []string{"this has no match"}}
		result, err := arista.Parse(req)
		assert.ErrorIs(t, err, types.ErrNoMatchingParser)
		assert.Nil(t, result)
	})
	t.Run("with extra", func(t *testing.T) {
		t.Parallel()
		req := &types.Request{Messages: []string{"peer 10.0.0.1 (VRF default AS 65000) old state OpenConfirm event Established new state Established"}, Extra: map[string]any{"key": "value"}}
		result, err := arista.Parse(req)
		require.NoError(t, err)
		for _, _log := range result {
			log, ok := _log.(*types.BGPLog)
			require.True(t, ok)
			assert.Equal(t, map[string]any{"key": "value"}, log.Extra)
		}
	})
	t.Run("with invalid", func(t *testing.T) {
		t.Parallel()
		req := &types.Request{
			Messages: []string{
				"peer 10.0.0.1 (VRF default AS 65000) old state OpenConfirm event Established new state Established",
				"peer 10.0.0.1 invalid",
			},
		}
		result, err := arista.Parse(req)
		require.ErrorIs(t, err, types.ErrIncompleteMatch)
		assert.Nil(t, result)
	})
}
