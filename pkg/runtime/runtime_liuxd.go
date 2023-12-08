package runtime

import (
	"context"
	"github.com/liuxd6825/dapr/pkg/apphealth"
)

func (a *DaprRuntime) AppHealthChanged(ctx context.Context) {
	a.appHealthChanged(ctx, apphealth.AppStatusHealthy)
}
