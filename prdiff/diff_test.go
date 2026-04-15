package prdiff_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
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

func TestClient_GetDiff_mkghtagPR(t *testing.T) {
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

	// https://github.com/suzuki-shunsuke/mkghtag/pull/607
	got, err := c.GetDiff(t.Context(), "suzuki-shunsuke", "mkghtag", 607)
	if err != nil {
		t.Fatalf("GetDiff: %v", err)
	}
	want, err := os.ReadFile("testdata/diff_607_api.txt")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	if d := cmp.Diff(string(want), got); d != "" {
		t.Errorf("diff mismatch (-want +got):\n%s", d)
	}
}
