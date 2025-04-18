//go:build allcomponents || stablecomponents

package components

import (
	"github.com/dapr/components-contrib/liuxd/nameresoluton/nacos"
	nrLoader "github.com/dapr/dapr/pkg/components/nameresolution"
)

func init() {
	nrLoader.DefaultRegistry.RegisterComponent(nacos.NewResolver, "nacos")
}
