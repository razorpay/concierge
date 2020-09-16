package ingress_driver

import (
	"concierge/models"
	"concierge/pkg"
	"github.com/gin-gonic/gin"
)

const (
	Looker = "looker"

	DefaultContext = "default"
	DefaultClass   = "default"
)

type IngressDriver interface {
	ShowAllowedIngress(ShowAllowedIngressRequest) (ShowAllowedIngressResponse, error)
	EnableLease(EnableLeaseRequest) (EnableLeaseResponse, error)
	DisableLease(DisableLeaseRequest) (DisableLeaseResponse, error)
	ShowIngressDetails(req ShowIngressDetailsRequest) (ShowIngressDetailsResponse, error)
	GetName() string
	isEnabled() bool
}

type ShowAllowedIngressRequest struct {
	Namespace string
}

type ShowAllowedIngressResponse struct {
	Ingresses []pkg.IngressList
}

type EnableLeaseRequest struct {
	Namespace  string
	Name       string
	GinContext *gin.Context
	User       *models.Users
}

type EnableLeaseResponse struct {
	UpdateStatusFlag bool
	Ingress          pkg.IngressList
	Identifier       string
}

type DisableLeaseRequest struct {
	Namespace       string
	Name            string
	LeaseIdentifier string
}

type DisableLeaseResponse struct {
	UpdateStatusFlag bool
	Ingress          pkg.IngressList
}

type ShowIngressDetailsRequest struct {
	Namespace string
	Name      string
}

type ShowIngressDetailsResponse struct {
	Ingress pkg.IngressList
}

func GetIngressDriverForNamespace(ns string) IngressDriver {
	switch ns {
	case Looker:
		return getLookerIngressDriver()
	default:
		return getKubernetesIngressDriver()
	}
}

func GetEnabledIngressDrivers() []IngressDriver {
	var enabledDrivers []IngressDriver
	for _, driver := range getAllDrivers() {
		if driver.isEnabled() {
			enabledDrivers = append(enabledDrivers, driver)
		}
	}
	return enabledDrivers
}

func getAllDrivers() []IngressDriver {
	drivers := []IngressDriver{getLookerIngressDriver(), getKubernetesIngressDriver()}
	return drivers
}
