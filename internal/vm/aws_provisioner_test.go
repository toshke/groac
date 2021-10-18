package vm

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssm_types "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/go-playground/assert/v2"
)

type MockSSMClient struct {
	ssmiface.SSMAPI
	getParamOut ssm.GetParameterOutput
}

type MockEC2Client struct {
	ec2iface.EC2API
}

func (m MockSSMClient) GetParameter(c context.Context, in *ssm.GetParameterInput, opts ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
	return &m.getParamOut, nil
}

func (m MockEC2Client) RunInstances(context.Context, *ec2.RunInstancesInput, ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error) {
	return nil, nil
}
func (m MockEC2Client) DescribeKeyPairs(context.Context, *ec2.DescribeKeyPairsInput, ...func(*ec2.Options)) (*ec2.DescribeKeyPairsOutput, error) {
	return nil, nil
}
func (m MockEC2Client) ImportKeyPair(context.Context, *ec2.ImportKeyPairInput, ...func(*ec2.Options)) (*ec2.ImportKeyPairOutput, error) {
	return nil, nil
}
func (m MockEC2Client) DescribeSubnets(context.Context, *ec2.DescribeSubnetsInput, ...func(*ec2.Options)) (*ec2.DescribeSubnetsOutput, error) {
	return nil, nil
}
func (m MockEC2Client) DescribeAvailabilityZones(context.Context, *ec2.DescribeAvailabilityZonesInput, ...func(*ec2.Options)) (*ec2.DescribeAvailabilityZonesOutput, error) {
	return nil, nil
}
func (m MockEC2Client) DescribeVpcs(context.Context, *ec2.DescribeVpcsInput, ...func(*ec2.Options)) (*ec2.DescribeVpcsOutput, error) {
	return nil, nil
}
func (m MockEC2Client) CreateSecurityGroup(context.Context, *ec2.CreateSecurityGroupInput, ...func(*ec2.Options)) (*ec2.CreateSecurityGroupOutput, error) {
	return nil, nil
}
func (m MockEC2Client) AuthorizeSecurityGroupEgress(context.Context, *ec2.AuthorizeSecurityGroupEgressInput, ...func(*ec2.Options)) (*ec2.AuthorizeSecurityGroupEgressOutput, error) {
	return nil, nil
}
func (m MockEC2Client) AuthorizeSecurityGroupIngress(context.Context, *ec2.AuthorizeSecurityGroupIngressInput, ...func(*ec2.Options)) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
	return nil, nil
}
func TestDiscoverAmiId(t *testing.T) {
	t.Run("verify explicit configuration", func(t *testing.T) {
		configuration := make(map[string]string)
		configuration["AmiId"] = "ami-12345"
		amiId, err := discoverAmiId(nil, configuration)
		assert.Equal(t, *amiId, "ami-12345")
		assert.IsEqual(err, nil)
	})

	t.Run("verify discovered configuration", func(t *testing.T) {
		configuration := make(map[string]string)
		client := new(MockSSMClient)
		client.getParamOut = ssm.GetParameterOutput{
			Parameter: &ssm_types.Parameter{
				Value: aws.String("ami-12345"),
			},
		}
		amiId, err := discoverAmiId(client, configuration)
		assert.Equal(t, *amiId, "ami-12345")
		assert.IsEqual(err, nil)
	})
}

func TestDiscoverSecurityGroupId(t *testing.T) {
	t.Run("verify explicit configuration", func(t *testing.T) {
		configuration := make(map[string]string)
		client := new(MockEC2Client)
		configuration["SecurityGroupIds"] = "sg-123,sg-345"
		vpcId := "vpc-123"
		sgIds, err := discoverSecurityGroupIds(client, &vpcId, configuration)
		assert.Equal(t, []string{"sg-123", "sg-345"}, sgIds)
		assert.IsEqual(err, nil)
	})

}
