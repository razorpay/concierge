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
func (k *KubeIngressDriver) ShowIngressDetails() {

}

func (k *KubeIngressDriver) GetName() string {
	return "kubeClient"
}
