package vm

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
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

func (provisioner *VmProvisionerAws) Provision(groacConfig map[string]string) *AwsVm {
	params := &ec2.RunInstancesInput{}

	params.MaxCount = aws.Int32(1)
	params.MinCount = aws.Int32(1)

	// generate key if does not exist
	connectionParams := newSshConnectionParams()
	connectionParams.initPrivateKey()

	//TODO: parametrise keypair name, allow configuration - provided value

	// create key if does not exist
	readKeyPairParams := &ec2.DescribeKeyPairsInput{KeyNames: []string{"groac-gitlab-runner-key"}}
	readKeyOutput, err := provisioner.ec2Client.DescribeKeyPairs(context.TODO(), readKeyPairParams)
	publicKeyMaterial, err := publicKeyMaterial(connectionParams.publicKeyPath)
	var keyPairName *string
	if val, ok := groacConfig["awsKeyPairName"]; ok {
		keyPairName = aws.String(val)
	} else {
		keyPairName = aws.String("groac-gitlab-runner-key")
	}
	// if the key exists with different finger print groac-gitlab-runner-key-timestamp key will be created
	if err != nil {
		// check if error claims that key does not exist
		if strings.Contains(err.Error(), "InvalidKeyPair.NotFound") {
			importKeyPairParams := &ec2.ImportKeyPairInput{
				KeyName:           keyPairName,
				PublicKeyMaterial: []byte(base64.StdEncoding.EncodeToString(publicKeyMaterial)),
			}
			importResponse, err := provisioner.ec2Client.ImportKeyPair(context.TODO(), importKeyPairParams)
			check(err)
			log.Info().Str("ProvisionedGroacKeyId", *importResponse.KeyPairId)
			log.Info().Str("ProvisionedGroacKeyName", *keyPairName)
		} else {
			check(err)
		}
	} else {
		// verify that keypair fingerprent matches local one
		awsFingerprint := *readKeyOutput.KeyPairs[0].KeyFingerprint
		localFingerprint := publicKeyFingerPrint(publicKeyMaterial)
		if awsFingerprint != localFingerprint {
			keyPairName = aws.String("groac-gitlab-runner-key-" + fmt.Sprintf("%d", time.Now().Unix()))
			importKeyPairParams := &ec2.ImportKeyPairInput{
				KeyName:           keyPairName,
				PublicKeyMaterial: []byte(base64.StdEncoding.EncodeToString(publicKeyMaterial)),
			}
			importResponse, err := provisioner.ec2Client.ImportKeyPair(context.TODO(), importKeyPairParams)
			check(err)
			log.Info().Str("ProvisionedGroacKeyId", *importResponse.KeyPairId)
			log.Info().Str("ProvisionedGroacKeyName", *keyPairName)
		}
	}

	params.KeyName = keyPairName

	//TODO: create Security Group if none given
	//TODO: discover default VPC if none given
	//TODO: pick default subnet if none given
	output, err := provisioner.ec2Client.RunInstances(context.TODO(), params)
	check(err)

	var awsVm AwsVm
	awsVm.connectionParams.hostname = *output.Instances[0].PrivateIpAddress

	return &awsVm
}
