package trace

import (
	"net/http"

	"google.golang.org/grpc/metadata"
)

const (
	HttpFormat = iota
	GrpcFormat
)

var (
	emptyHttpPropagator httpPropagator
	emptyGrpcPropagator grpcPropagator
)

type (
	Propagator interface {
		Extract(carrier interface{}) (Carrier, error)
		Inject(carrier interface{}) (Carrier, error)
	}

	httpPropagator struct{}
	grpcPropagator struct{}
)

func (h httpPropagator) Extract(carrier interface{}) (Carrier, error) {
	if c, ok := carrier.(http.Header); !ok {
		return nil, ErrInvalidCarrier
	} else {
		return httpCarrier(c), nil
	}
}

func (h httpPropagator) Inject(carrier interface{}) (Carrier, error) {
	if c, ok := carrier.(http.Header); ok {
		return httpCarrier(c), nil
	} else {
		return nil, ErrInvalidCarrier
	}
}

func (g grpcPropagator) Extract(carrier interface{}) (Carrier, error) {
	if c, ok := carrier.(metadata.MD); ok {
		return grpcCarrier(c), nil
	} else {
		return nil, ErrInvalidCarrier
	}
}

func (g grpcPropagator) Inject(carrier interface{}) (Carrier, error) {
	if c, ok := carrier.(metadata.MD); ok {
		return grpcCarrier(c), nil
	} else {
		return nil, ErrInvalidCarrier
	}
}
