package sse

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetApp(t *testing.T) {
	e := NewEvent("foo", "bar", "baz")
	assert.Equal(t, "foo", e.App)
}

func TestSSEify(t *testing.T) {
	e := NewEvent("foo", "bar", "baz")
	r, err := json.MarshalIndent(e, "data: ", "\t")
	assert.Nil(t, err)
	assert.Equal(t, string(r), e.SSEify())
}
