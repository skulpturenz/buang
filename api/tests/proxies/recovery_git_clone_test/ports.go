package proxiesrecoverygitclone

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
		p.dockerCompose = &ports.DockerComposeImpl{}
	}

	return p.dockerCompose
}

func (p *recoveryPorts) Git() ports.Git {
	if p.git == nil {
		p.git = &gitProxyImpl{}
	}

	return p.git
}
