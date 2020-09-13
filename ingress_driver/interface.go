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
	EnableUser(EnableUserRequest) (EnableUserResponse, error)
	DisableUser(DisableUserRequest) (DisableUserResponse, error)
	ShowIngressDetails(req ShowIngressDetailsRequest) (ShowIngressDetailsResponse, error)
	GetName() string
	isEnabled() bool
}

type CommonRequest struct {
	Namespace string
	Name      string
}

type ShowAllowedIngressRequest struct {
	User *models.Users
	CommonRequest
}

type ShowAllowedIngressResponse struct {
	Ingresses []pkg.IngressList
}

type EnableUserRequest struct {
	CommonRequest
	GinContext *gin.Context
	User       *models.Users
}

type EnableUserResponse struct {
	UpdateStatusFlag bool
	Ingress          pkg.IngressList
	Identifier       string
}

type DisableUserRequest struct {
	CommonRequest
	LeaseIdentifier string
}

type DisableUserResponse struct {
	UpdateStatusFlag bool
	Ingress          pkg.IngressList
}

type ShowIngressDetailsRequest struct {
	CommonRequest
}

type ShowIngressDetailsResponse struct {
	Ingress pkg.IngressList
}

func GetIngressDriverForNamespace(ns string) IngressDriver {
	switch ns {
	case Looker:
		return getLookerIngressDriver()
	default:
		return &KubeIngressDriver{}
	}
}

func GetEnabledIngressDrivers() []IngressDriver {
	var enabledDriverss []IngressDriver
	for _, driver := range getAllDrivers() {
		if driver.isEnabled() {
			enabledDriverss = append(enabledDriverss, driver)
		}
	}
	return enabledDriverss
}

func getAllDrivers() []IngressDriver {
	drivers := []IngressDriver{getLookerIngressDriver(), getKubernetesIngressDriver()}
	return drivers
}
