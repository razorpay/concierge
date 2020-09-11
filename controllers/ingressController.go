package controllers

import (
	"concierge/config"
	"concierge/database"
	"concierge/ingress_driver"
	"concierge/models"
	"concierge/pkg"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	log "github.com/sirupsen/logrus"
)

type info struct {
	Userinfo  models.Users
	Leaseinfo models.Leases
}

type ingressess []string

//ShowAllowedIngress ...
func ShowAllowedIngress(c *gin.Context) {
	User, _ := c.Get("User")
	ns, count := "", 0
	ns = c.Query("ns")
	var myIngress []pkg.IngressList
	namespaces := make(map[string]int)
	data := []pkg.IngressList{}

	req := ingress_driver.ShowAllowedIngressRequest{
		User:      User.(*models.Users),
		Namespace: ns,
	}
	for _, driver := range ingress_driver.GetIngressDrivers() {
		response, err := driver.ShowAllowedIngress(req)
		if err != nil {
			log.Errorf("Error listing ingresses for driver %s for user %s ", driver.GetName(), req.User)
		}
		data = append(data, response.Ingresses...)
	}

	for _, ingress := range data {
		if val, ok := namespaces[ingress.Namespace+":"+ingress.Name]; ok {
			myIngress[val].Context = myIngress[val].Context + "," + ingress.Context
			continue
		}
		namespaces[ingress.Namespace+":"+ingress.Name] = count
		myIngress = append(myIngress, ingress)
		count = count + 1
	}

	c.HTML(http.StatusOK, "showingresslist.gohtml", gin.H{
		"data":  myIngress,
		"user":  User,
		"token": csrf.Token(c.Request),
	})
}

//WhiteListIP ...
func WhiteListIP(c *gin.Context) {
	var leases []models.Leases
	errs := 0
	var err error
	var myIngress, data pkg.IngressList
	updateStatusflag, updateStatus := false, true
	User, _ := c.Get("User")
	ns := c.Param("ns")
	name := c.Param("name")

	expiry, _ := strconv.Atoi(c.PostForm("expiry"))
	if expiry > config.AppCfg.MaxExpiry {
		c.SetCookie("message", "Expiry time is incorrect", 10, "/", "", config.AppCfg.CookieSecure, config.AppCfg.CookieHTTPOnly)
		c.Redirect(http.StatusFound, "/ingress/"+ns+"/"+name)
		return
	}
	leases = GetActiveLeases(ns, name)
	for kubeContext, kubeClient := range config.KubeClients {
		clientset := kubeClient.ClientSet
		myclientset := pkg.MyClientSet{Clientset: clientset}
		data, err = myclientset.GetIngress(kubeContext, ns, name)
		if err != nil {
			errs = errs + 1
		}
		if data.Name != "" {
			myIngress = data
			break
		}
	}

	if errs >= len(config.KubeClients) {
		c.HTML(http.StatusOK, "manageingress.gohtml", gin.H{
			"data": myIngress,
			"user": User,
			"message": map[string]string{
				"class":   "Danger",
				"message": err.Error(),
			},
			"activeLeases": leases,
			"token":        csrf.Token(c.Request),
		})
		return
	}
	ips := c.Request.Header["X-Forwarded-For"][0]
	ip := strings.Split(ips, ",")[0]
	ip = ip + "/32"
	errs = 0
	for _, kubeClient := range config.KubeClients {
		clientset := kubeClient.ClientSet
		myclientset := pkg.MyClientSet{Clientset: clientset}
		updateStatus, err = myclientset.WhiteListIP(ns, name, ip)
		if err != nil {
			errs = errs + 1
		}
		if updateStatus {
			updateStatusflag = true
		}
	}

	if errs >= len(config.KubeClients) {
		c.HTML(http.StatusOK, "manageingress.gohtml", gin.H{
			"data": myIngress,
			"user": User,
			"message": map[string]string{
				"class":   "Danger",
				"message": "Your IP is already there",
			},
			"activeLeases": leases,
			"token":        csrf.Token(c.Request),
		})
		return
	}
	if updateStatusflag {
		msgInfo := "Whitelisted IP " + ip + " to ingress " + name + " in namespace " + ns + " for user " + User.(*models.Users).Email
		slackNotification(msgInfo, User.(*models.Users).Email)
		log.Info(msgInfo)
		if database.DB == nil {
			database.Conn()
		}

		lease := models.Leases{
			UserID:    User.(*models.Users).ID,
			LeaseIP:   ip,
			LeaseType: "Ingress",
			GroupID:   ns + ":" + name,
			Expiry:    uint(expiry),
		}

		database.DB.Create(&lease)
		leases = GetActiveLeases(ns, name)

		c.HTML(http.StatusOK, "manageingress.gohtml", gin.H{
			"data": myIngress,
			"user": User,
			"message": map[string]string{
				"class":   "Success",
				"message": "Lease is successfully taken",
			},
			"activeLeases": leases,
			"token":        csrf.Token(c.Request),
		})
		return
	}

	c.HTML(http.StatusOK, "manageingress.gohtml", gin.H{
		"data": myIngress,
		"user": User,
		"message": map[string]string{
			"class":   "Danger",
			"message": "Your IP is already present",
		},
		"activeLeases": leases,
		"token":        csrf.Token(c.Request),
	})
}

