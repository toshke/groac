package executor

import (
	"job"
)

type VmState struct {
	jobExecuted job.Job
	platform    string
}
