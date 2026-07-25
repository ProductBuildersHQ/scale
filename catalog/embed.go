// Package catalog embeds the reference SCALE framework catalog: the
// framework document plus the seed domains (api, observability, security).
package catalog

import (
	"embed"

	scale "github.com/ProductBuildersHQ/scale"
)

//go:embed framework.json domains/*.json external/*.json
var catalogFS embed.FS

// Default loads and validates the embedded reference framework.
func Default() (*scale.Framework, error) {
	return scale.LoadFrameworkFS(catalogFS, ".")
}
