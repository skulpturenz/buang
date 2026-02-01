package ports

type Ports interface {
	DockerCompose() DockerCompose
	Git() Git
}

var _ Ports = (*ports)(nil)

type ports struct {
}

func New() Ports {
	return ports{}
}

func (p ports) DockerCompose() DockerCompose {
	return &DockerComposeImpl{}
}

func (p ports) Git() Git {
	return &GitImpl{}
}
