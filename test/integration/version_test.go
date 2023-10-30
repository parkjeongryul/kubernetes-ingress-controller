//go:build integration_tests

package integration

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
)

func RunWhenKongVersion(t *testing.T, vRangeStr string, msg ...any) {
	t.Helper()

	vRange, err := kong.NewRange(vRangeStr)
	require.NoError(t, err)

	version := eventuallyGetKongVersion(t, proxyAdminURL)

	// We could parse version, clear the rc/alpha/beta suffixes and then compare
	// it but it seems unnecessary since gateway dev pre release images coming from
	// kong/kong-gateway-dev report the final version through Admin API anyway.
	// So when running 3.3.0.0-rc.3 we'll get 3.3.0.0.

	if !vRange(version) {
		if len(msg) > 0 {
			t.Log(msg...)
		}
		t.Skipf("skipping because Kong version %q is not within test's range %q: ", version, vRangeStr)
	}
}

func RunWhenKongDBMode(t *testing.T, dbmode string, msg ...any) {
	t.Helper()

	actual := eventuallyGetKongDBMode(t, proxyAdminURL)

	if actual != dbmode {
		if len(msg) > 0 {
			t.Log(msg...)
		}
		t.Skipf("skipping because Kong dbmode %q is different than requested %q", actual, dbmode)
	}
}

func RunWhenKongEnterprise(t *testing.T) {
	t.Helper()

	version := eventuallyGetKongVersion(t, proxyAdminURL)

	if !version.IsKongGatewayEnterprise() {
		t.Skipf("skipping because Kong is not running as Enterprise, detected version %q", version)
	}
}

func RunWhenKongExpressionRouter(t *testing.T) {
	if routerFlavor := eventuallyGetKongRouterFlavor(t, proxyAdminURL); routerFlavor != kongRouterFlavorExpressions {
		t.Skipf("skip test because expression router is disabled (current router flavor is: %q)", routerFlavor)
	}
}

func eventuallyGetKongVersion(t *testing.T, adminURL *url.URL) kong.Version {
	t.Helper()

	var (
		err     error
		version kong.Version
	)

	require.EventuallyWithT(t, func(t *assert.CollectT) {
		ctx, cancel := context.WithTimeout(context.Background(), test.RequestTimeout)
		defer cancel()
		version, err = helpers.GetKongVersion(ctx, adminURL, consts.KongTestPassword)
		assert.NoError(t, err)
	}, time.Minute, time.Second)
	return version
}

func eventuallyGetKongDBMode(t *testing.T, adminURL *url.URL) string {
	t.Helper()

	var (
		err    error
		dbmode string
	)

	require.EventuallyWithT(t, func(t *assert.CollectT) {
		ctx, cancel := context.WithTimeout(context.Background(), test.RequestTimeout)
		defer cancel()
		dbmode, err = helpers.GetKongDBMode(ctx, adminURL, consts.KongTestPassword)
		assert.NoError(t, err)
	}, time.Minute, time.Second)
	return dbmode
}

func eventuallyGetKongRouterFlavor(t *testing.T, adminURL *url.URL) string {
	t.Helper()

	var (
		err          error
		routerFlavor string
	)

	require.EventuallyWithT(t, func(t *assert.CollectT) {
		ctx, cancel := context.WithTimeout(context.Background(), test.RequestTimeout)
		defer cancel()
		routerFlavor, err = helpers.GetKongRouterFlavor(ctx, adminURL, consts.KongTestPassword)
		assert.NoError(t, err)
	}, time.Minute, time.Second)
	return routerFlavor
}
