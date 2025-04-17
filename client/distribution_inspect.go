package client // import "github.com/docker/docker/client"

import (
	"context"
	"net/http"
	"net/url"

	"github.com/bytedance/sonic"
	"github.com/docker/docker/api/types/registry"
)

// DistributionInspect returns the image digest with the full manifest.
func (cli *Client) DistributionInspect(ctx context.Context, imageRef, encodedRegistryAuth string) (registry.DistributionInspect, error) {
	if imageRef == "" {
		return registry.DistributionInspect{}, objectNotFoundError{object: "distribution", id: imageRef}
	}

	if err := cli.NewVersionError(ctx, "1.30", "distribution inspect"); err != nil {
		return registry.DistributionInspect{}, err
	}

	var headers http.Header
	if encodedRegistryAuth != "" {
		headers = http.Header{
			registry.AuthHeader: {encodedRegistryAuth},
		}
	}

	// Contact the registry to retrieve digest and platform information
	resp, err := cli.get(ctx, "/distribution/"+imageRef+"/json", url.Values{}, headers)
	defer ensureReaderClosed(resp)
	if err != nil {
		return registry.DistributionInspect{}, err
	}

	var distributionInspect registry.DistributionInspect
	err = sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&distributionInspect)
	return distributionInspect, err
}
