package prdiff_test

import (
	"context"
	"fmt"
	"log"

	"github.com/suzuki-shunsuke/go-pr-diff/prdiff"
)

func Example() {
	c, err := prdiff.NewClient(nil, "")
	if err != nil {
		log.Fatal(err)
	}
	// Get diff of https://github.com/suzuki-shunsuke/mkghtag/pull/1080
	diff, err := c.GetDiff(context.Background(), "suzuki-shunsuke", "mkghtag", 1080)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(diff)
	// Output:
	// diff --git a/.github/workflows/wc-test.yaml b/.github/workflows/wc-test.yaml
	// index 1f703973..cc980f16 100644
	// --- a/.github/workflows/wc-test.yaml
	// +++ b/.github/workflows/wc-test.yaml
	// @@ -12,7 +12,7 @@ jobs:
	//          uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd # v6.0.2
	//          with:
	//            persist-credentials: false
	// -      - uses: actions/setup-go@4b73464bb391d4059bd26b0524d20df3927bd417 # v6.3.0
	// +      - uses: actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c # v6.4.0
	//          with:
	//            go-version-file: go.mod
	//            cache: true
}
