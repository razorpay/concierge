package ingress_driver

import (
	"concierge/models"
	"concierge/pkg"
)

type ShowAllowedIngressRequest struct {
	User      *models.Users
	Namespace string
}

type ShowAllowedIngressResponse struct {
	Ingresses []pkg.IngressList
}

type IngressDriver interface {
	ShowAllowedIngress(ShowAllowedIngressRequest) (ShowAllowedIngressResponse, error)
	Whitelist()
	Terminate()
	ShowIngressDetails()
	GetName() string
}
type KubeIngressDriver struct {
}

type LookerIngressDriver struct {
}

func GetIngressDriverForNamespace(ns string) IngressDriver {
	switch ns {
	case "looker":
		return &LookerIngressDriver{}
	default:
		return &KubeIngressDriver{}
	}
}

func GetIngressDrivers() []IngressDriver {
	return []IngressDriver{&LookerIngressDriver{}, &KubeIngressDriver{}}
}
