// File: helpers.go
package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/golang/gddo/httputil/header"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var protoMarshaler = protojson.MarshalOptions{
	AllowPartial:    true,
	UseProtoNames:   true,
	EmitUnpopulated: true,
}

var protoUnmarshaler = protojson.UnmarshalOptions{
	AllowPartial:   true,
	DiscardUnknown: true,
}

func MarshalToJson(message proto.Message) ([]byte, error) {
	return protoMarshaler.Marshal(message)
}

func DecodeJsonToPbMessage(jsonBytes []byte, dst interface{}) *Status {
	err := json.Unmarshal(jsonBytes, dst)
	if err != nil {
		status := NewInternalServiceErrorStatus(err.Error(), "REST_BODY_DECODE_JSON")
		return &status
	}
	return nil
}

func DecodeJsonBodyAsObject(w http.ResponseWriter, r *http.Request, sizeLimit int64, dst interface{}) *Status {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			status := NewBadRequestStatus(msg)
			return &status
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, sizeLimit)

	jsonBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		status := NewBadRequestStatus(err.Error())
		return &status
	}
	return DecodeJsonToPbMessage(jsonBytes, dst)
}

func DecodeJsonBodyAsArray(w http.ResponseWriter, r *http.Request, sizeLimit int64, nextElement func() interface{}) *Status {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			status := NewBadRequestStatus(msg)
			return &status
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, sizeLimit)
	var jsonArray []interface{}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&jsonArray)
	if err != nil {
		status := NewBadRequestStatus(err.Error())
		return &status
	}

	for _, obj := range jsonArray {
		bytes, err := json.Marshal(obj)
		if err != nil {
			// this should not happen, but just for completeness
			status := NewBadRequestStatus(err.Error())
			return &status
		}
		status := DecodeJsonToPbMessage(bytes, nextElement())
		if status != nil {
			return status
		}
	}
	return nil
}
