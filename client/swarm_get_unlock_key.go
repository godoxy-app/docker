package client // import "github.com/docker/docker/client"

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/docker/docker/api/types"
)

// SwarmGetUnlockKey retrieves the swarm's unlock key.
func (cli *Client) SwarmGetUnlockKey(ctx context.Context) (types.SwarmUnlockKeyResponse, error) {
	resp, err := cli.get(ctx, "/swarm/unlockkey", nil, nil)
	defer ensureReaderClosed(resp)
	if err != nil {
		return types.SwarmUnlockKeyResponse{}, err
	}

	var response types.SwarmUnlockKeyResponse
	err = sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&response)
	return response, err
}
