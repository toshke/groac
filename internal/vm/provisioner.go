package vm

type VmProvisioner interface {
	Provision() *Vm
}
