package client // import "github.com/docker/docker/client"

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/docker/docker/api/types/swarm"
)

// SwarmInspect inspects the swarm.
func (cli *Client) SwarmInspect(ctx context.Context) (swarm.Swarm, error) {
	resp, err := cli.get(ctx, "/swarm", nil, nil)
	defer ensureReaderClosed(resp)
	if err != nil {
		return swarm.Swarm{}, err
	}

	var response swarm.Swarm
	err = sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&response)
	return response, err
}
