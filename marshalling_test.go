package gotils

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func GimmeString(i interface{}) (interface{}, error) {
	switch i.(type) {
	case nil:
		return nil, nil
	default:
		return fmt.Sprint(i), nil
	}
}
func GimmeTime(i interface{}) (interface{}, error) {
	switch i.(type) {
	case int:
		return time.Unix(0, int64(i.(int))*1000000).UTC(), nil
	case int64:
		return time.Unix(0, i.(int64)*1000000).UTC(), nil
	case string:
		format := "2006-01-02"
		return time.Parse(format, i.(string))
	}
	return time.Time{}, fmt.Errorf("could not make time")
}

type MixedType struct {
	Phone     string    `json:"phone,omitempty"`
	Country   string    `json:"country"`
	Amount    float64   `json:"amount,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func TestMarshalWithCastsCastingMultipleFuncs(t *testing.T) {
	castFns := map[string]func(interface{}) (interface{}, error){
		"Phone":     GimmeString,
		"Timestamp": GimmeTime,
	}

	j := []byte(`{"phone":50,"country":"ES","Timestamp":"2020-01-01","amount":200}`)

	newObj, _ := MarshalWithCasts(j, MixedType{}, castFns)
	mt := newObj.(MixedType)
	b, _ := json.Marshal(mt)

	assert.Equal(t, `{"phone":"50","country":"ES","amount":200,"timestamp":"2020-01-01T00:00:00Z"}`, string(b))
}

func TestMarshalWithZeroValues(t *testing.T) {
	castFns := map[string]func(interface{}) (interface{}, error){
		"Phone":     GimmeString,
		"Timestamp": GimmeTime,
	}

	j := []byte(`{"country":"ES","Timestamp":"2020-01-01","amount":200}`)

	newObj, _ := MarshalWithCasts(j, MixedType{}, castFns)
	mt := newObj.(MixedType)
	b, _ := json.Marshal(mt)

	assert.Equal(t, `{"country":"ES","amount":200,"timestamp":"2020-01-01T00:00:00Z"}`, string(b))
}

type SecretType struct {
	Phone     string `json:"phone,omitempty"`
	country   string
	Amount    float64   `json:"amount,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func TestMarshalWithUnexportedFields(t *testing.T) {
	castFns := map[string]func(interface{}) (interface{}, error){
		"Phone":     GimmeString,
		"Timestamp": GimmeTime,
	}

	j := []byte(`{"Timestamp":"2020-01-01","amount":200}`)

	newObj, _ := MarshalWithCasts(j, SecretType{}, castFns)
	mt := newObj.(SecretType)
	b, _ := json.Marshal(mt)

	assert.Equal(t, `{"amount":200,"timestamp":"2020-01-01T00:00:00Z"}`, string(b))
}

func TestMarshalWithErrorsInCastingFn(t *testing.T) {
	castFns := map[string]func(interface{}) (interface{}, error){
		"Phone":     GimmeString,
		"Timestamp": GimmeTime,
	}

	j := []byte(`{"phone":50,"country":"ES","Timestamp":"2020-01-018373","amount":200}`)

	_, err := MarshalWithCasts(j, MixedType{}, castFns)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "parsing time")
}
