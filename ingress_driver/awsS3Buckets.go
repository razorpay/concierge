package ingress_driver

import (
	"concierge/config"
	"concierge/pkg"
	"errors"
	"strings"
)

type AWSS3BucketIngressDriver struct {
	Buckets string
}

func getAWSS3BucketIngressDriver(S3Bucket string) IngressDriver {
	return &AWSS3BucketIngressDriver{S3Bucket}
}

func (k *AWSS3BucketIngressDriver) ShowAllowedIngress() (ShowAllowedIngressResponse, error) {
	// var err error
	client := pkg.GetAWSSession()
	resp := ShowAllowedIngressResponse{
		Buckets: []pkg.S3BucketList{},
	}

	result, err := client.GetS3Buckets()
	if err != nil {
		return resp, err
	}

	for _, S3Bucket := range result.Buckets {
		if client.IsConciergeTagTrue(*S3Bucket.Name) == nil {
			resp.Buckets = append(resp.Buckets, pkg.S3BucketList{
				Name:           *S3Bucket.Name,
				WhitelistedIps: nil,
			})
		}
	}

	return resp, nil
}

func (k *AWSS3BucketIngressDriver) EnableLease(req EnableLeaseRequest) (EnableLeaseResponse, error) {
	var err error
	resp := EnableLeaseResponse{}
	if !k.isEnabled() {
		return resp, errors.New("driver isn't enabled")
	}

	client := pkg.GetAWSSession()
	err = client.IsConciergeTagTrue(req.Name)
	if err != nil {
		return resp, err
	}

	ips := req.GinContext.Request.Header["X-Forwarded-For"][0]
	ip := strings.Split(ips, ",")[0]
	ip = ip + "/32"

	resp.UpdateStatusFlag, err = client.WhitelistIPInS3Bucket(req.Name, ip)
	if err != nil {
		return resp, err
	}

	resp.LeaseIdentifier = ip
	resp.LeaseType = k.GetLeaseType()

	return resp, nil
}

func (k *AWSS3BucketIngressDriver) DisableLease(req DisableLeaseRequest) (DisableLeaseResponse, error) {
	resp := DisableLeaseResponse{}
	if !k.isEnabled() {
		return resp, errors.New("driver isn't enabled")
	}
	client := pkg.GetAWSSession()
	err := client.IsConciergeTagTrue(req.Name)
	if err != nil {
		return resp, err
	}

	err = client.RevokeIPFromS3Bucket(req.Name, req.LeaseIdentifier)
	if err != nil {
		return resp, err
	}
	resp.UpdateStatusFlag = true
	return resp, nil
}

func (k *AWSS3BucketIngressDriver) isEnabled() bool {
	return config.AWSS3Buckets != ""
}

func (k *AWSS3BucketIngressDriver) GetLeaseType() string {
	return "awss3"
}

func (k *AWSS3BucketIngressDriver) ShowIngressDetails(req ShowIngressDetailsRequest) (ShowIngressDetailsResponse, error) {

	resp := ShowIngressDetailsResponse{}

	if !k.isEnabled() {
		return resp, errors.New("driver isn't enabled")
	}

	client := pkg.GetAWSSession()
	err := client.IsConciergeTagTrue(req.Name)
	if err != nil {
		return resp, err
	}

	resp.S3Bucket = pkg.S3BucketList{
		Name: k.Buckets,
	}
	return resp, nil
}

func (k *AWSS3BucketIngressDriver) GetName() string {
	return "aws"
}
