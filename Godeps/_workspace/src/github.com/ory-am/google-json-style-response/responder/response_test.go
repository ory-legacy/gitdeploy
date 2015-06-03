package responder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorResponseMarshal(t *testing.T) {
	e := ErrorResponse{
		ApiVersion: "1",
	}
	r, err := e.Marshal()
	assert.Nil(t, err)
	assert.Equal(t, r, []byte("{\"apiVersion\":\"1\",\"error\":{}}"))
}

func TestDataResponseMarshal(t *testing.T) {
	d := DataResponse{
		ApiVersion: "1",
	}
	r, e := d.Marshal()
	assert.Nil(t, e)
	t.Log(string(r))
	assert.Equal(t, r, []byte("{\"apiVersion\":\"1\"}"))
}
