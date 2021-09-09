package ingress_driver

import (
	"concierge/config"
	"concierge/pkg"
	"errors"

	log "github.com/sirupsen/logrus"
)

type LookerIngressDriver struct {
	ingress pkg.IngressList
}

var lookerIngressDriver *LookerIngressDriver

func getLookerIngressDriver() IngressDriver {
	host := config.LookerConfig.BaseUrl

	if lookerIngressDriver == nil {
		lookerIngressDriver = &LookerIngressDriver{ingress: struct {
			Name           string
			Namespace      string
			Context        string
			Host           string
			Class          string
			WhitelistedIps []string
		}{Name: Looker, Namespace: Looker, Context: DefaultContext, Host: host, Class: DefaultClass, WhitelistedIps: nil}}

	}
	return lookerIngressDriver
}

func (k *LookerIngressDriver) ShowAllowedIngress() (ShowAllowedIngressResponse, error) {
	resp := ShowAllowedIngressResponse{
		Looker: []pkg.IngressList{k.ingress},
	}

	return resp, nil
}

func (k *LookerIngressDriver) EnableLease(req EnableLeaseRequest) (EnableLeaseResponse, error) {
	log.Infof("Received Looker EnableLeaseRequest %v", req.User.Name)
	resp := EnableLeaseResponse{
		Looker:          k.ingress,
		LeaseIdentifier: req.User.Email,
		LeaseType:       k.GetLeaseType(),
	}

	client := pkg.GetLookerClient()

	users, searchErr := client.SearchUser(pkg.LookerSearchUserRequest{Email: req.User.Email})

	if searchErr != nil {
		return resp, searchErr
	}

	if len(users) == 0 {
		userCreationErr := client.CreateLookerUser(req.User.Email)

		userAttributeErr := client.UpdateLookerUserAttribute(req.User.Email)
		if userAttributeErr != nil {
			return resp, userAttributeErr
		}

		resp.UpdateStatusFlag = true

		return resp, userCreationErr
	}

	if len(users) > 1 {
		return resp, errors.New("you have multiple looker accounts. Please contact looker admins")
	}

	user := users[0]

	if !user.IsDisabled {
		userAttributeErr := client.UpdateLookerUserAttribute(req.User.Email)
		if userAttributeErr != nil {
			return resp, userAttributeErr
		}
		return resp, errors.New("looker user already enabled for user. Please visit looker")
	}

	patchedUser, patchErr := client.PatchUser(user.Id, pkg.LookerPatchUserRequest{IsDisabled: false})

	if patchErr != nil {
		return resp, patchErr
	}

	if patchedUser.IsDisabled {
		return resp, errors.New("failed to enable user. please contact admin")
	}

	resp.UpdateStatusFlag = true

	userAttributeErr := client.UpdateLookerUserAttribute(req.User.Email)
	if userAttributeErr != nil {
		return resp, userAttributeErr
	}

	return resp, nil
}

func (k *LookerIngressDriver) DisableLease(req DisableLeaseRequest) (DisableLeaseResponse, error) {
	log.Infof("Received Looker DisableLeaseRequest %v", req.LeaseIdentifier)
	response := DisableLeaseResponse{}

	client := pkg.GetLookerClient()

	users, searchErr := client.SearchUser(pkg.LookerSearchUserRequest{Email: req.LeaseIdentifier})

	if searchErr != nil {
		return DisableLeaseResponse{Ingress: k.ingress}, searchErr
	}

	for _, user := range users {
		if user.IsDisabled {
			continue
		}

		patchedUser, patchErr := client.PatchUser(user.Id, pkg.LookerPatchUserRequest{IsDisabled: true})

		if patchErr != nil || !patchedUser.IsDisabled {
			return response, errors.New("failed to delete lease. contact looker admin")
		}
	}

	response.UpdateStatusFlag = true

	return response, nil
}

func (k *LookerIngressDriver) ShowIngressDetails(ShowIngressDetailsRequest) (ShowIngressDetailsResponse, error) {
	return ShowIngressDetailsResponse{Looker: k.ingress}, nil
}

func (k *LookerIngressDriver) GetName() string {
	return Looker
}

func (k *LookerIngressDriver) GetLeaseType() string {
	return Looker
}

func (k *LookerIngressDriver) isEnabled() bool {
	return config.LookerConfig.IsEnabled
}
