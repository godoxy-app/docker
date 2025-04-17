package client // import "github.com/docker/docker/client"

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/docker/docker/api/types"
)

// ServerVersion returns information of the docker client and server host.
func (cli *Client) ServerVersion(ctx context.Context) (types.Version, error) {
	resp, err := cli.get(ctx, "/version", nil, nil)
	defer ensureReaderClosed(resp)
	if err != nil {
		return types.Version{}, err
	}

	var server types.Version
	err = sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&server)
	return server, err
}
