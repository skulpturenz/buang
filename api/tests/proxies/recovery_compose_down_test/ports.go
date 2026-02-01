package proxiesrecoverycomposedown

import "skulpture/buang/ports"

type recoveryPorts struct {
	*dockerComposeProxyImpl
}

func New() ports.Ports {
	return recoveryPorts{
		dockerComposeProxyImpl: &dockerComposeProxyImpl{},
	}
}
