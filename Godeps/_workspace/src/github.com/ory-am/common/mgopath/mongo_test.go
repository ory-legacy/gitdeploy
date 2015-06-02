package mgopath

import (
	"fmt"
	"github.com/ory-am/common/rand/sequence"
	"github.com/ory-am/dockertest"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConnect(t *testing.T) {
	containerID, ip, port := dockertest.SetupMongoContainer(t)
	dbName, err := sequence.RuneSequence(22, []rune("abcdefghijklmnopqrstuvxyzABCDEFGHIJKLMNOPQRSTUVXYZ1234567890"))
	assert.Nil(t, err)
	defer containerID.KillRemove(t)
	path := fmt.Sprintf("mongodb://%s:%d/%s", ip, port, string(dbName))
	db, name, err := Connect(path)
	assert.NotNil(t, db)
	assert.Equal(t, name, string(dbName))
	assert.Nil(t, err)
}
