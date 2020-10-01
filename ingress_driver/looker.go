package ingress_driver

import (
	"concierge/config"
	"concierge/pkg"
	"errors"
	log "github.com/sirupsen/logrus"
	"strings"
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
		Ingresses: []pkg.IngressList{k.ingress},
	}

	return resp, nil
}

func (k *LookerIngressDriver) EnableLease(req EnableLeaseRequest) (EnableLeaseResponse, error) {
	log.Infof("Received Looker EnableLeaseRequest %v", req.User.Name)
	resp := EnableLeaseResponse{
		Ingress:         k.ingress,
		LeaseIdentifier: req.User.Email,
		LeaseType:       k.GetLeaseType(),
	}

	client := pkg.GetLookerClient()

	users, searchErr := client.SearchUser(pkg.LookerSearchUserRequest{Email: req.User.Email})

	if searchErr != nil {
		return resp, searchErr
	}

	if len(users) == 0 {
		createUserResponse, err := k.createUser(req.User.Email)

		if err != nil {
			return resp, err
		}

		if createUserResponse.IsDisabled == true {
			return resp, errors.New("failed to create user on looker. please contact looker admins")
		}

		resp.UpdateStatusFlag = true

		return resp, nil
	}

	if len(users) > 1 {
		return resp, errors.New("You have multiple looker accounts. Please contact looker admins")
	}

	user := users[0]

	if user.IsDisabled == false {
		return resp, errors.New("Looker user already enabled for user. Please visit looker.")
	}

	patchedUser, patchErr := client.PatchUser(user.Id, pkg.LookerPatchUserRequest{IsDisabled: false})

	if patchErr != nil {
		return resp, patchErr
	}

	if patchedUser.IsDisabled == true {
		return resp, errors.New("Failed to enable user. please contact admin")
	}

	resp.UpdateStatusFlag = true

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
	return Looker
}

func (k *LookerIngressDriver) GetLeaseType() string {
	return Looker
}

func (k *LookerIngressDriver) isEnabled() bool {
	return config.LookerConfig.IsEnabled
}

func (k *LookerIngressDriver) createUser(email string) (*pkg.LookerCreateUserResponse, error) {

	client := pkg.GetLookerClient()

	firstName, lastName, err := parseFirstAndLastNameFromEmail(email)

	if err != nil {
		return nil, err
	}

	req := pkg.LookerCreateUserRequest{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	log.Infof("Creating looker account for email %s", email)
	resp, err := client.CreateUser(req)

	if err != nil {
		log.Errorf("failed to create looker account for email %s : %s", email, err.Error())
		return nil, err
	}

	log.Infof("created looker account for email:%s id:%s", email, resp.Id)
	log.Infof("creating looker user credentials email for email:%s", email)

	_, createUserCredEmailErr := client.CreateUserCredentialsEmail(resp.Id, pkg.LookerCreateCredentialsEmailRequest{
		Email: email,
		Type:  "email",
	})

	if createUserCredEmailErr != nil {
		log.Errorf("failed to create looker user credentials email for email %s : %s", email, createUserCredEmailErr.Error())
		return nil, createUserCredEmailErr
	}

	return resp, err
}

func parseFirstAndLastNameFromEmail(email string) (string, string, error) {
	emailSplit := strings.Split(email, "@")
	username := emailSplit[0]
	usernameSplit := strings.Split(username, ".")

	switch len(usernameSplit) {
	case 2: // email of the form "abc.xyz@razorpay.com"
		return usernameSplit[0], usernameSplit[1], nil
	case 1: // email of the form "abc@razorpay.com"
		return usernameSplit[0], "undefined", nil // if we dont send last name, email id doesnt get set on looker
	default:
		return "", "", errors.New("no looker account associated with this email id. unable to create a looker account automatically. contact looker admins")
	}
}
