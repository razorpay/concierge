package controllers

import (
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

type myClientSet struct {
	clientset *kubernetes.Clientset
}

func (c myClientSet) getIngresses(ns string) (map[int]ingressList, error) {
	ingressClient := c.clientset.ExtensionsV1beta1().Ingresses(ns)
	log.Infof("Listing ingress in namespace %s:\n", ns)

	ingressLists, err := ingressClient.List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	myIngress := make(map[int]ingressList)

	for index, ingress := range ingressLists.Items {
		if _, ok := ingress.Annotations["traefik.ingress.kubernetes.io/whitelist-source-range"]; ok {
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

func (c myClientSet) removeIngressIP(ns string, ingressName string, ip string) (bool, error) {
	ingressClient := c.clientset.ExtensionsV1beta1().Ingresses(ns)
	log.Infof("Removing IP %s from ingress %s in namespace %s:\n", ip, ingressName, ns)

	ingress, err := ingressClient.Get(ingressName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
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
	return false, nil
}

func (c myClientSet) whiteListIP(ns string, ingressName string, ip string) (bool, error) {
	ingressClient := c.clientset.ExtensionsV1beta1().Ingresses(ns)
	log.Infof("Whitelisting IP %s to ingress %s in namespace %s:\n", ip, ingressName, ns)

	ingress, err := ingressClient.Get(ingressName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
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
	return false, nil
}

func (c myClientSet) getIngress(ns string, ingressName string) (ingressList, error) {
	ingressClient := c.clientset.ExtensionsV1beta1().Ingresses(ns)
	ingress, err := ingressClient.Get(ingressName, metav1.GetOptions{})
	if err != nil {
		return ingressList{}, err
	}
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
