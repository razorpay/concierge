package ingress_driver

import (
	"concierge/pkg"
	"strings"
)

type AWSIngressDriver struct {
	SecurityGroup string
}

func getAWSIngressDriver(SecurityGroup string) IngressDriver {
	return &AWSIngressDriver{SecurityGroup}
}

func (k *AWSIngressDriver) ShowAllowedIngress() (ShowAllowedIngressResponse, error) {
	client := pkg.GetAWSSession()
	resp := ShowAllowedIngressResponse{
		SecurityGroups: []pkg.SecurityGroupList{},
	}

	result, err := client.GetSecurityGroups()
	if err != nil {
		return resp, err
	}

	for _, securityGroup := range result.SecurityGroups {
		resp.SecurityGroups = append(resp.SecurityGroups, pkg.SecurityGroupList{
			GroupId:        *securityGroup.GroupId,
			Name:           *securityGroup.GroupName,
			Description:    *securityGroup.Description,
			VpcId:          *securityGroup.VpcId,
			WhitelistedIps: nil,
		})
	}

	return resp, nil
}

func (k *AWSIngressDriver) EnableLease(req EnableLeaseRequest) (EnableLeaseResponse, error) {
	var err error
	resp := EnableLeaseResponse{}

	ips := req.GinContext.Request.Header["X-Forwarded-For"][0]
	ip := strings.Split(ips, ",")[0]
	ip = ip + "/32"

	client := pkg.GetAWSSession()
	err = client.WhitelistIP(k.SecurityGroup, ip, req.SecurityGroup, req.User.Email)
	if err != nil {
		return resp, err
	}

	resp.UpdateStatusFlag = true
	resp.LeaseIdentifier = ip
	resp.LeaseType = k.GetLeaseType()

	return resp, nil
}

func (k *AWSIngressDriver) DisableLease(req DisableLeaseRequest) (DisableLeaseResponse, error) {
	resp := DisableLeaseResponse{}
	client := pkg.GetAWSSession()
	err := client.RevokeIP(k.SecurityGroup, req.LeaseIdentifier, req.SecurityGroup)
	if err != nil {
		return resp, err
	}
	resp.UpdateStatusFlag = true
	return resp, nil
}

func (k *AWSIngressDriver) isEnabled() bool {
	return true
}

func (k *AWSIngressDriver) GetLeaseType() string {
	return "aws"
}

func (k *AWSIngressDriver) ShowIngressDetails(req ShowIngressDetailsRequest) (ShowIngressDetailsResponse, error) {

	client := pkg.GetAWSSession()
	resp := ShowIngressDetailsResponse{}

	result, err := client.GetSecurityGroupDetails(k.SecurityGroup)
	if err != nil {
		return resp, err
	}
	resp.SecurityGroup = pkg.SecurityGroupList{
		GroupId:        *result.SecurityGroups[0].GroupId,
		Name:           *result.SecurityGroups[0].GroupName,
		Description:    *result.SecurityGroups[0].Description,
		VpcId:          *result.SecurityGroups[0].VpcId,
		WhitelistedIps: nil,
	}

	return resp, nil
}

func (k *AWSIngressDriver) GetName() string {
	return "aws"
}
