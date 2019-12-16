package pkg

import (
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ingressList struct {
	Name           string
	Namespace      string
	Host           string
	Class          string
	WhitelistedIps []string
}

//MyClientSet ...
type MyClientSet struct {
	Clientset *kubernetes.Clientset
}

//GetIngresses ...
func (c MyClientSet) GetIngresses(ns string) (map[int]ingressList, error) {
	ingressClient := c.Clientset.ExtensionsV1beta1().Ingresses(ns)

	ingressLists, err := ingressClient.List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	myIngress := make(map[int]ingressList)

	for index, ingress := range ingressLists.Items {
		if _, ok := ingress.Annotations["concierge"]; ok && ingress.Annotations["concierge"] == "true" {
			var ingressHosts string
			for _, hosts := range ingress.Spec.Rules {
				if ingressHosts != "" {
					ingressHosts = ingressHosts + ", " + hosts.Host
				} else {
					ingressHosts = hosts.Host
				}
			}
			myIngress[index] = ingressList{
				ingress.Name,
				ingress.Namespace,
				ingressHosts,
				ingress.Annotations["kubernetes.io/ingress.class"],
				strings.Split(ingress.Annotations["traefik.ingress.kubernetes.io/whitelist-source-range"], ","),
			}
		}
	}
	return myIngress, nil
}

//RemoveIngressIP ...
func (c MyClientSet) RemoveIngressIP(ns string, ingressName string, ip string) (bool, error) {
	ingressClient := c.Clientset.ExtensionsV1beta1().Ingresses(ns)

	ingress, err := ingressClient.Get(ingressName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	if _, ok := ingress.Annotations["concierge"]; ok && ingress.Annotations["concierge"] == "true" {

		whitelistIps, ok := ingress.Annotations["traefik.ingress.kubernetes.io/whitelist-source-range"]
		if ok {
			whitelistIpsArr := strings.Split(whitelistIps, ",")
			for i := range whitelistIpsArr {
				whitelistIpsArr[i] = strings.TrimSpace(whitelistIpsArr[i])
			}
			var newWhitelistIpsArr []string
			for _, ips := range whitelistIpsArr {
				if ips != ip {
					newWhitelistIpsArr = append(newWhitelistIpsArr, ips)
				}
			}
			whitelistIps = strings.Join(newWhitelistIpsArr, ", ")
			ingress.Annotations["traefik.ingress.kubernetes.io/whitelist-source-range"] = whitelistIps
			_, updateErr := ingressClient.Update(ingress)

			if updateErr != nil {
				return false, updateErr
			}
			return true, nil
		}
	}
	return false, nil
}

//WhiteListIP ...
func (c MyClientSet) WhiteListIP(ns string, ingressName string, ip string) (bool, error) {
	ingressClient := c.Clientset.ExtensionsV1beta1().Ingresses(ns)

	ingress, err := ingressClient.Get(ingressName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	if _, ok := ingress.Annotations["concierge"]; ok && ingress.Annotations["concierge"] == "true" {

		whitelistIps, ok := ingress.Annotations["traefik.ingress.kubernetes.io/whitelist-source-range"]
		if ok {
			whitelistIpsArr := strings.Split(whitelistIps, ",")
			for i := range whitelistIpsArr {
				whitelistIpsArr[i] = strings.TrimSpace(whitelistIpsArr[i])
			}
			new := true
			for _, ips := range whitelistIpsArr {
				if ips == ip {
					log.Warn("Your IP is already present there")
					new = false
					break
				}
			}
			if new {
				whitelistIpsArr = append(whitelistIpsArr, ip)
				whitelistIps = strings.Join(whitelistIpsArr, ", ")
				ingress.Annotations["traefik.ingress.kubernetes.io/whitelist-source-range"] = whitelistIps
				_, updateErr := ingressClient.Update(ingress)

				if updateErr != nil {
					return false, updateErr
				}
				return true, nil
			}
			return false, nil
		}
	}
	return false, nil
}

//GetIngress ...
func (c MyClientSet) GetIngress(ns string, ingressName string) (ingressList, error) {
	ingressClient := c.Clientset.ExtensionsV1beta1().Ingresses(ns)
	ingress, err := ingressClient.Get(ingressName, metav1.GetOptions{})
	if err != nil {
		return ingressList{}, err
	}
	if _, ok := ingress.Annotations["concierge"]; ok && ingress.Annotations["concierge"] == "true" {
		var ingressHosts string
		for _, hosts := range ingress.Spec.Rules {
			if ingressHosts != "" {
				ingressHosts = ingressHosts + ", " + hosts.Host
			} else {
				ingressHosts = hosts.Host
			}
		}
		myIngress := ingressList{
			ingress.Name,
			ingress.Namespace,
			ingressHosts,
			ingress.Annotations["kubernetes.io/ingress.class"],
			[]string{},
		}
		return myIngress, nil
	}
	return ingressList{}, errors.New("Ingress not found")

}
