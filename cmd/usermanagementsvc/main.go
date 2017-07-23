package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/basvanbeek/pubsub/kafka"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"

	um "github.com/steven-ferrer/kafka-pubsub"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":10401", "http address of service")
		broker   = flag.String("kafka.broker", "localhost:9092", "kafka broker address")
	)
	flag.Parse()

	//Logging domain
	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stdout)
		logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}
	logger.Log("msg", "Service started")
	defer logger.Log("msg", "Service ended")

	//Business domain
	var service um.Service
	{
		topic := "UserCreate"

		pubber, err := kafka.NewKafkaPublisher(
			*broker,
			topic,
		)
		if err != nil {
			log.Fatalf(err.Error())
		}

		pubbers := um.Publishers{CreateUserPublisher: pubber}
		service = um.NewService(pubbers)
		service = um.ServiceLoggingMiddleware(logger)(service)
	}

	//subscribers
	var subbers um.Subscribers
	{
		topic := "UserCreate"

		createUserSubber, err := kafka.NewSubscriber(
			*broker,
			topic,
			kafka.OffsetOldest(),
		)
		if err != nil {
			log.Fatalf(err.Error())
		}

		subbers = um.Subscribers{
			CreateUserSubscriber: createUserSubber,
		}
	}

	//endpoints
	var createUserEndpoint endpoint.Endpoint
	{
		createUserEndpoint = um.MakeCreateUserEndpoint(service)
	}

	endpoints := um.Endpoints{
		CreateUserEndpoint: createUserEndpoint,
	}

	errc := make(chan error)

	//interrupt handler
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	//start the subscribers
	go func() {
		errSubbc := subbers.Start()
		errc <- (<-errSubbc)
	}()

	//http transport
	go func() {
		h := um.MakeHTTPHandler(endpoints)
		l := kitlog.With(logger, "transrpot", "http")
		l.Log("addr", *httpAddr)
		errc <- http.ListenAndServe(*httpAddr, h)
	}()

	log.Fatalf("%s", <-errc)
}
