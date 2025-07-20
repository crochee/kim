package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/sourcegraph/conc/pool"
	"github.com/spf13/viper"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/crochee/kim/cmd"
	"github.com/crochee/kim/internal/tracing"
)

func runRoot(ctx context.Context) error {
	issuer := fmt.Sprintf("http://localhost:%s/", "89000")
	r, err := SetupServer(issuer, "", []string{"http://localhost:3000/"})
	if err != nil {
		return err
	}
	g := pool.New().WithContext(ctx).WithCancelOnError()
	g.Go(func(ctx context.Context) error {
		return cmd.Operator(ctx)
	})
	g.Go(func(ctx context.Context) error {
		return trace(ctx)
	})
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
		BaseContext: func(net.Listener) context.Context {
			return logf.IntoContext(ctx, mainLog)
		},
	}
	g.Go(func(ctx context.Context) error {
		return srv.ListenAndServe()
	})
	g.Go(func(ctx context.Context) error {
		return srv.Shutdown(ctx)
	})
	if err := g.Wait(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func trace(ctx context.Context) error {
	otelEndponint := viper.GetString("otel-endpoint")
	if otelEndponint == "" {
		return nil
	}
	mainLog.Info("Initializing OpenTelemetry tracer provider", "otel-endpoint", otelEndponint)
	tp, err := tracing.InitTracer(ctx, otelEndponint)
	if err != nil {
		mainLog.Error(err, "Failed to initialize tracer")
		return err
	}
	defer func(ctx context.Context) {
		<-ctx.Done()
		if err := tp.Shutdown(ctx); err != nil {
			mainLog.Error(err, "Failed to shutdown tracer provider")
		}
	}(ctx)
	return nil
}
