package vm

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/rs/zerolog/log"
)

type VmProvisionerAws struct {
	securityGroupId string
	keyPairName     string
	instanceType    string
	amiId           string
	vpcId           string
	subnetIds       []string
	subnetPointer   int
	ec2Client       *ec2.Client
}

func NewVmProvisionerAws() *VmProvisionerAws {
	var provisioner VmProvisionerAws
	cfg, err := config.LoadDefaultConfig(context.TODO())
	check(err)

	client := ec2.NewFromConfig(cfg)
	provisioner.ec2Client = client

	return &provisioner
}

func (provisioner *VmProvisionerAws) Provision(groacConfig map[string]string, sequence int) *AwsVm {
	params := &ec2.RunInstancesInput{}

	params.MaxCount = aws.Int32(1)
	params.MinCount = aws.Int32(1)

	// generate key if does not exist
	connectionParams := newSshConnectionParams()
	connectionParams.initPrivateKey()

	//TODO: parametrise keypair name, allow configuration - provided value

	// create key if does not exist
	keyPairName, err := discoverKeyPair(provisioner.ec2Client, groacConfig, connectionParams.publicKeyPath)
	check(err)
	params.KeyName = keyPairName

	// subnet discovery
	subnetId, err := discoverSubnetId(provisioner.ec2Client, groacConfig, sequence)
	check(err)
	params.SubnetId = subnetId

	// security group discovery
	sgIds, err := discoverSecurityGroupIds(provisioner.ec2Client, groacConfig)
	check(err)
	params.SecurityGroupIds = sgIds
	//TODO: create Security Group if none given, allowing for SSH and Docker connection

	//TODO: configure context

	runInstanceOut, err := provisioner.ec2Client.RunInstances(context.TODO(), params)
	check(err)

	var awsVm AwsVm
	awsVm.connectionParams.hostname = *runInstanceOut.Instances[0].PrivateIpAddress

	return &awsVm
}

// return keypair used for groac. Discover based on default key pair name, and verify
// fingerprints. If fingerprint does not match, created new key pair for given name with
// timestamp suffix
func discoverKeyPair(client *ec2.Client, groacConfig map[string]string, publicKeyPath string) (*string, error) {
	var keyPairName *string
	if val, ok := groacConfig["awsKeyPairName"]; ok {
		keyPairName = aws.String(val)
	} else {
		keyPairName = aws.String("groac-gitlab-runner-key")
	}
	readKeyPairParams := &ec2.DescribeKeyPairsInput{KeyNames: []string{*keyPairName}}
	readKeyOutput, err := client.DescribeKeyPairs(context.TODO(), readKeyPairParams)
	if err != nil {
		return nil, err
	}

	publicKeyMaterial, err := publicKeyMaterial(publicKeyPath)
	if err != nil {
		return nil, err
	}

	// if the key exists with different finger print groac-gitlab-runner-key-timestamp key will be created
	if err != nil {
		// check if error claims that key does not exist
		if strings.Contains(err.Error(), "InvalidKeyPair.NotFound") {
			importKeyPairParams := &ec2.ImportKeyPairInput{
				KeyName:           keyPairName,
				PublicKeyMaterial: []byte(base64.StdEncoding.EncodeToString(publicKeyMaterial)),
			}
			importResponse, err := client.ImportKeyPair(context.TODO(), importKeyPairParams)
			if err != nil {
				return nil, err
			}
			log.Info().Str("ProvisionedGroacKeyId", *importResponse.KeyPairId)
			log.Info().Str("ProvisionedGroacKeyName", *keyPairName)
		} else {
			return nil, err
		}
	} else {
		// verify that keypair fingerprent matches local one
		awsFingerprint := *readKeyOutput.KeyPairs[0].KeyFingerprint
		localFingerprint := publicKeyFingerPrint(publicKeyMaterial)
		if awsFingerprint != localFingerprint {
			keyPairName = aws.String(*keyPairName + "-" + fmt.Sprintf("%d", time.Now().Unix()))
			importKeyPairParams := &ec2.ImportKeyPairInput{
				KeyName:           keyPairName,
				PublicKeyMaterial: []byte(base64.StdEncoding.EncodeToString(publicKeyMaterial)),
			}
			importResponse, err := client.ImportKeyPair(context.TODO(), importKeyPairParams)
			check(err)
			log.Info().Str("ProvisionedGroacKeyId", *importResponse.KeyPairId)
			log.Info().Str("ProvisionedGroacKeyName", *keyPairName)
		}
	}
	return keyPairName, nil
}

// pull out subnet from configuration, all discover default subnet. Azs are being
// picked by round robin based on sequence of the machine being provisioned
func discoverSubnetId(client *ec2.Client, groacConfig map[string]string, sequence int) (*string, error) {
	if val, ok := groacConfig["SubnetIds"]; ok {
		allValues := strings.Split(val, ",")
		return &allValues[sequence%len(allValues)], nil
	} else {
		// discover subnet
		// list all zones
		azs, _ := client.DescribeAvailabilityZones(context.TODO(), &ec2.DescribeAvailabilityZonesInput{})
		allValues := azs.AvailabilityZones
		selectAz := allValues[sequence%len(allValues)]
		log.Info().Str("selectAZ", *selectAz.ZoneName)
		vpcResp, _ := client.DescribeVpcs(context.TODO(),
			&ec2.DescribeVpcsInput{Filters: []types.Filter{
				{Name: aws.String("isDefault"), Values: []string{"true"}},
			}})
		if len(vpcResp.Vpcs) < 1 {
			return nil, errors.New("subnets not suplied via SubnetIds parameter, and no default VPC found ")
		}
		vpcId := vpcResp.Vpcs[0].VpcId
		log.Info().Str("defaultVpcId", *vpcId)
		subnetResponse, _ := client.DescribeSubnets(context.TODO(), &ec2.DescribeSubnetsInput{
			Filters: []types.Filter{
				{Name: aws.String("vpc-id"), Values: []string{*vpcId}},
				{Name: aws.String("availability-zone-id"), Values: []string{*selectAz.ZoneId}},
				{Name: aws.String("default-for-az"), Values: []string{"true"}},
			}})
		if len(subnetResponse.Subnets) < 1 {
			return nil, fmt.Errorf(
				"subnets not suplied via SubnetIds parameter, and no default subnet for AZ %s",
				*selectAz.ZoneName,
			)
		}
		return subnetResponse.Subnets[0].SubnetId, nil
	}
}

func discoverSecurityGroupIds(client *ec2.Client, groacConfig map[string]string) ([]string, error) {
	return nil, nil
}
