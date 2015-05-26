package sse

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetApp(t *testing.T) {
	e := NewEvent("foo", "bar")
	assert.Equal(t, "foo", e.GetApp())
}

func TestSSEify(t *testing.T) {
	e := NewEvent("foo", "bar")
	r, err := json.MarshalIndent(e, "data: ", "\t")
	assert.Nil(t, err)
	assert.Equal(t, string(r), e.SSEify())
}
