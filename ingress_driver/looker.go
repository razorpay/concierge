package ingress_driver

import "concierge/pkg"

func (k *LookerIngressDriver) ShowAllowedIngress(ShowAllowedIngressRequest) (ShowAllowedIngressResponse, error) {
	resp := ShowAllowedIngressResponse{
		Ingresses: []pkg.IngressList{{
			Name:           "Looker",
			Namespace:      "Looker",
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
func (k *LookerIngressDriver) ShowIngressDetails() {

}

func (k *LookerIngressDriver) GetName() string {
	return "looker"
}
