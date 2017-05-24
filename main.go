package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/awslabs/aws-sdk-go/service/ec2"
)

var (
	awsRegion = os.Getenv("AWS_REGION")
)

func tagWorker(nodeName string, nodeLabel string, nodeValue string) {
	args := []string{"label", "nodes", nodeName, nodeLabel + "=" + nodeValue, "--overwrite=true"}
	cmd := exec.Command("kubectl", args...)
	out, err := cmd.Output()
	if err != nil {
		errmsg := err.Error()
		exiterr, ok := err.(*exec.ExitError)
		if ok {
			errmsg = fmt.Sprintf("%s: %s", errmsg, string(exiterr.Stderr))
		}
		log.Printf("Error:%s\n", errmsg)
	}
	log.Printf("Output:%s\n", out)
}

func main() {
	c := ec2metadata.New(session.New())
	hostname, err := c.GetMetadata("hostname")
	if err != nil {
		log.Printf("Unable to retrieve the local hostname from the EC2 instance: %s\n", err)
		return
	}
	instanceID, err := c.GetMetadata("instance-id")
	if err != nil {
		log.Printf("Unable to retrieve the instance ID from the EC2 instance: %s\n", err)
		return
	}
	log.Printf("Hostname is %s\n", hostname)
	log.Printf("Instance ID is %s\n", instanceID)

	sess, err := session.NewSession(&aws.Config{Region: aws.String(awsRegion)})
	svc := ec2.New(sess)

	params := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}
	resp, err := svc.DescribeInstances(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			if instance.InstanceLifecycle != nil {
				var instanceLifecycle = "" + *instance.InstanceLifecycle
				log.Printf("Lifecycle is %s\n", instanceLifecycle)
				tagWorker(hostname, "spot-instance", "true")
			} else {
				tagWorker(hostname, "spot-instance", "false")
			}
		}
	}
}
