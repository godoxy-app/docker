package client // import "github.com/docker/docker/client"

import (
	"bytes"
	"context"
	"io"

	"github.com/bytedance/sonic"
	"github.com/docker/docker/api/types/swarm"
)

// SecretInspectWithRaw returns the secret information with raw data
func (cli *Client) SecretInspectWithRaw(ctx context.Context, id string) (swarm.Secret, []byte, error) {
	id, err := trimID("secret", id)
	if err != nil {
		return swarm.Secret{}, nil, err
	}
	if err := cli.NewVersionError(ctx, "1.25", "secret inspect"); err != nil {
		return swarm.Secret{}, nil, err
	}
	resp, err := cli.get(ctx, "/secrets/"+id, nil, nil)
	defer ensureReaderClosed(resp)
	if err != nil {
		return swarm.Secret{}, nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return swarm.Secret{}, nil, err
	}

	var secret swarm.Secret
	rdr := bytes.NewReader(body)
	err = sonic.ConfigDefault.NewDecoder(rdr).Decode(&secret)

	return secret, body, err
}
