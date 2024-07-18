package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stellaraf/go-parselog/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Request(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		var req *types.Request
		err := json.Unmarshal([]byte(`{"message":"IS-IS lost L2 adjacency to er02.hnl01.as14525.net on ae0.3613, reason: Aged out ","platform":"junos","source":"er01.gvl01.as14525.net","timestamp":"2024-07-13 21:57:59"}`), &req)
		require.NoError(t, err)
	})
	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()
		var req *types.Request
		err := json.Unmarshal([]byte(`["not a request"]`), &req)
		require.Error(t, err)
	})
	t.Run("missing fields", func(t *testing.T) {
		t.Parallel()
		var req *types.Request
		err := json.Unmarshal([]byte(`{"platform":"","source":"","timestamp":"2024-07-13 21:57:59","extra":{}}`), &req)
		require.Error(t, err)
		err = json.Unmarshal([]byte(`{"message":"","source":"","timestamp":"2024-07-13 21:57:59","extra":{}}`), &req)
		require.Error(t, err)
		err = json.Unmarshal([]byte(`{"message":"","platform":"","timestamp":"2024-07-13 21:57:59","extra":{}}`), &req)
		require.Error(t, err)
		err = json.Unmarshal([]byte(`{"message":"","platform":"","source":"","extra":{}}`), &req)
		require.Error(t, err)
		err = json.Unmarshal([]byte(`{"message":"","platform":"","source":"","timestamp":"2024-07-13 21:57:59","extra":""}`), &req)
		require.Error(t, err)
		err = json.Unmarshal([]byte(`{"message":"","platform":"","source":"","timestamp":"2024-07-13 21:57:59","extra":{}}`), &req)
		require.NoError(t, err)
		err = json.Unmarshal([]byte(`{"message":"","platform":"","source":"","timestamp":"2024-07-13 21:57:59"}`), &req)
		require.NoError(t, err)
		assert.IsType(t, map[string]any{}, req.Extra)
	})
	t.Run("invalid fields", func(t *testing.T) {
		t.Parallel()
		var req *types.Request
		err := json.Unmarshal([]byte(`{"message":0,"platform":"","source":"","timestamp":"2024-07-13 21:57:59","extra":{}}`), &req)
		require.Error(t, err)
		err = json.Unmarshal([]byte(`{"message":"","platform":false,"source":"","timestamp":"2024-07-13 21:57:59","extra":{}}`), &req)
		require.Error(t, err)
		err = json.Unmarshal([]byte(`{"message":"","platform":"","source":[],"timestamp":"2024-07-13 21:57:59","extra":{}}`), &req)
		require.Error(t, err)
		err = json.Unmarshal([]byte(`{"message":"","platform":"","source":"","timestamp":{},"extra":{}}`), &req)
		require.Error(t, err)
		err = json.Unmarshal([]byte(`{"message":"","platform":"","source":"","timestamp":"not a time","extra":{}}`), &req)
		require.Error(t, err)
	})
}
