package pkg

import (
	"concierge/config"
	"concierge/constants"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
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

type S3BucketList struct {
	Name           string
	WhitelistedIps []string
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

// ============================ SECURITY GROUPS START ============================================
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
			log.Error("Error: ", err.Error())
			return nil, err
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
			log.Error("Error: ", err.Error())
			return nil, err
		}
	}
	return result, nil
}

func (c *AWSClientSession) WhitelistIPInSecurityGroup(groupId string, ip string, securityGroup config.SecurityGroupIngress, email string) error {
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
			case "InvalidPermission.Duplicate":
				mesg := "lease creation failed! Does a similar lease already exist? Terminate that first"
				log.Error("Error: " + mesg)
				return errors.New(mesg)
			default:
				log.Error("Error: ", aerr.Error())
				return errors.New(aerr.Error())
			}
		} else {
			log.Error("Error: ", aerr.Error())
			return err
		}
	}
	return nil
}

func (c *AWSClientSession) RevokeIPFromSecurityGroup(groupId string, ip string, securityGroup config.SecurityGroupIngress) error {
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
			log.Error("Error: ", err.Error())
			return err
		}
	}
	return nil
}

// ============================ SECURITY GROUPS END ============================================

// ============================ S3BUCKET START =================================================
//GetS3Buckets ...
func (c *AWSClientSession) GetS3Buckets() (*s3.ListBucketsOutput, error) {
	var buckets []*s3.Bucket

	if config.AWSS3Buckets != "" {
		for _, bucketName := range strings.Split(config.AWSS3Buckets, ",") {
			buckets = append(buckets, &s3.Bucket{
				Name: aws.String(bucketName),
			})
		}
		return &s3.ListBucketsOutput{Buckets: buckets}, nil
	}
	return nil, errors.New("buckets not found")

	// Fallback to default method of fetching buckets from AWS via SDK
	// svc := s3.New(c.session)
	// input := &s3.ListBucketsInput{}

	// result, err := svc.ListBuckets(input)
	// if err != nil {
	// 	if aerr, ok := err.(awserr.Error); ok {
	// 		switch aerr.Code() {
	// 		default:
	// 			log.Error("Error: ", aerr.Error())
	// 			return nil, errors.New(aerr.Error())
	// 		}
	// 	} else {
	// 		log.Error("Error: ", err.Error())
	// 		return nil, err
	// 	}
	// }
	// return result, nil
}

func (c *AWSClientSession) GetS3BucketTags(bucket string) (*s3.GetBucketTaggingOutput, error) {
	svc := s3.New(c.session)
	input := &s3.GetBucketTaggingInput{
		Bucket: &bucket,
	}

	result, err := svc.GetBucketTagging(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Error("Error: ", aerr.Error())
				return nil, errors.New(aerr.Error())
			}
		} else {
			log.Error("Error: ", err.Error())
			return nil, err
		}
	}
	return result, nil
}

func (c *AWSClientSession) WhitelistIPInS3Bucket(bucket string, ip string) (bool, error) {
	svc := s3.New(c.session)
	var bucketPolicy *config.S3BucketPolicy
	input := &s3.GetBucketPolicyInput{
		Bucket: &bucket,
	}

	// Getting Bucket Policy
	result, err := svc.GetBucketPolicy(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Error("Error: ", aerr.Error())
				return false, errors.New(aerr.Error())
			}
		} else {
			log.Error("Error: ", err.Error())
			return false, err
		}
	}
	// convert json to map
	json.Unmarshal([]byte(aws.StringValue(result.Policy)), &bucketPolicy)
	for _, statement := range bucketPolicy.Statement {
		for key, value := range statement.Condition {
			if key == "NotIpAddressIfExists" {
				conditionStatement := value.(map[string]interface{})
				if val, ok := conditionStatement["aws:SourceIp"]; ok {
					ipAddress := val.([]interface{})
					new := true
					for _, ips := range ipAddress {
						if ips == ip {
							log.Warn("Your IP is already present there")
							new = false
							return false, nil
						}
					}
					if new {
						ipAddress = append(ipAddress, ip)
						conditionStatement["aws:SourceIp"] = ipAddress
					}
				}
			}
		}
	}
	// convert interface back to json
	policy, _ := json.Marshal(bucketPolicy)
	putInput := &s3.PutBucketPolicyInput{
		Bucket: &bucket,
		Policy: aws.String(string(policy)),
	}
	// updating bucket policy
	_, err = svc.PutBucketPolicy(putInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Error("Error: ", aerr.Error())
				return false, errors.New(aerr.Error())
			}
		} else {
			log.Error("Error: ", err.Error())
			return false, err
		}
	}
	return true, nil
}

func (c *AWSClientSession) RevokeIPFromS3Bucket(bucket string, ip string) error {
	svc := s3.New(c.session)
	var bucketPolicy *config.S3BucketPolicy
	input := &s3.GetBucketPolicyInput{
		Bucket: &bucket,
	}

	// Getting Bucket Policy
	result, err := svc.GetBucketPolicy(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Error("Error: ", aerr.Error())
				return errors.New(aerr.Error())
			}
		} else {
			log.Error("Error: ", err.Error())
			return err
		}
	}
	// convert json to map
	json.Unmarshal([]byte(aws.StringValue(result.Policy)), &bucketPolicy)
	for _, statement := range bucketPolicy.Statement {
		for key, value := range statement.Condition {
			if key == "NotIpAddressIfExists" {
				conditionStatement := value.(map[string]interface{})
				if val, ok := conditionStatement["aws:SourceIp"]; ok {
					ipAddress := val.([]interface{})
					var whitelistedIps []interface{}
					for _, ips := range ipAddress {
						if ips != ip {
							whitelistedIps = append(whitelistedIps, ips)
						}
					}
					conditionStatement["aws:SourceIp"] = whitelistedIps
				}
			}
		}
	}
	// convert interface back to json
	policy, _ := json.Marshal(bucketPolicy)
	putInput := &s3.PutBucketPolicyInput{
		Bucket: &bucket,
		Policy: aws.String(string(policy)),
	}
	// updating bucket policy
	_, err = svc.PutBucketPolicy(putInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Error("Error: ", aerr.Error())
				return errors.New(aerr.Error())
			}
		} else {
			log.Error("Error: ", err.Error())
			return err
		}
	}
	return nil
}

func (c *AWSClientSession) IsConciergeTagTrue(bucketName string) error {
	tags, err := c.GetS3BucketTags(bucketName)
	if err != nil {
		return err
	}

	for _, tag := range tags.TagSet {
		if strings.ToLower(*tag.Key) == "concierge" && strings.ToLower(*tag.Value) == "true" {
			return nil
		}
	}
	return errors.New("tag: concierge is not defined or set to true")
}

// ============================ S3BUCKET END ===================================================
