package eventhandlers

type InstanceHostMapData struct {
	Instance InstanceData `json:"instance"`
}

type InstanceData struct {
	Name string `json:"name"`
	Data struct {
		DockerContainer struct {
			HostConfig struct {
				NetworkMode string `json:"NetworkMode"`
			} `json:"HostConfig"`
		} `json:"dockerContainer"`
	}
}
