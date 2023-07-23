package amazon

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"main/pkg/configs"
)

type EC2DataTypes struct {
	instanceTypes []types.InstanceTypeInfo
}
type EC2Images struct {
	ami []types.Image
}

func GetEC2InstanceTypes() []types.InstanceTypeInfo {
	client := configs.GetEC2Client()
	ec2InstanceTypes := &EC2DataTypes{}
	pagingInstanceTypes(client, nil, ec2InstanceTypes)
	return ec2InstanceTypes.instanceTypes
}

func pagingInstanceTypes(client *ec2.Client, nextToken *string, ec2InstanceTypes *EC2DataTypes) {
	resp, err := client.DescribeInstanceTypes(context.TODO(), &ec2.DescribeInstanceTypesInput{
		NextToken: nextToken,
	})
	if err != nil {
		return
	}
	if resp.NextToken != nil {
		pagingInstanceTypes(client, resp.NextToken, ec2InstanceTypes)
	}
	ec2InstanceTypes.instanceTypes = append(ec2InstanceTypes.instanceTypes, resp.InstanceTypes...)
}

func GetEC2AMI() []types.Image {
	client := configs.GetEC2Client()
	ec2Images := &EC2Images{}
	pagingAMI(client, nil, ec2Images)
	return ec2Images.ami
}

func pagingAMI(client *ec2.Client, nextToken *string, ec2Images *EC2Images) {
	filterType := "owner-alias"
	publicFilterType := "is-public"
	ownerName := "amazon"
	resp, err := client.DescribeImages(context.TODO(), &ec2.DescribeImagesInput{
		NextToken: nextToken,
		Filters: []types.Filter{
			{
				Name:   &filterType,
				Values: []string{ownerName},
			},
			{
				Name:   &publicFilterType,
				Values: []string{"true"},
			},
		},
	})
	if err != nil {
		return
	}
	ec2Images.ami = append(ec2Images.ami, resp.Images...)
}
