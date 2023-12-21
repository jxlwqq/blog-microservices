package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcprometheus "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	grpclogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	grpcvalidator "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	v1 "github.com/jxlwqq/blog-microservices/api/protobuf/post/v1"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	flag "github.com/spf13/pflag"
	_ "go.uber.org/automaxprocs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var flagConfig = flag.String("config", "./configs/config.yaml", "config file")

func main() {
	flag.Parse()
	logger := log.New()
	conf, err := config.Load(*flagConfig)
	if err != nil {
		logger.Fatal(err)
	}

	// Setup metrics.
	srvMetrics := grpcprometheus.NewServerMetrics(
		grpcprometheus.WithServerHandlingTimeHistogram(
			grpcprometheus.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)
	reg := prometheus.NewRegistry()
	reg.MustRegister(srvMetrics)

	postServer, err := InitServer(logger, conf)
	if err != nil {
		logger.Fatal(err)
	}
	healthServer := health.NewServer()

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		grpcrecovery.UnaryServerInterceptor(),
		srvMetrics.UnaryServerInterceptor(),
		grpcvalidator.UnaryServerInterceptor(),
		grpclogging.UnaryServerInterceptor(logger),
	))

	v1.RegisterPostServiceServer(grpcServer, postServer)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	lis, err := net.Listen("tcp", conf.Post.Server.GRPC.Port)

	// Start gRPC server
	ch := make(chan os.Signal, 1)
	logger.Infof("gPRC Listening on port %s", conf.Post.Server.GRPC.Port)
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	// todo Start HTTP server

	// Start Metrics server
	logger.Infof("Metrics Listening on port %s", conf.Post.Server.Metrics.Port)
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics e.g. to support exemplars.
			EnableOpenMetrics: true,
		},
	))
	metricsServer := &http.Server{
		Addr:    conf.Post.Server.Metrics.Port,
		Handler: metricsMux,
	}
	go func() {
		if err = metricsServer.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			logger.Fatal(err)
		}
	}()

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	grpcServer.GracefulStop()
	if err = metricsServer.Shutdown(ctx); err != nil {
		logger.Fatal(err)
	}
	<-ctx.Done()
	close(ch)
	fmt.Println("Graceful Shutdown end")
}
