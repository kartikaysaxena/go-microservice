package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kartikaysaxena/go-microservice/client"
	ops "github.com/kartikaysaxena/go-microservice/operations"
	"github.com/kartikaysaxena/go-microservice/proto"
	redis "github.com/kartikaysaxena/go-microservice/redis"
	server "github.com/kartikaysaxena/go-microservice/server"
	log "github.com/sirupsen/logrus"
	redisLib "github.com/go-redis/redis/v8"
)


var (
	grpcAddr = flag.String("grpc", ":4000", "gRPC server address")
)



func getIPAddr(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}

	return strings.Split(IPAddress, ":")[0]
}

// User represents a user entit


// func RateLimiter(rateLimit float32, timeWindow float32) error {

// }

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		realIP := getIPAddr(r)

		log.WithFields(log.Fields{
			"IP": realIP,
		}).Info("Request received")

		log.Printf("Started %s %s", r.Method, r.RequestURI)

		ctx := context.Background()

		redisKey := "request_count:" + realIP


		count, err := redis.RedisClient.Get(ctx, redisKey).Int()
		if err != nil && err != redisLib.Nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		log.WithFields(log.Fields{
			"count": count,
		}).Infof("Request count for IP %s", realIP)


		if count >= server.RateLimit {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		pipe := redis.RedisClient.TxPipeline()
		pipe.Incr(ctx, redisKey)

		if count == 0 {
			pipe.Expire(ctx, redisKey, server.TimeWindow)
		}
		_, err = pipe.Exec(ctx)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)

		log.Printf("Completed in %v", time.Since(start))
	})
}

// CreateUser creates a new user

func main() {
	redis.RedisClient = redis.InitRedis()


	grpcClient, err := client.NewGRPCClient(":4000")
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}

	go func ()  {
		for {
			time.Sleep(5 * time.Second)
			resp, err := grpcClient.RateLimiter(context.Background(), &proto.RateRequest{
				RateLimit: 5,
				Minutes: 1,
			})
			if err != nil {
				log.Fatalf("Failed to rate limit: %v", err)
			}
			fmt.Println("Count:", resp.Count)
		}
	}()

	go server.MakeGRPCServerAndRun(server.NewRateModifier(&server.ModifyRateImplementation{}), *grpcAddr)

	http.Handle("/user/list", loggingMiddleware(http.HandlerFunc(ops.ListUsers)))
	http.Handle("/user/create", loggingMiddleware(http.HandlerFunc(ops.CreateUser)))
	http.Handle("/user/update", loggingMiddleware(http.HandlerFunc(ops.UpdateUser)))
	http.Handle("/user/delete", loggingMiddleware(http.HandlerFunc(ops.DeleteUser)))

	log.WithFields(log.Fields{
		"port": 8080,
	}).Info("Starting server")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
