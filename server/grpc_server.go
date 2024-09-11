package grpc_server

import (
	"context"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/kartikaysaxena/go-microservice/proto"
	"github.com/kartikaysaxena/go-microservice/redis"
	"google.golang.org/grpc"
)

var (
	RateLimit  = 10
	TimeWindow = 1 * time.Minute
)

type ModifyRate interface {
	RateModifier(ctx context.Context, rateRequest *proto.RateRequest) (*proto.RateResponse, error)
}

func NewRateModifier(svc ModifyRate) ModifyRate {
	return &ModifyRateImplementation{}
}


type GRPCRateModifier struct {
	svc ModifyRate
	proto.UnimplementedRateModifierServer
}

func MakeGRPCServerAndRun(svc ModifyRate, listenAddr string) error {
	grpcServer := NewGRPCRateModifier(svc)
	
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err 
	}
	server := grpc.NewServer()
	proto.RegisterRateModifierServer(server, grpcServer)

	return server.Serve(ln)

}

func NewGRPCRateModifier(svc ModifyRate) *GRPCRateModifier {
	return &GRPCRateModifier{
		svc: svc,
	}
}

type KeyValue struct {
    Key   string
    Value int64
}

type ModifyRateImplementation struct{}

func (s *ModifyRateImplementation) RateModifier(ctx context.Context, rateRequest *proto.RateRequest) (*proto.RateResponse, error) {

	keys, err := redis.RedisClient.Keys(ctx, "request_count:*").Result()
	if err != nil {
		log.Fatalf("Failed to fetch keys: %v", err)
	}

	var keyValues []KeyValue

	for _, key := range keys {
		valueStr, err := redis.RedisClient.Get(ctx, key).Result()
		if err != nil {
			log.Printf("Failed to get value for key %s: %v", key, err)
			continue
		}

		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			log.Printf("Failed to parse value for key %s: %v", key, err)
			continue
		}

		keyValues = append(keyValues, KeyValue{Key: key, Value: value})
	}

	
	// Filter and sort by values greater than 5
	rateLimit := int(rateRequest.RateLimit)
	filteredValues := filterValues(keyValues, int64(rateLimit))

	count := len(filteredValues)	
	
	RateLimit = int(rateRequest.RateLimit)
	TimeWindow = time.Duration(rateRequest.Minutes) * time.Minute

	return &proto.RateResponse{
		Count: float32(count),
	}, nil
}


func (s *GRPCRateModifier) RateLimiter(ctx context.Context, in *proto.RateRequest) (*proto.RateResponse, error) {
	return s.svc.RateModifier(ctx, in)
}


func filterValues(keyValues []KeyValue, threshold int64) []KeyValue {
    // Filter keys based on threshold
    var filtered []KeyValue
    for _, kv := range keyValues {
        if kv.Value >= threshold {
            filtered = append(filtered, kv)
        }
    }

    return filtered
}


