package client // import "github.com/docker/docker/client"

import (
	"context"
	"net/url"

	"github.com/bytedance/sonic"
	"github.com/docker/docker/api/types/registry"
)

// RegistryLogin authenticates the docker server with a given docker registry.
// It returns unauthorizedError when the authentication fails.
func (cli *Client) RegistryLogin(ctx context.Context, auth registry.AuthConfig) (registry.AuthenticateOKBody, error) {
	resp, err := cli.post(ctx, "/auth", url.Values{}, auth, nil)
	defer ensureReaderClosed(resp)

	if err != nil {
		return registry.AuthenticateOKBody{}, err
	}

	var response registry.AuthenticateOKBody
	err = sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&response)
	return response, err
}
