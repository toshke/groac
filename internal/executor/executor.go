package executor

import (
	"internal/job"
)

type VmState struct {
	jobExecuted job.Job
	platform    string
}
