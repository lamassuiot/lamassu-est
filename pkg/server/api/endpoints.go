package api

import (
	"context"
	"crypto/x509"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/opentracing"
	stdopentracing "github.com/opentracing/opentracing-go"
)

type Endpoints struct {
	HealthEndpoint       endpoint.Endpoint
	GetCAsEndpoint       endpoint.Endpoint
	EnrollerEndpoint     endpoint.Endpoint
	ReenrollerEndpoint   endpoint.Endpoint
	ServerKeyGenEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service, otTracer stdopentracing.Tracer) Endpoints {
	var healthEndpoint endpoint.Endpoint
	{
		healthEndpoint = MakeHealthEndpoint(s)
		healthEndpoint = opentracing.TraceServer(otTracer, "Health")(healthEndpoint)
	}

	var getCasEndpoint endpoint.Endpoint
	{
		getCasEndpoint = MakeGetCAsEndpoint(s)
		getCasEndpoint = opentracing.TraceServer(otTracer, "GetCAs")(getCasEndpoint)
	}

	var enrollEndpoint endpoint.Endpoint
	{
		enrollEndpoint = MakeEnrollEndpoint(s)
		enrollEndpoint = opentracing.TraceServer(otTracer, "Enroll")(enrollEndpoint)
	}

	var reenrollEndpoint endpoint.Endpoint
	{
		reenrollEndpoint = MakeReenrollEndpoint(s)
		reenrollEndpoint = opentracing.TraceServer(otTracer, "Reenroll")(reenrollEndpoint)
	}
	var serverkeygenEndpoint endpoint.Endpoint
	{
		serverkeygenEndpoint = MakeServerKeyGenEndpoint(s)
		serverkeygenEndpoint = opentracing.TraceServer(otTracer, "Serverkeygen")(serverkeygenEndpoint)
	}
	return Endpoints{
		HealthEndpoint:       healthEndpoint,
		GetCAsEndpoint:       getCasEndpoint,
		EnrollerEndpoint:     enrollEndpoint,
		ReenrollerEndpoint:   reenrollEndpoint,
		ServerKeyGenEndpoint: serverkeygenEndpoint,
	}
}

func MakeHealthEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		healthy := s.Health(ctx)
		return HealthResponse{Healthy: healthy}, nil
	}
}

func MakeGetCAsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		cas, err := s.CACerts(ctx, "", nil)
		return GetCasResponse{Certs: cas}, err
	}
}

func MakeEnrollEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(EnrollRequest)
		cas, err := s.Enroll(ctx, req.csr, req.aps, req.crt, nil)
		return EnrollReenrollResponse{Cert: cas}, err
	}
}

func MakeReenrollEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ReenrollRequest)
		cas, err := s.Reenroll(ctx, req.crt, req.csr, "", nil)
		return EnrollReenrollResponse{Cert: cas}, err
	}
}

func MakeServerKeyGenEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ServerKeyGenRequest)
		cas, key, err := s.ServerKeyGen(ctx, req.csr, req.aps, nil)
		return ServerKeyGenResponse{Cert: cas, Key: key}, err
	}
}

type EmptyRequest struct{}

type EnrollRequest struct {
	csr *x509.CertificateRequest
	aps string
	crt *x509.Certificate
}

type ReenrollRequest struct {
	csr *x509.CertificateRequest
	crt *x509.Certificate
}
type ServerKeyGenRequest struct {
	csr *x509.CertificateRequest
	aps string
}

type HealthResponse struct {
	Healthy bool  `json:"healthy,omitempty"`
	Err     error `json:"err,omitempty"`
}

type GetCasResponse struct {
	Certs []*x509.Certificate
}
type EnrollReenrollResponse struct {
	Cert *x509.Certificate
}
type ServerKeyGenResponse struct {
	Cert *x509.Certificate
	Key  []byte
}
