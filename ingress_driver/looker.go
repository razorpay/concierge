package ingress_driver

import (
	"concierge/pkg"
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

	// todo make call to looker

	return EnableUserResponse{
		UpdateStatusFlag: true,
		Ingress:          k.ingress,
		Identifier:       req.User.Email,
	}, nil
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
