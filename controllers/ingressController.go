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

//ShowAllowedIngress ...
func ShowAllowedIngress(c *gin.Context) {
	User, _ := c.Get("User")
	var data ingress_driver.ShowAllowedIngressResponse

	log.Infof("Listing securityGroups, ingresses, looker for user %s\n", User.(*models.Users).Email)

	for _, driver := range ingress_driver.GetEnabledIngressDrivers() {
		response, err := driver.ShowAllowedIngress()
		if err != nil {
			log.Errorf("Error listing resources for driver %s for user %s ", driver.GetName(), User)
		} else {
			switch driver.GetLeaseType() {
			case "looker":
				data.Looker = response.Looker
			case "aws":
				data.SecurityGroups = response.SecurityGroups
			default:
				data.Ingresses = response.Ingresses
			}
		}
	}
	c.HTML(http.StatusOK, "showingresslist.gohtml", gin.H{
		"data":  data,
		"user":  User,
		"token": csrf.Token(c.Request),
	})
}

//WhiteListIP ...
func WhiteListIP(c *gin.Context) {
	var leases []models.Leases
	var securityGroup config.SecurityGroupIngress
	User, _ := c.Get("User")

	driver := c.Param("driver")
	ns := c.Param("ns")
	name := c.Param("name")

	if driver == "aws" {
		securityGroup = config.SecurityGroupIngress{
			RuleType: c.PostForm("rule_type"),
		}
		switch securityGroup.RuleType {
		case "ssh":
			securityGroup.Protocol = "tcp"
			securityGroup.PortFrom = 22
			securityGroup.PortTo = 22
		case "https":
			securityGroup.Protocol = "tcp"
			securityGroup.PortFrom = 443
			securityGroup.PortTo = 443
		case "custom":
			securityGroup.Protocol = c.PostForm("protocol")
			securityGroup.PortFrom, _ = strconv.ParseInt(c.PostForm("port_from"), 10, 64)
			securityGroup.PortTo, _ = strconv.ParseInt(c.PostForm("port_to"), 10, 64)
		}
	}

	expiry, _ := strconv.Atoi(c.PostForm("expiry"))
	if expiry > config.AppCfg.MaxExpiry {
		c.SetCookie("message", "Expiry time is incorrect", 10, "/", "", config.AppCfg.CookieSecure, config.AppCfg.CookieHTTPOnly)
		c.Redirect(http.StatusFound, "/resources/"+driver+"/"+ns+"/"+name)
		return
	}
	leases = GetActiveLeases(driver, ns, name)

	showIngressDetailsResponse, err := ingress_driver.GetIngressDriverForNamespace(driver, ns).
		ShowIngressDetails(ingress_driver.ShowIngressDetailsRequest{Name: name})

	if err != nil {
		c.HTML(http.StatusOK, "manageingress.gohtml", gin.H{
			"data":   showIngressDetailsResponse,
			"user":   User,
			"driver": driver,
			"message": map[string]string{
				"class":   "Danger",
				"message": err.Error(),
			},
			"activeLeases": leases,
			"token":        csrf.Token(c.Request),
		})
		return
	}

	enableUserRequest := ingress_driver.EnableLeaseRequest{
		Name:          name,
		GinContext:    c,
		User:          User.(*models.Users),
		SecurityGroup: securityGroup,
	}

	enableUserResponse, enableUserErr := ingress_driver.GetIngressDriverForNamespace(driver, ns).EnableLease(enableUserRequest)

	if enableUserErr != nil {
		c.HTML(http.StatusOK, "manageingress.gohtml", gin.H{
			"data":   showIngressDetailsResponse,
			"user":   User,
			"driver": driver,
			"message": map[string]string{
				"class":   "Danger",
				"message": enableUserErr.Error(),
			},
			"activeLeases": leases,
			"token":        csrf.Token(c.Request),
		})
		return
	}

	if enableUserResponse.UpdateStatusFlag {
		msgInfo := "Whitelisted " + enableUserResponse.LeaseIdentifier + " to " + driver + " " + name + " in namespace " + ns + " for user " + User.(*models.Users).Email
		slackNotification(msgInfo, User.(*models.Users).Email)
		log.Info(msgInfo)
		if database.DB == nil {
			database.Conn()
		}

		lease := models.Leases{
			UserID:          User.(*models.Users).ID,
			LeaseIdentifier: enableUserResponse.LeaseIdentifier,
			LeaseType:       enableUserResponse.LeaseType,
			GroupID:         ns + ":" + name,
			Expiry:          uint(expiry),
		}

		if driver == "aws" {
			lease.Protocol = securityGroup.Protocol
			lease.PortFrom = strconv.FormatInt(securityGroup.PortFrom, 10)
			lease.PortTo = strconv.FormatInt(securityGroup.PortTo, 10)
		}

		database.DB.Create(&lease)

		leases = GetActiveLeases(driver, ns, name)

		c.HTML(http.StatusOK, "manageingress.gohtml", gin.H{
			"data":   showIngressDetailsResponse,
			"user":   User,
			"driver": driver,
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
		"data":   showIngressDetailsResponse,
		"user":   User,
		"driver": driver,
		"message": map[string]string{
			"class":   "Danger",
			"message": "Your IP/User is already present",
		},
		"activeLeases": leases,
		"token":        csrf.Token(c.Request),
	})
}

//DeleteIPFromIngress ...
func DeleteIPFromIngress(c *gin.Context) {
	var err error
	User, _ := c.Get("User")
	driver := c.Param("driver")
	ns := c.Param("ns")
	name := c.Param("name")
	leaseID, _ := strconv.Atoi(c.Param("id"))
	ID := uint(leaseID)
	leases := GetActiveLeases(driver, ns, name)

	showIngressDetailsResponse, err := ingress_driver.GetIngressDriverForNamespace(driver, ns).
		ShowIngressDetails(ingress_driver.ShowIngressDetailsRequest{Name: name})

	if err != nil {
		c.HTML(http.StatusOK, "manageingress.gohtml", gin.H{
			"data":   showIngressDetailsResponse,
			"user":   User,
			"driver": driver,
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
		err := errors.New("unauthorized, Trying to delete a lease of other user")
		log.Error("Error: ", err)
		c.HTML(http.StatusOK, "manageingress.gohtml", gin.H{
			"data":   showIngressDetailsResponse,
			"user":   User,
			"driver": driver,
			"message": map[string]string{
				"class":   "Danger",
				"message": err.Error(),
			},
			"activeLeases": leases,
			"token":        csrf.Token(c.Request),
		})
		return
	}
	leaseIdentifier := myCurrentLease.LeaseIdentifier

	resp, respErr := DeleteLeases(driver, ns, name, myCurrentLease, ID)
	if respErr != nil {
		c.HTML(http.StatusOK, "manageingress.gohtml", gin.H{
			"data":   showIngressDetailsResponse,
			"user":   User,
			"driver": driver,
			"message": map[string]string{
				"class":   "Danger",
				"message": respErr.Error(),
			},
			"activeLeases": leases,
			"token":        csrf.Token(c.Request),
		})
		return
	}
	if resp.UpdateStatusFlag {
		msgInfo := "Removed IP " + leaseIdentifier + " from " + driver + " " + name + " in namespace " + ns + " for user " + User.(*models.Users).Email
		slackNotification(msgInfo, User.(*models.Users).Email)
		log.Info(msgInfo)
		leases = GetActiveLeases(driver, ns, name)
		c.HTML(http.StatusOK, "manageingress.gohtml", gin.H{
			"data":         showIngressDetailsResponse,
			"user":         User,
			"driver":       driver,
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
		"data":   showIngressDetailsResponse,
		"user":   User,
		"driver": driver,
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

	User, _ := c.Get("User")
	driver := c.Param("driver")
	ns := c.Param("ns")
	name := c.Param("name")
	leases := GetActiveLeases(driver, ns, name)

	resp, err := ingress_driver.GetIngressDriverForNamespace(driver, ns).
		ShowIngressDetails(ingress_driver.ShowIngressDetailsRequest{Name: name})

	if err != nil {
		c.HTML(http.StatusNotFound, "manageingress.gohtml", gin.H{
			"message": map[string]string{
				"class":   "Danger",
				"message": err.Error(),
			},
			"user":         User,
			"driver":       driver,
			"activeLeases": leases,
			"token":        csrf.Token(c.Request),
		})
		return
	}
	message, cookieErr := c.Cookie("message")
	if cookieErr == nil {
		c.HTML(http.StatusOK, "manageingress.gohtml", gin.H{
			"data":         resp,
			"user":         User,
			"driver":       driver,
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
		"data":         resp,
		"user":         User,
		"driver":       driver,
		"activeLeases": leases,
		"token":        csrf.Token(c.Request),
	})
}

//GetActiveLeases ...
func GetActiveLeases(driver string, ns string, name string) []models.Leases {
	if database.DB == nil {
		database.Conn()
	}

	leases := []models.Leases{}

	leaseTypes := ingress_driver.GetLeaseTypes()

	query := database.DB.Preload("User").
		Where("lease_type in (?)", leaseTypes)

	if ns == "" && name == "" {
		query.Find(&leases)
	} else {
		query.Where(models.Leases{
			GroupID: ns + ":" + name,
		}).Find(&leases)
	}
	myleases := []models.Leases{}
	for i, lease := range leases {
		splitGroupID := strings.Split(lease.GroupID, ":")
		ns = splitGroupID[0]
		// Since old data has groupID without colon format.
		name = splitGroupID[0]
		if len(splitGroupID) > 1 {
			name = splitGroupID[1]
		}
		t := uint(lease.CreatedAt.Unix()) + lease.Expiry
		if t < uint(time.Now().Unix()) {
			leases[i].Expiry = uint(0)
			resp, err := DeleteLeases(driver, ns, name, lease, lease.ID)

			if resp.UpdateStatusFlag {
				log.Infof("Removed expired IP %s from %s %s in namespace %s for User %s\n", lease.LeaseIdentifier, driver, name, ns, lease.User.Email)
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

//DeleteLeases driver, ...
func DeleteLeases(driver string, ns string, name string, myCurrentLease models.Leases, ID uint) (ingress_driver.DisableLeaseResponse, error) {
	if database.DB == nil {
		database.Conn()
	}
	portFrom, _ := strconv.ParseInt(myCurrentLease.PortFrom, 10, 64)
	portTo, _ := strconv.ParseInt(myCurrentLease.PortTo, 10, 64)
	req := ingress_driver.DisableLeaseRequest{
		Name:            name,
		LeaseIdentifier: myCurrentLease.LeaseIdentifier,
		SecurityGroup: config.SecurityGroupIngress{
			Protocol: myCurrentLease.Protocol,
			PortFrom: portFrom,
			PortTo:   portTo,
		},
	}

	resp, err := ingress_driver.GetIngressDriverForNamespace(driver, ns).DisableLease(req)

	if resp.UpdateStatusFlag {
		database.DB.Delete(models.Leases{
			ID: ID,
		})
		log.Infof("Removing IP %s from database\n", myCurrentLease.LeaseIdentifier)
	}
	return resp, err
}

//ClearExpiredLeases ...
func ClearExpiredLeases(c *gin.Context) {
	GetActiveLeases("", "", "")
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
			"attachments": {
				payload,
			},
		},
	}
	payloads.SlackNotification(slackWebhookURL)
}
