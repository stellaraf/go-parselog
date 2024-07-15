package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stellaraf/go-parselog/types"
	"github.com/stretchr/testify/require"
)

func Test_Request(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		raw := []byte(`{"message":"IS-IS lost L2 adjacency to er02.hnl01.as14525.net on ae0.3613, reason: Aged out ","platform":"junos","source":"er01.gvl01.as14525.net","timestamp":"2024-07-13 21:57:59"}`)
		var req *types.Request
		err := json.Unmarshal(raw, &req)
		require.NoError(t, err)
	})
	t.Run("invalid timestamp json", func(t *testing.T) {
		t.Parallel()
		raw := []byte(`{"timestamp":0}`)
		var req *types.Request
		err := json.Unmarshal(raw, &req)
		require.Error(t, err)
	})
	t.Run("invalid timestamp format", func(t *testing.T) {
		t.Parallel()
		raw := []byte(`{"timestamp":"2024-07-13T21:57:59Z"}`)
		var req *types.Request
		err := json.Unmarshal(raw, &req)
		require.Error(t, err)
	})
}
