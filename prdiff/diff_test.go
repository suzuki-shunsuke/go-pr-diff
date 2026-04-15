package prdiff_test

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/suzuki-shunsuke/go-pr-diff/prdiff"
)

type authTransport struct {
	token string
	base  http.RoundTripper
}

func (t *authTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r2 := r.Clone(r.Context())
	r2.Header.Set("Authorization", "Bearer "+t.token)
	return t.base.RoundTrip(r2) //nolint:wrapcheck
}

func TestClient_GetDiff_aquaRegistryPR(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skip network test in -short mode")
	}

	var hc *http.Client
	if tok := os.Getenv("GITHUB_TOKEN"); tok != "" {
		hc = &http.Client{Transport: &authTransport{token: tok, base: http.DefaultTransport}}
	}

	c, err := prdiff.NewClient(hc, "")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	diff, err := c.GetDiff(t.Context(), "aquaproj", "aqua-registry", 51977)
	if err != nil {
		t.Fatalf("GetDiff: %v", err)
	}
	if diff == "" {
		t.Fatal("got empty diff")
	}
	for _, want := range []string{
		"diff --git",
		"pkgs/trzsz/trzsz-go",
	} {
		if !strings.Contains(diff, want) {
			t.Errorf("diff does not contain %q", want)
		}
	}
}
