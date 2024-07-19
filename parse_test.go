package parselog_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stellaraf/go-parselog"
	"github.com/stellaraf/go-parselog/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Parse(t *testing.T) {
	t.Run("junos base", func(t *testing.T) {
		t.Parallel()
		now := time.Now()
		req := &types.Request{
			Messages:  []string{"BGP peer 2604:c0c0:3000::13e2 (Internal AS 14525) changed state from OpenConfirm to Established (event RecvKeepAlive) (instance master)"},
			Timestamp: now,
			Platform:  "junos",
			Source:    "er01.gvl01.as14525.net",
		}
		result, err := parselog.Parse(req)
		require.NoError(t, err)
		log, ok := result[0].(*types.BGPLog)
		require.True(t, ok)
		assert.Equal(t, req.Timestamp, log.Timestamp)
		assert.Equal(t, "2604:c0c0:3000::13e2", log.Remote)
		assert.Equal(t, "14525", log.RemoteAS)
		assert.True(t, log.Up())
		assert.Equal(t, "master", log.Table)
		assert.True(t, result[0].Is(parselog.BGPLogType))
	})
	t.Run("arista base", func(t *testing.T) {
		t.Parallel()
		now := time.Now()
		req := &types.Request{
			Messages:  []string{"L2 Neighbor State Change for SystemID 1004.2550.1100 on Et5 to UP"},
			Timestamp: now,
			Platform:  "arista_eos",
			Source:    "leaf0401",
		}
		result, err := parselog.Parse(req)
		require.NoError(t, err)
		log, ok := result[0].(*types.ISISLog)
		require.True(t, ok)
		assert.Equal(t, req.Timestamp, log.Timestamp)
		assert.Equal(t, "1004.2550.1100", log.Remote)
		assert.Equal(t, "Et5", log.Interface)
		assert.True(t, log.Up())
		assert.True(t, result[0].Is(parselog.ISISLogType))
	})
	t.Run("junos json", func(t *testing.T) {
		t.Parallel()
		raw := []byte(`{"message":"IS-IS lost L2 adjacency to er02.hnl01.as14525.net on ae0.3613, reason: Aged out ","platform":"junos","source":"er01.gvl01.as14525.net","timestamp":"2024-07-13 21:57:59","extra":{"key":"value"}}`)
		var req *types.Request
		err := json.Unmarshal(raw, &req)
		require.NoError(t, err)
		result, err := parselog.Parse(req)
		require.NoError(t, err)
		log, ok := result[0].(*types.ISISLog)
		require.True(t, ok)
		assert.Equal(t, req.Timestamp, log.Timestamp)
		assert.Equal(t, "er02.hnl01.as14525.net", log.Remote)
		assert.Equal(t, "ae0.3613", log.Interface)
		assert.True(t, log.Down())
		assert.True(t, result[0].Is(parselog.ISISLogType))
	})
	t.Run("arista json", func(t *testing.T) {
		t.Parallel()
		raw := []byte(`{"message":"peer 10.0.0.1 (VRF default AS 65000) old state OpenConfirm event Established new state Established","platform":"arista_eos","source":"leaf0401","timestamp":"2024-07-13 21:57:59","extra":{"key":"value"}}`)
		var req *types.Request
		err := json.Unmarshal(raw, &req)
		require.NoError(t, err)
		result, err := parselog.Parse(req)
		require.NoError(t, err)
		log, ok := result[0].(*types.BGPLog)
		require.True(t, ok)
		assert.Equal(t, req.Timestamp, log.Timestamp)
		assert.Equal(t, "10.0.0.1", log.Remote)
		assert.Equal(t, "65000", log.RemoteAS)
		assert.True(t, log.Up())
		assert.True(t, result[0].Is(parselog.BGPLogType))
	})
	t.Run("junos json multiple logs", func(t *testing.T) {
		t.Parallel()
		raw := []byte(`{"message":"IS-IS new L2 adjacency to er02.phx01.as14525.net on xe-0/0/13.3607__  IS-IS new L2 adjacency to er01.phx01.as14525.net on xe-0/0/13.3607__ ","platform":"junos","source":"er01.gvl01.as14525.net","timestamp":"2024-07-13 21:57:59","extra":{"key":"value"}}`)
		var req *types.Request
		err := json.Unmarshal(raw, &req)
		require.NoError(t, err)
		result, err := parselog.Parse(req)
		require.NoError(t, err)
		for _, _log := range result {
			log, ok := _log.(*types.ISISLog)
			require.True(t, ok)
			assert.Equal(t, req.Timestamp, log.Timestamp)
			assert.NotEmpty(t, log.Local)
			assert.NotEmpty(t, log.Interface)
			assert.True(t, log.Up())
			assert.True(t, log.Is(parselog.ISISLogType))
		}
	})
	t.Run("junos json multiple logs", func(t *testing.T) {
		t.Parallel()
		raw := []byte(`{"message":"L2 Neighbor State Change for SystemID 1004.2550.1100 on Et5 to UP__  L2 Neighbor State Change for SystemID 1004.2550.1100 on Et5 to UP__ ","platform":"arista_eos","source":"leaf0401","timestamp":"2024-07-13 21:57:59","extra":{"key":"value"}}`)
		var req *types.Request
		err := json.Unmarshal(raw, &req)
		require.NoError(t, err)
		result, err := parselog.Parse(req)
		require.NoError(t, err)
		for _, _log := range result {
			log, ok := _log.(*types.ISISLog)
			require.True(t, ok)
			assert.Equal(t, req.Timestamp, log.Timestamp)
			assert.Equal(t, "leaf0401", log.Local)
			assert.Equal(t, "1004.2550.1100", log.Remote)
			assert.Equal(t, "Et5", log.Interface)
			assert.True(t, log.Up())
			assert.True(t, log.Is(parselog.ISISLogType))
		}
	})
	t.Run("no matching platform", func(t *testing.T) {
		t.Parallel()
		req := &types.Request{Platform: "no-match"}
		_, err := parselog.Parse(req)
		assert.ErrorIs(t, err, types.ErrNoMatchingPlatform)
	})
}
