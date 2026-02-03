package proxiesrecoveryimagepull

import (
	"skulpture/buang/ports"
)

type recoveryPorts struct {
	dockerCompose ports.DockerCompose
	git           ports.Git
}

func New() ports.Ports {
	return &recoveryPorts{}
}

func (p *recoveryPorts) DockerCompose() ports.DockerCompose {
	if p.dockerCompose == nil {
		p.dockerCompose = &dockerComposeProxyImpl{}
	}

	return p.dockerCompose
}

func (p *recoveryPorts) Git() ports.Git {
	if p.git == nil {
		p.git = &ports.GitImpl{}
	}

	return p.git
}
