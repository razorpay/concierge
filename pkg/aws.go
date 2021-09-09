package pkg

import (
	"concierge/config"
	"concierge/constants"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
)

var awsclient *AWSClientSession

type SecurityGroupList struct {
	Name                 string
	GroupId              string
	Description          string
	SecurityGroupIngress *config.SecurityGroupIngress
	VpcId                string
	WhitelistedIps       []string
}

//AWSClientSession ...
type AWSClientSession struct {
	session *session.Session
}

func GetAWSSession() *AWSClientSession {
	if awsclient == nil {
		awsclient = &AWSClientSession{
			session: session.Must(session.NewSession()),
		}
	}
	return awsclient
}

//GetSecurityGroups ...
func (c *AWSClientSession) GetSecurityGroups() (*ec2.DescribeSecurityGroupsOutput, error) {
	svc := ec2.New(c.session)
	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String(constants.AWSTagName),
				Values: []*string{
					aws.String("true"),
				},
			},
		},
	}

	result, err := svc.DescribeSecurityGroups(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Error("Error: ", aerr.Error())
				return nil, errors.New(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Error("Error: ", err.Error())
			return nil, errors.New(aerr.Error())
		}
	}
	return result, nil
}

func (c *AWSClientSession) GetSecurityGroupDetails(groupId string) (*ec2.DescribeSecurityGroupsOutput, error) {
	svc := ec2.New(c.session)
	input := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{
			aws.String(groupId),
		},
	}
	result, err := svc.DescribeSecurityGroups(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Error("Error: ", aerr.Error())
				return nil, errors.New(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Error("Error: ", err.Error())
			return nil, errors.New(aerr.Error())
		}
	}
	return result, nil
}

func (c *AWSClientSession) WhitelistIP(groupId string, ip string, securityGroup config.SecurityGroupIngress, email string) error {
	svc := ec2.New(c.session)
	time := time.Now().Format("01-02-2006 15:04:05")
	input := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: aws.String(groupId),
		IpPermissions: []*ec2.IpPermission{
			{
				FromPort:   aws.Int64(securityGroup.PortFrom),
				IpProtocol: aws.String(securityGroup.Protocol),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp:      aws.String(ip),
						Description: aws.String("Created by " + email + " at " + time),
					},
				},
				ToPort: aws.Int64(securityGroup.PortTo),
			},
		},
	}
	_, err := svc.AuthorizeSecurityGroupIngress(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Error("Error: ", aerr.Error())
				return errors.New(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Error("Error: ", err.Error())
			return errors.New(aerr.Error())
		}
	}
	return nil
}

func (c *AWSClientSession) RevokeIP(groupId string, ip string, securityGroup config.SecurityGroupIngress) error {
	svc := ec2.New(c.session)
	input := &ec2.RevokeSecurityGroupIngressInput{
		GroupId: aws.String(groupId),
		IpPermissions: []*ec2.IpPermission{
			{
				FromPort:   aws.Int64(securityGroup.PortFrom),
				IpProtocol: aws.String(securityGroup.Protocol),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp: aws.String(ip),
					},
				},
				ToPort: aws.Int64(securityGroup.PortTo),
			},
		},
	}
	_, err := svc.RevokeSecurityGroupIngress(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Error("Error: ", aerr.Error())
				return errors.New(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Error("Error: ", err.Error())
			return errors.New(aerr.Error())
		}
	}
	return nil
}
