package responder

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	r := New("1.0")
	assert.NotNil(t, r)
}

type DataMock struct {
	A string
	B string
}

func TestWriteData(t *testing.T) {
	responder := New("1.0")
	data := DataMock{
		A: "a",
		B: "b",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o := responder.Success(data)
		err := responder.Write(w, o)
		assert.Nil(t, err)
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.Nil(t, err)
	httpResponse, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	res.Body.Close()

	dataResponse := new(DataResponse)
	assert.Nil(t, json.Unmarshal(httpResponse, dataResponse))
	dm := dataResponse.Data.(map[string]interface{})
	assert.Equal(t, data.A, dm["A"])
}

func TestWriteError(t *testing.T) {
	errorCode := 500
	errorMessage := "foobar"
	responder := New("1.0")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o := responder.Error(errorCode, errorMessage)
		responder.AddError(ErrorItem{
			Message: "baz",
		})
		err := responder.Write(w, o)
		assert.Nil(t, err)
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.Nil(t, err)
	httpResponse, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	res.Body.Close()

	errorResponse := new(ErrorResponse)
	assert.Nil(t, json.Unmarshal(httpResponse, errorResponse))
	assert.Equal(t, errorCode, errorResponse.Error.Code)
}
