package container

type (
	SummaryTrimmed struct {
		ID         string `json:"Id"`
		Names      []string
		Image      string
		Ports      []Port
		Labels     map[string]string
		State      ContainerState
		Status     string
		HostConfig struct {
			NetworkMode string `json:",omitempty"`
		}
		NetworkSettings *NetworkSettingsSummaryTrimmed
		Mounts          []MountPointTrimmed
	}

	// NetworkSettingsSummary provides a summary of container's networks
	// in /containers/json
	NetworkSettingsSummaryTrimmed struct {
		Networks map[string]*struct {
			IPAddress string
		}
	}

	MountPointTrimmed struct {
		Destination string
	}
)
