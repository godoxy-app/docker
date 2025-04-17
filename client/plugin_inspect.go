package client // import "github.com/docker/docker/client"

import (
	"bytes"
	"context"
	"io"

	"github.com/bytedance/sonic"
	"github.com/docker/docker/api/types"
)

// PluginInspectWithRaw inspects an existing plugin
func (cli *Client) PluginInspectWithRaw(ctx context.Context, name string) (*types.Plugin, []byte, error) {
	name, err := trimID("plugin", name)
	if err != nil {
		return nil, nil, err
	}
	resp, err := cli.get(ctx, "/plugins/"+name+"/json", nil, nil)
	defer ensureReaderClosed(resp)
	if err != nil {
		return nil, nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	var p types.Plugin
	rdr := bytes.NewReader(body)
	err = sonic.ConfigDefault.NewDecoder(rdr).Decode(&p)
	return &p, body, err
}
