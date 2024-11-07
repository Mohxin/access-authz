package cachemanager

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/volvo-cars/go-request/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	serviceName    = "cache-manager"
	attempts       = 3
	requestTimeout = 5 * time.Second
)

type PartnerType string

const (
	PartnerTypeNone  PartnerType = "NONE"
	PartnerTypeParma PartnerType = "PARMA"
	PartnerTypeNsc   PartnerType = "NSC"
)

var PartnerTypeNames = map[string]PartnerType{
	"NONE":  PartnerTypeNone,
	"PARMA": PartnerTypeParma,
	"NSC":   PartnerTypeNsc,
}

func (p PartnerType) String() string {
	return string(p)
}

type requester interface {
	Get(ctx context.Context, url string, headers map[string]string, opts ...request.RequestOption) (*request.HTTPResponse, error)
}

type collector interface {
	ObserveRequestTimeWithOp(dependency, operation, method, route string, status int, duration time.Duration)
}

type tracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}

type CacheGateway struct {
	cfg    *Config
	tracer tracer
	client requester
}

func New(cfg *Config, collector collector) *CacheGateway {
	credentials := clientcredentials.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		TokenURL:     cfg.TokenURL,
		Scopes:       cfg.Scopes,
	}

	return &CacheGateway{
		cfg:    cfg,
		tracer: otel.Tracer("gateway/cache"),
		client: request.NewRequestHandler(
			request.WithHTTPClient(request.NewPooledOAuth2ClientWithOtel(credentials, request.WithServiceAttribute(serviceName))),
			request.WithAttempts(attempts),
			request.WithTimeout(requestTimeout),
			request.OnResponse(func(req *http.Request, res *http.Response, operation string, duration time.Duration) {
				collector.ObserveRequestTimeWithOp(serviceName, operation, req.Method, req.URL.Path, res.StatusCode, duration)
			}),
		),
	}
}

func (g *CacheGateway) GetPartnersByCodes(ctx context.Context, partnerCodes []string, partnerType string) ([]*Partner, error) {
	codes := strings.Join(partnerCodes, ",")
	operationName := "cache.GetPartnersByCodes"
	ctx, span := g.tracer.Start(ctx, operationName, trace.WithAttributes(attribute.String("partnerCode", codes), attribute.String("partnerType", partnerType)))
	defer span.End()

	URL, err := url.Parse(g.cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	query := URL.Query()
	query.Add("codes", codes)
	query.Add("type", toPartnerType(partnerType))
	URL.RawQuery = query.Encode()

	opts := []request.RequestOption{
		request.WithOperation(operationName),
	}

	resp, err := g.client.Get(ctx, URL.String(), nil, opts...)
	if err != nil {
		return nil, err
	}

	obj, err := request.Unmarshal[PartnersResponse](resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get partner response code:[%d] error message:%s", resp.StatusCode, obj.Error.Message)
	}

	return obj.Data, nil
}

func toPartnerType(typ string) string {
	if t, ok := PartnerTypeNames[typ]; ok {
		return t.String()
	}
	return PartnerTypeNone.String()
}
