// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: hello.proto

package proto

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "go-micro.dev/v4/api"
	client "go-micro.dev/v4/client"
	server "go-micro.dev/v4/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Hello service

func NewHelloEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Hello service

type HelloService interface {
	Greet(ctx context.Context, in *GreetRequest, opts ...client.CallOption) (*GreetResponse, error)
}

type helloService struct {
	c    client.Client
	name string
}

func NewHelloService(name string, c client.Client) HelloService {
	return &helloService{
		c:    c,
		name: name,
	}
}

func (c *helloService) Greet(ctx context.Context, in *GreetRequest, opts ...client.CallOption) (*GreetResponse, error) {
	req := c.c.NewRequest(c.name, "Hello.Greet", in)
	out := new(GreetResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Hello service

type HelloHandler interface {
	Greet(context.Context, *GreetRequest, *GreetResponse) error
}

func RegisterHelloHandler(s server.Server, hdlr HelloHandler, opts ...server.HandlerOption) error {
	type hello interface {
		Greet(ctx context.Context, in *GreetRequest, out *GreetResponse) error
	}
	type Hello struct {
		hello
	}
	h := &helloHandler{hdlr}
	return s.Handle(s.NewHandler(&Hello{h}, opts...))
}

type helloHandler struct {
	HelloHandler
}

func (h *helloHandler) Greet(ctx context.Context, in *GreetRequest, out *GreetResponse) error {
	return h.HelloHandler.Greet(ctx, in, out)
}
