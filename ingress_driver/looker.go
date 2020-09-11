package ingress_driver

import "concierge/pkg"

func (k *LookerIngressDriver) ShowAllowedIngress(ShowAllowedIngressRequest) (ShowAllowedIngressResponse, error) {
	resp := ShowAllowedIngressResponse{
		Ingresses: []pkg.IngressList{{
			Name:           "looker",
			Namespace:      "looker",
			Context:        "default",
			Host:           "default",
			Class:          "default",
			WhitelistedIps: nil,
		}},
	}

	return resp, nil
}

func (k *LookerIngressDriver) Whitelist() {

}
func (k *LookerIngressDriver) Terminate() {

}
func (k *LookerIngressDriver) ShowIngressDetails(ShowIngressDetailsRequest) (ShowIngressDetailsResponse, error) {

	return ShowIngressDetailsResponse{
		Ingress: pkg.IngressList{
			Name:           "looker",
			Namespace:      "looker",
			Context:        "default",
			Host:           "https://looker.razorpay.com/",
			Class:          "default",
			WhitelistedIps: nil,
		},
	}, nil
}

func (k *LookerIngressDriver) GetName() string {
	return "looker"
}
