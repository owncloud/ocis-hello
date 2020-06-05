// Code generated by protoc-gen-microweb. DO NOT EDIT.
// source: proto.proto

package proto

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/golang/protobuf/jsonpb"
)

type webHelloHandler struct {
	r chi.Router
	h HelloHandler
}

func (h *webHelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

func (h *webHelloHandler) Greet(w http.ResponseWriter, r *http.Request) {

	req := &GreetRequest{}

	resp := &GreetResponse{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusPreconditionFailed)
		return
	}

	if err := h.h.Greet(
		r.Context(),
		req,
		resp,
	); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, resp)
}

func RegisterHelloWeb(r chi.Router, i HelloHandler, middlewares ...func(http.Handler) http.Handler) {
	handler := &webHelloHandler{
		r: r,
		h: i,
	}

	r.MethodFunc("POST", "/api/v0/hello/greet", handler.Greet)
}

// GreetRequestJSONMarshaler describes the default jsonpb.Marshaler used by all
// instances of GreetRequest. This struct is safe to replace or modify but
// should not be done so concurrently.
var GreetRequestJSONMarshaler = new(jsonpb.Marshaler)

// MarshalJSON satisfies the encoding/json Marshaler interface. This method
// uses the more correct jsonpb package to correctly marshal the message.
func (m *GreetRequest) MarshalJSON() ([]byte, error) {
	if m == nil {
		return json.Marshal(nil)
	}

	buf := &bytes.Buffer{}

	if err := GreetRequestJSONMarshaler.Marshal(buf, m); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var _ json.Marshaler = (*GreetRequest)(nil)

// GreetRequestJSONUnmarshaler describes the default jsonpb.Unmarshaler used by all
// instances of GreetRequest. This struct is safe to replace or modify but
// should not be done so concurrently.
var GreetRequestJSONUnmarshaler = new(jsonpb.Unmarshaler)

// UnmarshalJSON satisfies the encoding/json Unmarshaler interface. This method
// uses the more correct jsonpb package to correctly unmarshal the message.
func (m *GreetRequest) UnmarshalJSON(b []byte) error {
	return GreetRequestJSONUnmarshaler.Unmarshal(bytes.NewReader(b), m)
}

var _ json.Unmarshaler = (*GreetRequest)(nil)

// GreetResponseJSONMarshaler describes the default jsonpb.Marshaler used by all
// instances of GreetResponse. This struct is safe to replace or modify but
// should not be done so concurrently.
var GreetResponseJSONMarshaler = new(jsonpb.Marshaler)

// MarshalJSON satisfies the encoding/json Marshaler interface. This method
// uses the more correct jsonpb package to correctly marshal the message.
func (m *GreetResponse) MarshalJSON() ([]byte, error) {
	if m == nil {
		return json.Marshal(nil)
	}

	buf := &bytes.Buffer{}

	if err := GreetResponseJSONMarshaler.Marshal(buf, m); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var _ json.Marshaler = (*GreetResponse)(nil)

// GreetResponseJSONUnmarshaler describes the default jsonpb.Unmarshaler used by all
// instances of GreetResponse. This struct is safe to replace or modify but
// should not be done so concurrently.
var GreetResponseJSONUnmarshaler = new(jsonpb.Unmarshaler)

// UnmarshalJSON satisfies the encoding/json Unmarshaler interface. This method
// uses the more correct jsonpb package to correctly unmarshal the message.
func (m *GreetResponse) UnmarshalJSON(b []byte) error {
	return GreetResponseJSONUnmarshaler.Unmarshal(bytes.NewReader(b), m)
}

var _ json.Unmarshaler = (*GreetResponse)(nil)
