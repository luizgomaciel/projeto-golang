package jobs

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestAccountValidationIsComplete(t *testing.T) {
	account := NewAccount(2)

	account.ID = uuid.NewV4().String()
	err := account.Validate()

	require.Nil(t, err)
}
