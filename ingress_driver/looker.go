package ingress_driver

import (
	"concierge/pkg"
	"errors"
	log "github.com/sirupsen/logrus"
)

func (k *LookerIngressDriver) ShowAllowedIngress(ShowAllowedIngressRequest) (ShowAllowedIngressResponse, error) {
	resp := ShowAllowedIngressResponse{
		Ingresses: []pkg.IngressList{k.ingress},
	}

	return resp, nil
}

func (k *LookerIngressDriver) EnableUser(req EnableUserRequest) (EnableUserResponse, error) {
	log.Infof("Received Looker EnableUserRequest %v", req.User.Name)
	resp := EnableUserResponse{
		Ingress:    k.ingress,
		Identifier: req.User.Email,
	}

	client := pkg.GetLookerClient()

	users, searchErr := client.SearchUser(pkg.LookerSearchUserRequest{Email: req.User.Email})

	if len(users) == 0 {
		return resp, errors.New("You dont have a looker account. Please contact Looker admins")
	}

	for _, u := range users {
		if u.IsDisabled == false {
			return resp, errors.New("Looker user already present for user. Please contact admin")
		}
	}

	user := users[0] // lets create only for the 0th user. if there are multiple users, its bug on looker side

	patchedUser, patchErr := client.PatchUser(user.Id, pkg.LookerPatchUserRequest{IsDisabled: false})

	if patchErr != nil {
		return resp, patchErr
	}

	if patchedUser.IsDisabled == true {
		return resp, errors.New("Failed to enable user. please contact admin")
	}

	resp.UpdateStatusFlag = true

	if searchErr != nil {
		return resp, searchErr
	}

	return resp, nil
}

func (k *LookerIngressDriver) DisableUser(req DisableUserRequest) (DisableUserResponse, error) {
	log.Infof("Received Looker DisableUserRequest %v", req.LeaseIdentifier)
	response := DisableUserResponse{}

	client := pkg.GetLookerClient()

	users, searchErr := client.SearchUser(pkg.LookerSearchUserRequest{Email: req.LeaseIdentifier})

	if searchErr != nil {
		return DisableUserResponse{Ingress: k.ingress}, searchErr
	}

	for _, user := range users {
		if user.IsDisabled {
			continue
		}

		patchedUser, patchErr := client.PatchUser(user.Id, pkg.LookerPatchUserRequest{IsDisabled: true})

		if patchErr != nil || patchedUser.IsDisabled == false {
			return response, errors.New("failed to delete lease. contact looker admin")
		}
	}

	response.UpdateStatusFlag = true

	return response, nil
}

func (k *LookerIngressDriver) ShowIngressDetails(ShowIngressDetailsRequest) (ShowIngressDetailsResponse, error) {
	return ShowIngressDetailsResponse{Ingress: k.ingress}, nil
}

func (k *LookerIngressDriver) GetName() string {
	return "looker"
}
