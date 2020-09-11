package ingress_driver

import (
	"concierge/config"
	"concierge/pkg"
	log "github.com/sirupsen/logrus"
)

func (k *KubeIngressDriver) ShowAllowedIngress(req ShowAllowedIngressRequest) (ShowAllowedIngressResponse, error) {
	ns := req.Namespace
	User := req.User

	resp := ShowAllowedIngressResponse{
		Ingresses: []pkg.IngressList{},
	}

	log.Infof("Listing ingress in namespace %s for user %s\n", ns, User.Email)

	var data []pkg.IngressList
	for kubeContext, kubeClient := range config.KubeClients {
		clientset := kubeClient.ClientSet
		myclientset := pkg.MyClientSet{Clientset: clientset}
		d, _ := myclientset.GetIngresses(kubeContext, ns)
		data = append(data, d...)
	}
	return resp, nil
}

func (k *KubeIngressDriver) Whitelist() {

}
func (k *KubeIngressDriver) Terminate() {

}
func (k *KubeIngressDriver) ShowIngressDetails(req ShowIngressDetailsRequest) (ShowIngressDetailsResponse, error) {
	var myIngress, data pkg.IngressList
	var err error
	errs := 0

	resp := ShowIngressDetailsResponse{}

	for kubeContext, kubeClient := range config.KubeClients {
		clientset := kubeClient.ClientSet
		myclientset := pkg.MyClientSet{Clientset: clientset}
		data, err = myclientset.GetIngress(kubeContext, req.Namespace, req.Name)

		if data.Name != "" {
			myIngress = data
			break
		}
		if err != nil {
			errs = errs + 1
		}
	}
	resp.Ingress = myIngress

	if errs >= len(config.KubeClients) {
		return resp, err
	}
	return resp, nil
}

func (k *KubeIngressDriver) GetName() string {
	return "kubeClient"
}
