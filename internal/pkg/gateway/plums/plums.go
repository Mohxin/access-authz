package plums

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/volvo-cars/go-request/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	serviceName    = "plums"
	usersPath      = "users/by-cdsid"
	rolesPath      = "roles"
	attempts       = 3
	requestTimeout = 5 * time.Second
)

const (
	unableToParseBaseURLMessage      = "unable to parse base url: %w"
	unableToGetUserFromPlumsMessage  = "unable to get user from plums: %w"
	unableToGetRolesFromPlumsMessage = "unable to get roles from plums: %w"
)

var ErrUserNotFound = errors.New("plums: user not found")

type collector interface {
	ObserveRequestTimeWithOp(dependency, operation, method, route string, status int, duration time.Duration)
}

type requester interface {
	Get(ctx context.Context, url string, headers map[string]string, opts ...request.RequestOption) (*request.HTTPResponse, error)
}

type tracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}

type PlumsGateway struct {
	cfg           *Config
	tracer        tracer
	client        requester
	userKeyHeader map[string]string
}

func New(cfg *Config, collector collector) *PlumsGateway {
	credentials := clientcredentials.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		TokenURL:     cfg.TokenURL,
		Scopes:       cfg.Scopes,
		EndpointParams: map[string][]string{
			"audience": {cfg.Audience},
		},
	}

	return &PlumsGateway{
		cfg:    cfg,
		tracer: otel.Tracer("gateway/plums"),
		userKeyHeader: map[string]string{
			"user-key": cfg.UserKey,
		},
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

func (g *PlumsGateway) GetUserByCDSID(ctx context.Context, cdsid string) (*User, error) {
	ctx, span := g.tracer.Start(ctx, "plums.GetUserByCDSID", trace.WithAttributes(attribute.String("user_cdsid", cdsid)))
	defer span.End()

	u, err := url.Parse(g.cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf(unableToParseBaseURLMessage, err)
	}

	u.Path = path.Join(u.Path, usersPath, cdsid)

	response, err := g.client.Get(ctx, u.String(), g.userKeyHeader)
	if err != nil {
		return nil, fmt.Errorf(unableToGetUserFromPlumsMessage, err)
	}

	if response.StatusCode == http.StatusNotFound {
		return nil, ErrUserNotFound
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user from plums: %s", response.Body)
	}

	return request.Unmarshal[*User](response.Body)
}

func (g *PlumsGateway) GetRoles(ctx context.Context) (*[]Role, error) {
	u, err := url.Parse(g.cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse BaseUrl: %w", err)
	}

	u.Path = path.Join(u.Path, rolesPath)

	response, err := g.client.Get(ctx, u.String(), g.userKeyHeader)
	if err != nil {
		return nil, fmt.Errorf(unableToGetRolesFromPlumsMessage, err)
	}

	if response.StatusCode == http.StatusNotFound {
		return nil, ErrUserNotFound
	}

	return request.Unmarshal[*[]Role](response.Body)
}
