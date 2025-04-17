package client // import "github.com/docker/docker/client"

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/docker/docker/api/types/volume"
)

// VolumeCreate creates a volume in the docker host.
func (cli *Client) VolumeCreate(ctx context.Context, options volume.CreateOptions) (volume.Volume, error) {
	resp, err := cli.post(ctx, "/volumes/create", nil, options, nil)
	defer ensureReaderClosed(resp)
	if err != nil {
		return volume.Volume{}, err
	}

	var vol volume.Volume
	err = sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&vol)
	return vol, err
}
