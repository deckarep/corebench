/*
Open Source Initiative OSI - The MIT License (MIT):Licensing
The MIT License (MIT)
Copyright (c) 2018 Ralph Caraveo (deckarep@gmail.com)
Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package providers

import (
	"os"
	"strings"
	"time"
	log "github.com/sirupsen/logrus"
	"github.com/deckarep/corebench/pkg/ssh"
	"github.com/deckarep/corebench/pkg/utility"
	"context"
	"fmt"
	"bufio"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	// "github.com/aws/aws-sdk-go-v2/service/pricing"
)

type AwsProvider struct {
	client       *ec2.EC2
	cfn          *cloudformation.CloudFormation
  // TODO: add pricing to size cmd
	// pricing      *pricing.Pricing
	instanceType string
	repoLastPath string
	Keyfile      string
	sshKeys      string
}
var (
	privateKey string
	pairName string = "corebench"
	awsRegion string = "us-east-1"
)

func NewAwsProvider() Provider {
	cfg, err := external.LoadDefaultAWSConfig(
		external.WithSharedConfigProfile("default"))
	if err != nil {
			panic("unable to load SDK config, " + err.Error())
	}
  // TODO: force region if sizecmd to get around issue with pricing api composing endpoint out of default region
	cfg.Region = awsRegion

	return &AwsProvider{
		client: ec2.New(cfg),
		cfn: cloudformation.New(cfg),
		// TODO: add pricing to size cmd
		// pricing: pricing.New(cfg),
	}
}

func (p *AwsProvider) List(ctx context.Context) error {
	svc := p.client
	input := &ec2.DescribeInstancesInput{
		Filters: []ec2.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []string{"running", "pending"},
			},
			{
				Name: aws.String("tag-value"),
				Values: []string{"corebench"},
			},
	   },
    }
	req := svc.DescribeInstancesRequest(input)
	result, err := req.Send()
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return nil
	 }
	if len(result.Reservations) == 0 {
		log.Infof("No Instances found. Check the AWS Console!")
		return nil
	} else {
		  for _, reservation := range result.Reservations {
				for _, instance := range reservation.Instances {
        if *instance.PublicDnsName != "" {
					fmt.Printf("corebench instance %v:", *instance.InstanceId)
					fmt.Println("\n Type:             ", instance.InstanceType)
					fmt.Println(" Keypair Name:     ", *instance.KeyName)
					fmt.Println(" PublicDNS:        ", *instance.PublicDnsName)
					fmt.Println(" IP Address:       ", *instance.PublicIpAddress)
				 }
				}
			 }
		 }
	return nil
}

func (p *AwsProvider) SetKeys(keys []string) {
  if p.checkKeypair(keys) != false {
     p.deleteKeypair(keys)
	}
	p.sshKeys = p.createKeypair(keys)
	p.saveKeypair(p.sshKeys)

}

func (p *AwsProvider) Sizes(ctx context.Context) error {

	// TODO: implement pricing api integration, including horrorshow error handling
	// svc := p.pricing
	// 	input := &pricing.GetAttributeValuesInput{
	// 		AttributeName: aws.String("volumeType"),
	// 		MaxResults:    aws.Int64(2),
	// 		ServiceCode:   aws.String("AmazonEC2"),
	// 	}
	//
	// 	req := svc.GetAttributeValuesRequest(input)
	// 	result, err := req.Send()
	// 	if err != nil {
	// 		if aerr, ok := err.(awserr.Error); ok {
	// 			switch aerr.Code() {
	// 			case pricing.ErrCodeInternalErrorException:
	// 				fmt.Println(pricing.ErrCodeInternalErrorException, aerr.Error())
	// 			case pricing.ErrCodeInvalidParameterException:
	// 				fmt.Println(pricing.ErrCodeInvalidParameterException, aerr.Error())
	// 			case pricing.ErrCodeNotFoundException:
	// 				fmt.Println(pricing.ErrCodeNotFoundException, aerr.Error())
	// 			case pricing.ErrCodeInvalidNextTokenException:
	// 				fmt.Println(pricing.ErrCodeInvalidNextTokenException, aerr.Error())
	// 			case pricing.ErrCodeExpiredNextTokenException:
	// 				fmt.Println(pricing.ErrCodeExpiredNextTokenException, aerr.Error())
	// 			default:
	// 				fmt.Println(aerr.Error())
	// 			}
	// 		} else {
	// 			fmt.Println(err.Error())
	// 		}
	// 		return nil
	// 	}

	return nil
}

func (p *AwsProvider) Term(ctx context.Context, settings ProviderTermSettings) error {
	p.cleanup(pairName)
	return nil
}

func (p *AwsProvider) cleanup(keyname string) error {
	svc := p.cfn
	input := &cloudformation.DeleteStackInput{
				StackName:    aws.String(keyname),
			}
	req := svc.DeleteStackRequest(input)
 	_, err := req.Send()
  // the below error check is essentially redundant, aws api will only return error on malformed request
  p.genericAwsErrorCheck(err)
	log.Infof("Cleaning up resources...")
	log.Infof("Stack deletion request sent for \"%v\"", *input.StackName)
	return nil
}

func (p *AwsProvider) createKeypair(keypair []string) string {
	svc := p.client
	input := &ec2.CreateKeyPairInput{
			 KeyName: aws.String(pairName),
	    }
	req := svc.CreateKeyPairRequest(input)
	result, err := req.Send()
	 if err != nil {
			 if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "InvalidKeyPair.Duplicate" {
					 log.Infof("Keypair %q already exists", pairName)
					 return ""
			 }
			 log.Fatalf("Unable to create key pair: %s, %v.", pairName, err)
	 }
	 log.Infof("Created keypair %q", *result.KeyName)
	 privateKey = *result.KeyMaterial
return privateKey
}

func (p *AwsProvider) saveKeypair(privateKey string) error {
	filename := fmt.Sprintf("%s.pem", pairName)
	fileHandle, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create file %s!", filename)
	}
	os.Chmod(filename, 0600)
	writer := bufio.NewWriter(fileHandle)
	defer fileHandle.Close()

  fmt.Fprintln(writer, privateKey)
	writer.Flush()

	log.Infof("Saved new keypair to file: %s", filename)
	p.Keyfile = filename
	return nil
}

func (p *AwsProvider) checkKeypair(keypair []string) bool {
	doesexist := false
	svc := p.client
	input := &ec2.DescribeKeyPairsInput{
			 KeyNames: strings.Split(pairName, ","),
			}
	req := svc.DescribeKeyPairsRequest(input)
	_, err := req.Send()
	if err != nil {
			if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "InvalidKeyPair.NotFound" {
					log.Infof("Key pair %q does not exist", pairName)
					doesexist = false
					return doesexist
			}
	}
	doesexist = true
	log.Infof("Keypair exists! Cleaning up and creating new keypair...")
  return doesexist
}

func (p *AwsProvider) deleteKeypair(keypair []string) error {
	svc := p.client
	input := &ec2.DeleteKeyPairInput{
			 KeyName: aws.String(pairName),
			}
	req := svc.DeleteKeyPairRequest(input)
	_, err := req.Send()
  if err != nil {
      if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "InvalidKeyPair.Duplicate" {
          log.Infof("Key pair %q does not exist, creating..", pairName)
      }
      log.Fatalf("Unable to delete key pair: %s, %v.", pairName, err)
  }
		return nil
  }

func (p *AwsProvider) processCfnTemplate(settings ProviderSpinSettings) string {
	p.repoLastPath = utility.GitPathLast(settings.GitURL())
	finalCfnTemplate :=
		strings.Replace(CfnTemplate, "${go-version}", fmt.Sprintf(goVersionFmt, settings.GoVersion()), -1)
	finalCfnTemplate =
		strings.Replace(finalCfnTemplate, "${git-repo}", settings.GitURL(), -1)
	finalCfnTemplate =
		strings.Replace(finalCfnTemplate, "${git-repo-last-path}", p.repoLastPath, -1)
	finalCfnTemplate =
		strings.Replace(finalCfnTemplate, "${keypair}", pairName, -1)
	finalCfnTemplate =
		strings.Replace(finalCfnTemplate, "${instancetype}", settings.InstanceTypeString(), -1)
	finalCfnTemplate =
		strings.Replace(finalCfnTemplate, "${awsregion}", awsRegion, -1)

	return finalCfnTemplate
}

func (p *AwsProvider) processBenchCommandTemplate(settings ProviderSpinSettings) string {
	AwsBenchCmd :=
		strings.Replace(AwsBenchCommandTemplate, "${git-repo}", settings.GitURL(), -1)
	AwsBenchCmd =
		strings.Replace(AwsBenchCmd, "${cpu-count}", settings.Cpus(), -1)
	AwsBenchCmd =
		strings.Replace(AwsBenchCmd, "${benchmem-setting}", settings.BenchMemString(), -1)
	AwsBenchCmd =
		strings.Replace(AwsBenchCmd, "${bench-regex}", settings.Regex(), -1)
	AwsBenchCmd =
		strings.Replace(AwsBenchCmd, "${bench-count}", fmt.Sprintf("%d", settings.Count()), -1)

	// Should be last, turns on the benchstat summary
	if settings.Stat() {
		AwsBenchCmd = AwsBenchCmd + AwsBenchStatTemplate
	}

	return fmt.Sprintf("%s && %s", AwsBenchReadyScript, AwsBenchCmd)
}

func (p *AwsProvider) genericAwsErrorCheck(err error) error {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Fatalf(aerr.Error())
			}
		} else {
			log.Fatalf(err.Error())
		}
	}
	return nil
}

func (p *AwsProvider) Spinup(ctx context.Context, settings ProviderSpinSettings) error {
	var instanceid string
  svc := p.cfn
	finalCfnTemplate := p.processCfnTemplate(settings)
	input := &cloudformation.CreateStackInput{
				StackName:    aws.String(pairName),
	      TemplateBody: aws.String(finalCfnTemplate),
	    }
	req := svc.CreateStackRequest(input)
	log.Infof("About to provision Cloudformation stack \"%v\" with instance type \"%v\"", *input.StackName, settings.InstanceTypeString())
	if !utility.PromptConfirmation("Continue provisioning? (Yy)es/(Nn)o") {
    log.Info("Cleaning up...")
		p.deleteKeypair(strings.Split(pairName, ","))
		log.Info("Quitting")
		return nil
	}
	result, err := req.Send()
	if p.genericAwsErrorCheck(err) == nil {
		log.Infof("Stack creation request sent: %v\n", result)
	}

 notready := true
 for notready {
  statusinput := &cloudformation.DescribeStacksInput{
					StackName:  aws.String(pairName),
	}
	statusreq := svc.DescribeStacksRequest(statusinput)
	statusresult, err := statusreq.Send()
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return nil
	 } else {
		 for _, stackstatus := range statusresult.Stacks {
       log.Infof("Waiting for stack resource creation to complete...")
			 if stackstatus.StackStatus == "CREATE_COMPLETE" {
			    log.Infof("Good news! Stack status is now %v\n", stackstatus.StackStatus)
					notready = false
		   }	else {
				time.Sleep(30 * time.Second)
			    }
	   }
	   }
 }

	 var chosenIP string
	 const maxDialAttempts = 20

	 // Spin wait - TODO: make this more graceful.
advance_to_ssh:
	for {
	 svc := p.client
	 input := &ec2.DescribeInstancesInput{
	 	Filters: []ec2.Filter{
	 		{
	 			Name: aws.String("instance-state-name"),
	 			Values: []string{"running", "pending", "stopped"},
	 		},
	 		{
	 			Name: aws.String("tag-value"),
	 			Values: []string{"corebench"},
	 		},
	 	 },
	 	}
	 req := svc.DescribeInstancesRequest(input)
	 result, err := req.Send()
	 p.genericAwsErrorCheck(err)

  var ip string
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if len(reservation.Instances) == 0 {
				fmt.Println("")
				log.Infof("No Instances found! Check the AWS Console", err)
				return nil
			}
			ip = *instance.PublicIpAddress
			instanceid = *instance.InstanceId
		}
	}

	if ip != "" {
		chosenIP = ip
		err = ssh.PollSSH(chosenIP + ":22")
		if err == nil {
			break advance_to_ssh
		}
	}
 }
   time.Sleep(time.Second * 30)
	 log.Infof("Instance %v is provisioned and reachable at ip: %v\n", instanceid, chosenIP)
	 log.Info("Instance benchmark starting momentarily...\n")

	 AwsBenchCmd := p.processBenchCommandTemplate(settings)
	 chosenIP = fmt.Sprintf("ubuntu@%s", chosenIP)
	 err = ssh.ExecuteSSH(chosenIP, AwsBenchCmd)
	 if err != nil {
	 	log.Fatalln("Failed to SSH: ", err)
	 }
	 if !settings.LeaveRunning() {
		 defer func () {
		 p.cleanup(pairName)
	   }()
			 } else {
				 log.Infof("Leaving AWS resources running! Execute \"ssh %s -i %s\" to connect to the instance", chosenIP, p.Keyfile)
			 }
 return nil
}
