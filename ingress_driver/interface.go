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
	ShowAllowedIngress() (ShowAllowedIngressResponse, error)
	EnableLease(EnableLeaseRequest) (EnableLeaseResponse, error)
	DisableLease(DisableLeaseRequest) (DisableLeaseResponse, error)
	ShowIngressDetails(req ShowIngressDetailsRequest) (ShowIngressDetailsResponse, error)
	GetName() string
	GetLeaseType() string
	isEnabled() bool
}

type ShowAllowedIngressResponse struct {
	Ingresses []pkg.IngressList
}

type EnableLeaseRequest struct {
	Name       string
	GinContext *gin.Context
	User       *models.Users
}

type EnableLeaseResponse struct {
	UpdateStatusFlag bool
	Ingress          pkg.IngressList
	LeaseIdentifier  string
	LeaseType        string
}

type DisableLeaseRequest struct {
	Name            string
	LeaseIdentifier string
}

type DisableLeaseResponse struct {
	UpdateStatusFlag bool
	Ingress          pkg.IngressList
}

type ShowIngressDetailsRequest struct {
	Name string
}

type ShowIngressDetailsResponse struct {
	Ingress pkg.IngressList
}

func GetIngressDriverForNamespace(ns string) IngressDriver {
	switch ns {
	case Looker:
		return getLookerIngressDriver()
	default:
		return getKubernetesIngressDriver(ns)
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

func GetLeaseTypes() []string {
	var leaseTypes []string

	for _, driver := range getAllDrivers() {
		leaseTypes = append(leaseTypes, driver.GetLeaseType())
	}
	return leaseTypes
}

func getAllDrivers() []IngressDriver {
	drivers := []IngressDriver{getLookerIngressDriver(), getKubernetesIngressDriver("")}
	return drivers
}
