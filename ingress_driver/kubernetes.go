package ingress_driver

import (
	"concierge/config"
	"concierge/pkg"
	"errors"
	"strings"
)

type KubeIngressDriver struct {
	Namespace string
}

func getKubernetesIngressDriver(namespace string) IngressDriver {
	return &KubeIngressDriver{namespace}
}

func (k *KubeIngressDriver) ShowAllowedIngress() (ShowAllowedIngressResponse, error) {
	ns := k.Namespace
	count := 0

	namespaces := make(map[string]int)

	resp := ShowAllowedIngressResponse{
		Ingresses: []pkg.IngressList{},
	}

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

func (k *KubeIngressDriver) EnableLease(req EnableLeaseRequest) (EnableLeaseResponse, error) {
	var err error
	resp := EnableLeaseResponse{}
	updateStatus := false

	errs := 0

	ips := req.GinContext.Request.Header["X-Forwarded-For"][0]
	ip := strings.Split(ips, ",")[0]
	ip = ip + "/32"
	errs = 0

	for _, kubeClient := range config.KubeClients {
		clientset := kubeClient.ClientSet
		myclientset := pkg.MyClientSet{Clientset: clientset}
		updateStatus, err = myclientset.WhiteListIP(k.Namespace, req.Name, ip)
		if err != nil {
			errs = errs + 1
		}
		if updateStatus {
			resp.UpdateStatusFlag = true

		}
	}

	resp.LeaseIdentifier = ip
	resp.LeaseType = "Ingress"

	if errs >= len(config.KubeClients) {
		return EnableLeaseResponse{}, errors.New("Your IP is already there")
	}

	return resp, nil
}

func (k *KubeIngressDriver) DisableLease(req DisableLeaseRequest) (DisableLeaseResponse, error) {
	dbflag := false
	resp := DisableLeaseResponse{UpdateStatusFlag: true}
	var err error
	errs := 0

	for _, kubeClient := range config.KubeClients {
		clientset := kubeClient.ClientSet
		myclientset := pkg.MyClientSet{Clientset: clientset}
		_, dbflag, err = myclientset.RemoveIngressIP(k.Namespace, req.Name, req.LeaseIdentifier)
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
		data, err = myclientset.GetIngress(kubeContext, k.Namespace, req.Name)

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
	return "kubernetes"
}

func (k *KubeIngressDriver) isEnabled() bool {
	return true
}
