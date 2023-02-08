package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"

	"github.com/milas/gha-oidc-bridge/pkg/api/gha"
)

func main() {
	ctx := context.Background()

	const hostname = ""
	const port = "45321"

	e := echo.New()
	e.HideBanner = true
	e.Debug = true

	e.POST("/oidc/gha", oidcHandler(ctx))

	e.Logger.Fatal(e.Start(net.JoinHostPort(hostname, port)))
}

func oidcHandler(ctx context.Context) echo.HandlerFunc {
	provider, err := oidc.NewProvider(ctx, "https://token.actions.githubusercontent.com")
	if err != nil {
		log.Fatalf("oidc discovery: %v", err)
	}
	verifier := provider.Verifier(&oidc.Config{SkipClientIDCheck: true})
	return func(c echo.Context) error {
		var tokenReq gha.TokenExchangeRequest
		if err := c.Bind(&tokenReq); err != nil {
			c.Logger().Warnf("Binding request: %v", err)
			return err
		}

		if tokenReq.Value == "" {
			return echo.ErrBadRequest
		}

		_, err := verifier.Verify(ctx, tokenReq.Value)
		if err != nil {
			c.Logger().Warnf("Verifying token: %v", err)
			return echo.ErrUnauthorized
		}

		return c.NoContent(http.StatusNoContent)
	}
}
