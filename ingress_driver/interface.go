package ingress_driver

import (
	"concierge/config"
	"concierge/models"
	"concierge/pkg"

	"github.com/gin-gonic/gin"
)

const (
	Looker           = "looker"
	AWSSecurityGroup = "awssg"
	Ingress          = "ingress"
	AWSS3Bucket      = "awss3"

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
	Ingresses      []pkg.IngressList
	SecurityGroups []pkg.SecurityGroupList
	Buckets        []pkg.S3BucketList
	Looker         []pkg.IngressList
}

type EnableLeaseRequest struct {
	Name          string
	SecurityGroup config.SecurityGroupIngress
	GinContext    *gin.Context
	User          *models.Users
}

type EnableLeaseResponse struct {
	UpdateStatusFlag bool
	Ingress          pkg.IngressList
	Looker           pkg.IngressList
	SecurityGroup    pkg.SecurityGroupList
	LeaseIdentifier  string
	LeaseType        string
}

type DisableLeaseRequest struct {
	Name            string
	SecurityGroup   config.SecurityGroupIngress
	LeaseIdentifier string
}

type DisableLeaseResponse struct {
	UpdateStatusFlag bool
	Ingress          pkg.IngressList
	Looker           pkg.IngressList
	SecurityGroup    pkg.SecurityGroupList
}

type ShowIngressDetailsRequest struct {
	Name string
}

type ShowIngressDetailsResponse struct {
	Ingress       pkg.IngressList
	Looker        pkg.IngressList
	SecurityGroup pkg.SecurityGroupList
	S3Bucket      pkg.S3BucketList
}

func GetIngressDriverForNamespace(driver string, ns string) IngressDriver {
	switch driver {
	case Looker:
		return getLookerIngressDriver()
	case AWSSecurityGroup:
		return getAWSSecurityGroupIngressDriver(ns)
	case AWSS3Bucket:
		return getAWSS3BucketIngressDriver(ns)
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
	drivers := []IngressDriver{getLookerIngressDriver(), getKubernetesIngressDriver(""), getAWSSecurityGroupIngressDriver(""), getAWSS3BucketIngressDriver("")}
	return drivers
}
