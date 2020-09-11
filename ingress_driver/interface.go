package ingress_driver

import (
	"concierge/models"
	"concierge/pkg"
)

const (
	Looker = "looker"
)

type ShowAllowedIngressRequest struct {
	User      *models.Users
	Namespace string
}

type ShowAllowedIngressResponse struct {
	Ingresses []pkg.IngressList
}

type ShowIngressDetailsRequest struct {
	Namespace string
	Name      string
}

type ShowIngressDetailsResponse struct {
	Ingress pkg.IngressList
}

type IngressDriver interface {
	ShowAllowedIngress(ShowAllowedIngressRequest) (ShowAllowedIngressResponse, error)
	Whitelist()
	Terminate()
	ShowIngressDetails(req ShowIngressDetailsRequest) (ShowIngressDetailsResponse, error)
	GetName() string
}
type KubeIngressDriver struct {
}

type LookerIngressDriver struct {
}

func GetIngressDriverForNamespace(ns string) IngressDriver {
	switch ns {
	case Looker:
		return &LookerIngressDriver{}
	default:
		return &KubeIngressDriver{}
	}
}

func GetIngressDrivers() []IngressDriver {
	return []IngressDriver{&LookerIngressDriver{}, &KubeIngressDriver{}}
}
