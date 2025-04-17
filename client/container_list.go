package client // import "github.com/docker/docker/client"

import (
	"context"
	"net/url"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

// ContainerList returns the list of containers in the docker host.
func (cli *Client) ContainerList(ctx context.Context, options container.ListOptions) ([]container.SummaryTrimmed, error) {
	query := url.Values{}

	if options.All {
		query.Set("all", "1")
	}

	if options.Limit > 0 {
		query.Set("limit", strconv.Itoa(options.Limit))
	}

	if options.Since != "" {
		query.Set("since", options.Since)
	}

	if options.Before != "" {
		query.Set("before", options.Before)
	}

	if options.Size {
		query.Set("size", "1")
	}

	if options.Filters.Len() > 0 {
		//nolint:staticcheck // ignore SA1019 for old code
		filterJSON, err := filters.ToParamWithVersion(cli.version, options.Filters)
		if err != nil {
			return nil, err
		}

		query.Set("filters", filterJSON)
	}

	resp, err := cli.get(ctx, "/containers/json", query, nil)
	defer ensureReaderClosed(resp)
	if err != nil {
		return nil, err
	}

	var containers []container.SummaryTrimmed
	err = sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&containers)
	if err != nil {
		return nil, err
	}

	return containers, nil
}