//DeleteIPFromIngress ...
func DeleteIPFromIngress(c *gin.Context) {
	errs := 0
	var err error
	updateStatusflag := false
	User, _ := c.Get("User")
	ns := c.Param("ns")
	name := c.Param("name")
	leaseID, err := strconv.Atoi(c.Param("id"))
	ID := uint(leaseID)
	leases := GetActiveLeases(ns, name)
	var myIngress, data pkg.IngressList

	for kubeContext, kubeClient := range config.KubeClients {
		clientset := kubeClient.ClientSet
		myclientset := pkg.MyClientSet{Clientset: clientset}
		data, err = myclientset.GetIngress(kubeContext, ns, name)
		if err != nil {
			errs = errs + 1
		}
		if data.Name != "" {
			myIngress = data
			break
		}
	}

	if errs >= len(config.KubeClients) {
		c.HTML(http.StatusOK, "manageingress.gohtml", gin.H{
			"data": myIngress,
			"user": User,
			"message": map[string]string{
				"class":   "Danger",
				"message": err.Error(),
			},
			"activeLeases": leases,
			"token":        csrf.Token(c.Request),
		})
		return
	}

	if database.DB == nil {
		database.Conn()
	}
	myCurrentLease := models.Leases{}
	database.DB.Where(models.Leases{
		ID: ID,
	}).Find(&myCurrentLease)
	if myCurrentLease.UserID != User.(*models.Users).ID {
		err := errors.New("Unauthorized, Trying to delete a lease of other user")
		log.Error("Error: ", err)
		c.HTML(http.StatusOK, "manageingress.gohtml", gin.H{
			"data": myIngress,
			"user": User,
			"message": map[string]string{
				"class":   "Danger",
				"message": err.Error(),
			},
			"activeLeases": leases,
			"token":        csrf.Token(c.Request),
		})
		return
	}
	ip := myCurrentLease.LeaseIP

	updateStatusflag, errs, err = DeleteLeases(ns, name, ip, ID)
	if errs >= len(config.KubeClients) {
		c.HTML(http.StatusOK, "manageingress.gohtml", gin.H{
			"data": myIngress,
			"user": User,
			"message": map[string]string{
				"class":   "Danger",
				"message": err.Error(),
			},
			"activeLeases": leases,
			"token":        csrf.Token(c.Request),
		})
		return
	}
	if updateStatusflag {
		msgInfo := "Removed IP " + ip + " from ingress " + name + " in namespace " + ns + " for user " + User.(*models.Users).Email
		slackNotification(msgInfo, User.(*models.Users).Email)
		log.Info(msgInfo)
		leases = GetActiveLeases(ns, name)
		c.HTML(http.StatusOK, "manageingress.gohtml", gin.H{
			"data":         myIngress,
			"user":         User,
			"activeLeases": leases,
			"message": map[string]string{
				"class":   "Success",
				"message": "Lease is successfully deleted",
			},
			"token": csrf.Token(c.Request),
		})
		return
	}
	c.HTML(http.StatusOK, "manageingress.gohtml", gin.H{
		"data": myIngress,
		"user": User,
		"message": map[string]string{
			"class":   "Danger",
			"message": "There is some error in deleting your IP, Try again or contact admin",
		},
		"activeLeases": leases,
		"token":        csrf.Token(c.Request),
	})
}

