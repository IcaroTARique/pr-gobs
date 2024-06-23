package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/IcaroTARique/pr-locate-weather/configs"
	"github.com/IcaroTARique/pr-locate-weather/internal/infra/api/cep"
	"github.com/IcaroTARique/pr-locate-weather/internal/infra/api/weather"
	"github.com/IcaroTARique/pr-locate-weather/internal/infra/webserver/handler"
	"github.com/IcaroTARique/pr-locate-weather/internal/observability"
	InternalZipkin "github.com/IcaroTARique/pr-locate-weather/internal/observability/zipkin"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	conf, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	observabilityProvider := observability.NewProvider(conf.OtelServiceName, conf.OtelExporterEndpoint)

	//Zipkin
	url := flag.String("zipkin", fmt.Sprintf("http://%s:9411/api/v2/spans", conf.ZipkinDomain), "zipkin url")
	zipkinProvider := InternalZipkin.NewProvider(*url)
	flag.Parse()

	sigCh := make(chan os.Signal, 1)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	conn, err := observabilityProvider.InitConn()
	if err != nil {
		log.Fatal(err)
	}
	var serviceName = semconv.ServiceNameKey.String(conf.OtelServiceName)
	res, err := resource.New(
		ctx,
		resource.WithAttributes(serviceName),
	)

	//Zipkin
	zipShutdown, zipTracer, err := zipkinProvider.InitTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := zipShutdown(ctx); err != nil {
			log.Fatalf("failed to zipShutdown TracerProvider: %s", err)
		}
	}()

	//Zipkin
	zipWeatherShutdown, zipWeatherTracer, err := zipkinProvider.InitTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := zipWeatherShutdown(ctx); err != nil {
			log.Fatalf("failed to zipShutdown TracerProvider: %s", err)
		}
	}()

	shutdown, err := observabilityProvider.InitTracerProvider(ctx, res, conn)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %s", err)
		}
	}()
	shutdownMeterProvider, err := observabilityProvider.InitMeterProvider(ctx, res, conn)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdownMeterProvider(ctx); err != nil {
			log.Fatalf("failed to shutdown MeterProvider: %s", err)
		}
	}()

	tracer := otel.Tracer("microsservice1-tracer")
	//Zipkin
	ztracer := zipTracer.Tracer("zipkin-tracer")
	zWeatherTrace := zipWeatherTracer.Tracer("zipkin-weather-tracer")

	apiCep := cep.NewApiCep(conf.OtelServiceName)
	apiWeather := weather.NewApiWeather(conf.OtelServiceName)
	temperatureHandler := handler.NewApiTemperatureResponse(apiCep, apiWeather, tracer, ztracer, zWeatherTrace, conf.OtelServiceName)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Handle("/metrics", promhttp.Handler())

	r.Get("/temperature/{cep}", temperatureHandler.GetTemperatureHandler)

	go func() {
		if err := http.ListenAndServe(":"+conf.WebServerPort, r); err != nil {
			panic(err)
		}
	}()

	select {
	case <-sigCh:
		log.Println("Shutting down gracefully, CTRL+C pressed...")
	case <-ctx.Done():
		log.Println("Shutting down due to other reason...")
	}

	// Create a timeout context for the graceful shutdown
	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

}
