package jobs

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestValidateIfAccountIsError(t *testing.T) {
	account := NewAccount(2)

	err := account.Validate()

	require.Error(t, err)
}

func TestAccountValidationIsComplete(t *testing.T) {
	account := NewAccount(2)

	account.ID = uuid.NewV4().String()
	err := account.Validate()

	require.Nil(t, err)
}
