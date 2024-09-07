package main

import (
	"context"
	"github.com/anthdm/pricefetcher/proto"
	"google.golang.org/grpc"
)

func makeGRPCServerAndRun(svc PriceFetcher) error {

	grpcPriceFetcher := NewGRPCPriceFetcherServer(svc)

	opts := []grpc.ServerOption()
	server := grpc.NewServer(opts)
}

// based on service.go
type GRPCPriceFetcherServer struct {
	svc PriceFetcher
}

func NewGRPCPriceFetcherServer(svc PriceFetcher) *GRPCPriceFetcherServer {
	return &GRPCPriceFetcherServer{
		svc,
	}
}

func (s *GRPCPriceFetcherServer) FetchPrice(ctx context.Context, req *proto.PriceRequest) (*proto.PriceResponse, error) {
	price, err := s.svc.FetchPrice(ctx, req.Ticker)
	if err != nil {
		return nil, err
	}

	resp := &proto.PriceResponse{
		Ticker: req.Ticker,
		Price:  float32(price),
	}

	return resp, err
}
