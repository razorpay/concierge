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
	count := 0

	namespaces := make(map[string]int)

	resp := ShowAllowedIngressResponse{
		Ingresses: []pkg.IngressList{},
	}

	log.Infof("Listing ingress in namespace %s for user %s\n", ns, User.Email)

	for kubeContext, kubeClient := range config.KubeClients {
		clientset := kubeClient.ClientSet
		myclientset := pkg.MyClientSet{Clientset: clientset}
		data, _ := myclientset.GetIngresses(kubeContext, ns)
		for _, ingress := range data {
			if val, ok := namespaces[ingress.Namespace+":"+ingress.Name]; ok {
				resp.Ingresses[val].Context = resp.Ingresses[val].Context + "," + ingress.Context
				continue
			}
			namespaces[ingress.Namespace+":"+ingress.Name] = count
			resp.Ingresses = append(resp.Ingresses, ingress)
			count = count + 1
		}
	}
	return resp, nil
}

func (k *KubeIngressDriver) EnableUser(req EnableUserRequest) (EnableUserResponse, error) {
	var err error
	resp := EnableUserResponse{}
	updateStatus := false

	errs := 0

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

	resp.Identifier = ip

	if errs >= len(config.KubeClients) {
		return EnableUserResponse{}, errors.New("Your IP is already there")
	}

	return resp, nil
}

func (k *KubeIngressDriver) DisableUser(req DisableUserRequest) (DisableUserResponse, error) {
	dbflag := false
	resp := DisableUserResponse{UpdateStatusFlag: true}
	var err error
	errs := 0

	for _, kubeClient := range config.KubeClients {
		clientset := kubeClient.ClientSet
		myclientset := pkg.MyClientSet{Clientset: clientset}
		_, dbflag, err = myclientset.RemoveIngressIP(req.Namespace, req.Name, req.LeaseIdentifier)
		if err != nil {
			errs = errs + 1
		}
		if dbflag {
			resp.UpdateStatusFlag = false
		}
	}

	if errs >= len(config.KubeClients) {
		return resp, err
	} else {
		return resp, nil
	}
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
