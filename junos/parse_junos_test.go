package junos_test

import (
	"testing"
	"time"

	"github.com/stellaraf/go-parselog/junos"
	"github.com/stellaraf/go-parselog/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseISIS(t *testing.T) {
	t.Run("up", func(t *testing.T) {
		t.Parallel()
		now := time.Now()
		msg := "IS-IS new L2 adjacency to er02.hnl01.as14525.net on ae0.3613"
		result, err := junos.ParseISIS(msg, "er01.gvl01.as14525.net", now, nil)
		require.NoError(t, err)
		attrs := result.Attrs()
		assert.Equal(t, types.UP, attrs["state"])
		assert.Equal(t, "er02.hnl01.as14525.net", attrs["remote"])
		assert.Equal(t, "ae0.3613", attrs["interface"])
		assert.Equal(t, msg, attrs["original"])
		assert.Empty(t, attrs["reason"])
		assert.False(t, result.Down())
		assert.True(t, result.Up())
	})
	t.Run("down", func(t *testing.T) {
		t.Parallel()
		now := time.Now()
		msg := "IS-IS lost L2 adjacency to er02.hnl01.as14525.net on ae0.3613, reason: Aged out"
		result, err := junos.ParseISIS(msg, "er01.gvl01.as14525.net", now, nil)
		require.NoError(t, err)
		attrs := result.Attrs()
		assert.Equal(t, types.DOWN, attrs["state"])
		assert.Equal(t, "er02.hnl01.as14525.net", attrs["remote"])
		assert.Equal(t, "ae0.3613", attrs["interface"])
		assert.Equal(t, "Aged out", attrs["reason"])
		assert.Equal(t, msg, attrs["original"])
		assert.False(t, result.Up())
		assert.True(t, result.Down())
	})
	t.Run("missing fields", func(t *testing.T) {
		t.Parallel()
		_, err := junos.ParseISIS("IS-IS lost L2 adjacency to er02.hnl01.as14525.net", "", time.Now(), nil)
		assert.ErrorIs(t, err, types.ErrIncompleteMatch)
	})
}

func Test_ParseBGP(t *testing.T) {
	t.Run("up", func(t *testing.T) {
		t.Parallel()
		msg := "BGP peer 2604:c0c0:3000::13e2 (Internal AS 14525) changed state from OpenConfirm to Established (event RecvKeepAlive) (instance master)"
		result, err := junos.ParseBGP(msg, "er01.gvl01.as14525.net", time.Now(), nil)
		require.NoError(t, err)
		attrs := result.Attrs()
		assert.Equal(t, types.UP, attrs["state"])
		assert.Equal(t, "2604:c0c0:3000::13e2", attrs["remote"])
		assert.Equal(t, "14525", attrs["remote_as"])
		assert.Equal(t, "master", attrs["table"])
		assert.Equal(t, msg, attrs["original"])
		assert.False(t, result.Down())
		assert.True(t, result.Up())
	})
	t.Run("down", func(t *testing.T) {
		t.Parallel()
		msg := "BGP peer 2604:c0c0:3000::13e2 (Internal AS 14525) changed state from Established to Idle (event RecvNotify) (instance master)"
		result, err := junos.ParseBGP(msg, "er01.gvl01.as14525.net", time.Now(), nil)
		require.NoError(t, err)
		attrs := result.Attrs()
		assert.Equal(t, types.DOWN, attrs["state"])
		assert.Equal(t, "2604:c0c0:3000::13e2", attrs["remote"])
		assert.Equal(t, "14525", attrs["remote_as"])
		assert.Equal(t, "master", attrs["table"])
		assert.Equal(t, msg, attrs["original"])
		assert.False(t, result.Up())
		assert.True(t, result.Down())
	})
	t.Run("missing fields", func(t *testing.T) {
		t.Parallel()
		_, err := junos.ParseBGP("BGP peer 2604:c0c0:3000::13e2 (Internal AS 14525)", "", time.Now(), nil)
		assert.ErrorIs(t, err, types.ErrIncompleteMatch)
	})
}

func Test_Parse(t *testing.T) {
	t.Run("isis", func(t *testing.T) {
		t.Parallel()
		req := &types.Request{Messages: []string{"IS-IS new L2 adjacency to er02.hnl01.as14525.net on ae0.3613"}}
		result, err := junos.Parse(req)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
	t.Run("bgp", func(t *testing.T) {
		t.Parallel()
		req := &types.Request{Messages: []string{"BGP peer 2604:c0c0:3000::13e2 (Internal AS 14525) changed state from OpenConfirm to Established (event RecvKeepAlive) (instance master)"}}
		result, err := junos.Parse(req)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
	t.Run("no match", func(t *testing.T) {
		t.Parallel()
		req := &types.Request{Messages: []string{"this has no match"}}
		result, err := junos.Parse(req)
		assert.ErrorIs(t, err, types.ErrNoMatchingParser)
		assert.Nil(t, result)
	})
	t.Run("with extra", func(t *testing.T) {
		t.Parallel()
		req := &types.Request{Messages: []string{"BGP peer 2604:c0c0:3000::13e2 (Internal AS 14525) changed state from OpenConfirm to Established (event RecvKeepAlive) (instance master)"}, Extra: map[string]any{"key": "value"}}
		result, err := junos.Parse(req)
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
				"BGP peer 2604:c0c0:3000::13e2 (Internal AS 14525) changed state from OpenConfirm to Established (event RecvKeepAlive) (instance master)",
				"IS-IS incomplete",
			},
		}
		result, err := junos.Parse(req)
		require.Error(t, err)
		assert.ErrorIs(t, err, types.ErrIncompleteMatch)
		assert.Nil(t, result)
	})
}
