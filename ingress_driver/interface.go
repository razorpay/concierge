package ingress_driver

import (
	"concierge/models"
	"concierge/pkg"
	"github.com/gin-gonic/gin"
)

const (
	Looker     = "looker"
	LookerHost = "https://looker.razorpay.com"

	DefaultContext = "default"
	DefaultClass   = "default"
)

type IngressDriver interface {
	ShowAllowedIngress(ShowAllowedIngressRequest) (ShowAllowedIngressResponse, error)
	EnableUser(EnableUserRequest) (EnableUserResponse, error)
	DisableUser(DisableUserRequest) (DisableUserResponse, error)
	ShowIngressDetails(req ShowIngressDetailsRequest) (ShowIngressDetailsResponse, error)
	GetName() string
}
type KubeIngressDriver struct {
}

type LookerIngressDriver struct {
	ingress pkg.IngressList
}

type ShowAllowedIngressRequest struct {
	User      *models.Users
	Namespace string
}

type ShowAllowedIngressResponse struct {
	Ingresses []pkg.IngressList
}

type EnableUserRequest struct {
	Namespace  string
	Name       string
	GinContext *gin.Context
	User       *models.Users
}

type EnableUserResponse struct {
	UpdateStatusFlag bool
	Ingress          pkg.IngressList
	IdentifierType   string
	Identifier       string
}

type DisableUserRequest struct {
}

type DisableUserResponse struct {
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
		return &KubeIngressDriver{}
	}
}

func getLookerIngressDriver() IngressDriver {
	return &LookerIngressDriver{ingress: struct {
		Name           string
		Namespace      string
		Context        string
		Host           string
		Class          string
		WhitelistedIps []string
	}{Name: Looker, Namespace: Looker, Context: DefaultContext, Host: LookerHost, Class: DefaultClass, WhitelistedIps: nil}}
}

func GetIngressDrivers() []IngressDriver {
	return []IngressDriver{getLookerIngressDriver(), &KubeIngressDriver{}}
}
