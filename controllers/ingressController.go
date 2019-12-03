package controllers

import (
	"concierge/config"
	"concierge/database"
	"concierge/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	layout = "2006-01-02 15:04:05 -0700 MST"
)

type info struct {
	Userinfo  models.Users
	Leaseinfo models.Leases
}

//ShowAllowedIngress ...
func ShowAllowedIngress(c *gin.Context) {
	clientset := config.KubeClient.ClientSet
	User, _ := c.Get("User")

	myclientset := myClientSet{clientset}
	ns := ""
	ns = c.Query("ns")
	data, err := myclientset.getIngresses(ns)
	if err != nil {
		log.Error("Error", err)
		return
	}
	c.HTML(http.StatusOK, "showingresslist.gohtml", gin.H{
		"data": data,
		"user": User,
	})
}

//WhiteListIP ...
func WhiteListIP(c *gin.Context) {
	clientset := config.KubeClient.ClientSet

	User, _ := c.Get("User")
	myclientset := myClientSet{clientset}
	ns := c.Param("ns")
	name := c.Param("name")
	expiry, _ := strconv.Atoi(c.PostForm("expiry"))
	data, err := myclientset.getIngress(ns, name)
	if err != nil {
		log.Error("Error", err)
		return
	}
	ip := c.Request.Header["X-Forwarded-For"][0]
	updateStatus, err := myclientset.whiteListIP(ns, name, ip)
	var leases []models.Leases
	if updateStatus {
		db, err := database.Conn()
		if err != nil {
			log.Error("Error", err)
		}
		defer db.Close()

		lease := models.Leases{
			UserID:    User.(*models.Users).ID,
			LeaseIP:   ip,
			LeaseType: "Ingress",
			GroupID:   ns + ":" + name,
			Expiry:    uint(expiry),
		}

		db.Create(&lease)
		leases = GetActiveLeases(ns, name)
		c.HTML(http.StatusOK, "manageingress.gohtml", gin.H{
			"data": data,
			"user": User,
			"message": map[string]string{
				"class":   "Success",
				"message": "Lease is successfully taken",
			},
			"activeLeases": leases,
		})
		return
	}
	if err != nil {
		log.Error("Error", err)
		return
	}
	leases = GetActiveLeases(ns, name)
	c.HTML(http.StatusOK, "manageingress.gohtml", gin.H{
		"data": data,
		"user": User,
		"message": map[string]string{
			"class":   "Danger",
			"message": "Your IP is already present",
		},
		"activeLeases": leases,
	})
}

//DeleteIPFromIngress ...
func DeleteIPFromIngress(c *gin.Context) {
	clientset := config.KubeClient.ClientSet

	User, _ := c.Get("User")

	myclientset := myClientSet{clientset}
	ns := c.Param("ns")
	name := c.Param("name")
	leaseID, err := strconv.Atoi(c.Param("id"))
	ID := uint(leaseID)
	data, err := myclientset.getIngress(ns, name)
	if err != nil {
		log.Error("Error", err)
		return
	}
	updateStatus, err := DeleteLeases(ns, name, c.Request.Header["X-Forwarded-For"][0], ID)
	leases := GetActiveLeases(ns, name)
	if updateStatus {
		c.HTML(http.StatusOK, "manageingress.gohtml", gin.H{
			"data":         data,
			"user":         User,
			"activeLeases": leases,
			"message": map[string]string{
				"class":   "Success",
				"message": "Lease is successfully deleted",
			},
		})
		return
	}
	if err != nil {
		log.Error("Error", err)
		return
	}
}

//IngressDetails ...
func IngressDetails(c *gin.Context) {
	clientset := config.KubeClient.ClientSet

	User, _ := c.Get("User")
	myclientset := myClientSet{clientset}
	ns := c.Param("ns")
	name := c.Param("name")
	leases := GetActiveLeases(ns, name)
	data, err := myclientset.getIngress(ns, name)
	if err != nil {
		log.Error("Error", err)
		c.HTML(http.StatusNotFound, "manageingress.gohtml", gin.H{
			"message": map[string]string{
				"class":   "Danger",
				"message": err.Error(),
			},
			"user":         User,
			"activeLeases": leases,
		})
		return
	}
	c.HTML(http.StatusOK, "manageingress.gohtml", gin.H{
		"data":         data,
		"user":         User,
		"activeLeases": leases,
	})
}

//GetActiveLeases ...
func GetActiveLeases(ns string, name string) []models.Leases {
	db, err := database.Conn()
	if err != nil {
		log.Error("Error", err)
	}
	defer db.Close()
	leases := []models.Leases{}
	if ns == "" && name == "" {
		db.Preload("User").Where(models.Leases{
			LeaseType: "Ingress",
		}).Find(&leases)
	} else {
		db.Preload("User").Where(models.Leases{
			LeaseType: "Ingress",
			GroupID:   ns + ":" + name,
		}).Find(&leases)
	}
	myleases := []models.Leases{}
	for i, lease := range leases {
		t := uint(lease.CreatedAt.Unix()) + lease.Expiry
		if t < uint(time.Now().Unix()) {
			leases[i].Expiry = uint(0)
			DeleteLeases(ns, name, lease.LeaseIP, lease.ID)
		} else {
			leases[i].Expiry = t - uint(time.Now().Unix())
			myleases = append(myleases, leases[i])
		}
	}
	return myleases
}

//DeleteLeases ...
func DeleteLeases(ns string, name string, ip string, ID uint) (bool, error) {
	clientset := config.KubeClient.ClientSet

	myclientset := myClientSet{clientset}
	db, err := database.Conn()
	if err != nil {
		log.Error("Error", err)
	}
	defer db.Close()

	db.Delete(models.Leases{
		ID: ID,
	})
	updateStatus, err := myclientset.removeIngressIP(ns, name, ip)
	return updateStatus, err
}

//ClearExpiredLeases ...
func ClearExpiredLeases(c *gin.Context) {
	GetActiveLeases("", "")
	c.String(200, "Done")
}
