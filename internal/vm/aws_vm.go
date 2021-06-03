package vm

type AwsVm struct {
	Vm
	connectionParams sshConnectionParams
	InstanceId       string
	KeyPairName      string
}

func (*AwsVm) Terminate() {

}
