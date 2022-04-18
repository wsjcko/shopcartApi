// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: shopcartApi.proto

package go_micro_api_shop_cart

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

// Api Endpoints for ShopCartApi service

func NewShopCartApiEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for ShopCartApi service

type ShopCartApiService interface {
	FindAll(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
}

type shopCartApiService struct {
	c    client.Client
	name string
}

func NewShopCartApiService(name string, c client.Client) ShopCartApiService {
	return &shopCartApiService{
		c:    c,
		name: name,
	}
}

func (c *shopCartApiService) FindAll(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "ShopCartApi.FindAll", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ShopCartApi service

type ShopCartApiHandler interface {
	FindAll(context.Context, *Request, *Response) error
}

func RegisterShopCartApiHandler(s server.Server, hdlr ShopCartApiHandler, opts ...server.HandlerOption) error {
	type shopCartApi interface {
		FindAll(ctx context.Context, in *Request, out *Response) error
	}
	type ShopCartApi struct {
		shopCartApi
	}
	h := &shopCartApiHandler{hdlr}
	return s.Handle(s.NewHandler(&ShopCartApi{h}, opts...))
}

type shopCartApiHandler struct {
	ShopCartApiHandler
}

func (h *shopCartApiHandler) FindAll(ctx context.Context, in *Request, out *Response) error {
	return h.ShopCartApiHandler.FindAll(ctx, in, out)
}
