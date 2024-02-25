package trace

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"strings"
	"sync"
)

var (
	once            sync.Once
	defaultProvider trace.TracerProvider
	namedTracer     = map[string]trace.Tracer{}
	namedTracerLock sync.RWMutex
	promCounterVec  *prometheus.CounterVec
	gVersion        string
)

func StartTracer(project, srvName, jaegerAddr, version string) {
	if project == "" || srvName == "" || jaegerAddr == "" {
		slog.Error("not init jaeger.because project/srvName/jaegerAddr is empty")
		return
	}
	gVersion = version
	once.Do(func() {
		promCounterVec = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("%s_trace_span_count", project),
				Help: fmt.Sprintf("%s_trace_span_count", project),
			},
			nil,
		)
		prometheus.MustRegister(promCounterVec)
		prom.AddCollector(promCounterVec)

		// Create the Jaeger exporter
		var endpointOption jaeger.EndpointOption
		if strings.HasPrefix(jaegerAddr, "http") {
			// HTTP.
			endpointOption = jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerAddr))
		} else {
			tmp := strings.SplitN(jaegerAddr, ":", 2)
			var (
				host = tmp[0]
				port = tmp[1]
			)

			endpointOption = jaeger.WithAgentEndpoint(
				jaeger.WithAgentHost(host), jaeger.WithAgentPort(port),
			)
		}

		exp, err := jaeger.New(endpointOption)
		if err != nil {
			return
		}
		var opts []sdktrace.TracerProviderOption
		// Always be sure to batch in production.
		opts = append(opts, sdktrace.WithBatcher(&spanExporter{exp}))
		opts = append(opts, sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(srvName),
		)))

		defaultProvider = sdktrace.NewTracerProvider(opts...)
		otel.SetTracerProvider(defaultProvider)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, propagation.Baggage{}))
		otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
			xlog.Errorf(context.TODO(), "[otel] error: %s ", err)
		}))
	})
	return
}

func getTracerProvider(name string) trace.Tracer {
	if defaultProvider == nil {
		return nil
	}
	namedTracerLock.RLock()
	t, ok := namedTracer[name]
	namedTracerLock.RUnlock()
	if ok {
		return t
	}
	namedTracerLock.Lock()
	if t = defaultProvider.Tracer(name, trace.WithInstrumentationVersion(gVersion)); t != nil {
		namedTracer[name] = t
	}
	namedTracerLock.Unlock()

	return t
}
