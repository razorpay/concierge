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
	log.Infof("Received EnableUserRequest %v", req.User.Name)
	resp := EnableUserResponse{
		Ingress:    k.ingress,
		Identifier: req.User.Email,
	}

	client := pkg.GetLookerClient()

	users, err := client.SearchUser(pkg.LookerSearchUserRequest{Email: req.User.Email})

	if len(users) == 0 {
		// todo: add user creation flow here
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

	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (k *LookerIngressDriver) DisableUser(req DisableUserRequest) (DisableUserResponse, error) {

	// todo make call to looker

	return DisableUserResponse{
		UpdateStatusFlag: true,
	}, nil
}

func (k *LookerIngressDriver) ShowIngressDetails(ShowIngressDetailsRequest) (ShowIngressDetailsResponse, error) {
	return ShowIngressDetailsResponse{Ingress: k.ingress}, nil
}

func (k *LookerIngressDriver) GetName() string {
	return "looker"
}
