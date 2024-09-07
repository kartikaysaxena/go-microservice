package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/anthdm/pricefetcher/client"
	"github.com/anthdm/pricefetcher/proto"
)

func main() {
	// client := client.New("http://localhost:3000")

	// price, err := client.FetchPrice(context.Background(), "ET")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("%+v\n", price)

	var (
		jsonAddr = flag.String("listenaddr", ":3000", "listen address the service is running")
		grpcAddr = flag.String("grpc", ":4000", "listen address of the grpc transport")
		svc      = NewLoggingService(NewMetricService(&priceFetcher{}))
		ctx      = context.Background()
	)

	flag.Parse()

	grpcClient, err := client.NewGRPCClient(":4000")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			time.Sleep(3 * time.Second)
			resp, err := grpcClient.FetchPrice(ctx, &proto.PriceRequest{Ticker: "BTC"})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("\n", resp)
		}
	}()

	go makeGRPCServerAndRun(svc, *grpcAddr)

	jsonServer := NewJSONAPIServer(*jsonAddr, svc)
	jsonServer.Run()
}
