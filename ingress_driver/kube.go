package ingress_driver

import (
	"concierge/config"
	"concierge/pkg"
	"errors"
	log "github.com/sirupsen/logrus"
	"strings"
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

func (k *KubeIngressDriver) EnableUser(req EnableUserRequest) (EnableUserResponse, error) {
	var data pkg.IngressList
	var err error
	var resp EnableUserResponse
	updateStatus := false

	errs := 0

	for kubeContext, kubeClient := range config.KubeClients {
		clientset := kubeClient.ClientSet
		myclientset := pkg.MyClientSet{Clientset: clientset}
		data, err = myclientset.GetIngress(kubeContext, req.Namespace, req.Name)
		if err != nil {
			errs = errs + 1
		}
		if data.Name != "" {
			resp.Ingress = data
			break
		}
	}

	if errs >= len(config.KubeClients) {
		return EnableUserResponse{}, err
	}

	ips := req.GinContext.Request.Header["X-Forwarded-For"][0]
	ip := strings.Split(ips, ",")[0]
	ip = ip + "/32"
	errs = 0

	for _, kubeClient := range config.KubeClients {
		clientset := kubeClient.ClientSet
		myclientset := pkg.MyClientSet{Clientset: clientset}
		updateStatus, err = myclientset.WhiteListIP(req.Namespace, req.Name, ip)
		if err != nil {
			errs = errs + 1
		}
		if updateStatus {
			resp.UpdateStatusFlag = true

		}
	}
	resp.IdentifierType = "Ingress"
	resp.Identifier = ip

	if errs >= len(config.KubeClients) {
		return EnableUserResponse{}, errors.New("Your IP is already there")
	}

	return resp, nil
}

func (k *KubeIngressDriver) DisableUser(req DisableUserRequest) (DisableUserResponse, error) {
	return DisableUserResponse{}, nil
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
