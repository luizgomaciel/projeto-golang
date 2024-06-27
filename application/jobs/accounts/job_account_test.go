package jobs_test

import (
	jobs "encoder/application/jobs/accounts"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewJobAccounts(t *testing.T) {
	job, err := jobs.NewJobAccount(2, 10)

	require.NotNil(t, job)
	require.Nil(t, err)
	require.EqualValues(t, len(job.Accounts), 10)
}
