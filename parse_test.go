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
	t.Run("base", func(t *testing.T) {
		t.Parallel()
		now := time.Now()
		req := &types.Request{
			Message:   "BGP peer 2604:c0c0:3000::13e2 (Internal AS 14525) changed state from OpenConfirm to Established (event RecvKeepAlive) (instance master)",
			Timestamp: types.Timestamp{Time: now},
			Platform:  "junos",
			Source:    "er01.gvl01.as14525.net",
		}
		result, err := parselog.Parse(req)
		require.NoError(t, err)
		log, ok := result.(*types.BGPLog)
		require.True(t, ok)
		assert.Equal(t, req.Timestamp.Time, log.Timestamp)
		assert.Equal(t, "2604:c0c0:3000::13e2", log.Remote)
		assert.Equal(t, "14525", log.RemoteAS)
		assert.True(t, log.Up())
		assert.Equal(t, "master", log.Table)
	})
	t.Run("json", func(t *testing.T) {
		t.Parallel()
		raw := []byte(`{"message":"IS-IS lost L2 adjacency to er02.hnl01.as14525.net on ae0.3613, reason: Aged out ","platform":"junos","source":"er01.gvl01.as14525.net","timestamp":"2024-07-13 21:57:59","extra":{"key":"value"}}`)
		var req *types.Request
		err := json.Unmarshal(raw, &req)
		require.NoError(t, err)
	})

	t.Run("no matching platform", func(t *testing.T) {
		t.Parallel()
		req := &types.Request{Platform: "no-match"}
		_, err := parselog.Parse(req)
		assert.ErrorIs(t, err, types.ErrNoMatchingPlatform)
	})
}
