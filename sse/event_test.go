package sse

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/json"
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
