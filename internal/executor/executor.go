package executor

import (
	. "github.com/toshke/groac/internal/job"
)

type VmState struct {
	jobExecuted Job
	platform    string
}