//IngressDetails ...
func IngressDetails(c *gin.Context) {
	errs := 0
	var err error
	User, _ := c.Get("User")
	ns := c.Param("ns")
	name := c.Param("name")
	leases := GetActiveLeases(ns, name)
	var myIngress, data pkg.IngressList

	for kubeContext, kubeClient := range config.KubeClients {
		clientset := kubeClient.ClientSet
		myclientset := pkg.MyClientSet{Clientset: clientset}
		data, err = myclientset.GetIngress(kubeContext, ns, name)

		if data.Name != "" {
			myIngress = data
			break
		}
		if err != nil {
			errs = errs + 1
		}
	}

	if errs >= len(config.KubeClients) {
		c.HTML(http.StatusNotFound, "manageingress.gohtml", gin.H{
			"message": map[string]string{
				"class":   "Danger",
				"message": err.Error(),
			},
			"user":         User,
			"activeLeases": leases,
			"token":        csrf.Token(c.Request),
		})
		return
	}
	message, cookieErr := c.Cookie("message")
	if cookieErr == nil {
		c.HTML(http.StatusOK, "manageingress.gohtml", gin.H{
			"data":         myIngress,
			"user":         User,
			"activeLeases": leases,
			"token":        csrf.Token(c.Request),
			"message": map[string]string{
				"class":   "Danger",
				"message": message,
			},
		})
		return
	}
	c.HTML(http.StatusOK, "manageingress.gohtml", gin.H{
		"data":         myIngress,
		"user":         User,
		"activeLeases": leases,
		"token":        csrf.Token(c.Request),
	})
}

//GetActiveLeases ...
func GetActiveLeases(ns string, name string) []models.Leases {
	if database.DB == nil {
		database.Conn()
	}

	leases := []models.Leases{}
	if ns == "" && name == "" {
		database.DB.Preload("User").Where(models.Leases{
			LeaseType: "Ingress",
		}).Find(&leases)
	} else {
		database.DB.Preload("User").Where(models.Leases{
			LeaseType: "Ingress",
			GroupID:   ns + ":" + name,
		}).Find(&leases)
	}
	myleases := []models.Leases{}
	for i, lease := range leases {
		splitGroupID := strings.Split(lease.GroupID, ":")
		ns = splitGroupID[0]
		name = splitGroupID[1]
		t := uint(lease.CreatedAt.Unix()) + lease.Expiry
		if t < uint(time.Now().Unix()) {
			leases[i].Expiry = uint(0)
			updateStatusflag, _, err := DeleteLeases(ns, name, lease.LeaseIP, lease.ID)

			if updateStatusflag {
				log.Infof("Removed expired IP %s from ingress %s in namespace %s for User %s\n", lease.LeaseIP, name, ns, lease.User.Email)
			} else {
				log.Error("Error: ", err)
			}
		} else {
			leases[i].Expiry = t - uint(time.Now().Unix())
			myleases = append(myleases, leases[i])
		}
	}
	return myleases
}

//DeleteLeases ...
func DeleteLeases(ns string, name string, ip string, ID uint) (bool, int, error) {
	updateStatusflag, dbflag := true, false
	var err error
	errs := 0
	if database.DB == nil {
		database.Conn()
	}

	for _, kubeClient := range config.KubeClients {
		clientset := kubeClient.ClientSet
		myclientset := pkg.MyClientSet{Clientset: clientset}
		_, dbflag, err = myclientset.RemoveIngressIP(ns, name, ip)
		if err != nil {
			errs = errs + 1
		}
		if dbflag {
			updateStatusflag = false
		}
	}

	if updateStatusflag {
		database.DB.Delete(models.Leases{
			ID: ID,
		})
		log.Infof("Removing IP %s from database\n", ip)
	}
	return updateStatusflag, errs, err
}

//ClearExpiredLeases ...
func ClearExpiredLeases(c *gin.Context) {
	GetActiveLeases("", "")
	c.String(200, "Done")
}

func slackNotification(msg string, user string) {
	slackWebhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	if slackWebhookURL == "" {
		return
	}
	payload := pkg.Payload{
		Title:      "Concierge",
		Pretext:    msg,
		Text:       msg,
		Color:      "#36a64f",
		AuthorName: user,
		TitleLink:  "",
		Footer:     "Concierge",
		Timestamp:  strconv.FormatInt(time.Now().Unix(), 10),
	}
	payloads := pkg.Payloads{
		Attachments: map[string][]pkg.Payload{
			"attachments": []pkg.Payload{
				payload,
			},
		},
	}
	payloads.SlackNotification(slackWebhookURL)
}
