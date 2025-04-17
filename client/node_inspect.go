package client // import "github.com/docker/docker/client"

import (
	"bytes"
	"context"
	"io"

	"github.com/bytedance/sonic"
	"github.com/docker/docker/api/types/swarm"
)

// NodeInspectWithRaw returns the node information.
func (cli *Client) NodeInspectWithRaw(ctx context.Context, nodeID string) (swarm.Node, []byte, error) {
	nodeID, err := trimID("node", nodeID)
	if err != nil {
		return swarm.Node{}, nil, err
	}
	resp, err := cli.get(ctx, "/nodes/"+nodeID, nil, nil)
	defer ensureReaderClosed(resp)
	if err != nil {
		return swarm.Node{}, nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return swarm.Node{}, nil, err
	}

	var response swarm.Node
	rdr := bytes.NewReader(body)
	err = sonic.ConfigDefault.NewDecoder(rdr).Decode(&response)
	return response, body, err
}
