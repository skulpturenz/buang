package ports

type Ports interface {
	DockerCompose
}

var _ Ports = (*ports)(nil)

type ports struct {
	*DockerComposeImpl
}

func New() Ports {
	return ports{
		DockerComposeImpl: &DockerComposeImpl{},
	}
}
