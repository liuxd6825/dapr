/*
Copyright 2021 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"github.com/dapr/kit/logger"
	"github.com/liuxd6825/dapr/cmd/operator/options"
	"github.com/liuxd6825/dapr/pkg/buildinfo"
	"github.com/liuxd6825/dapr/pkg/concurrency"
	"github.com/liuxd6825/dapr/pkg/metrics"
	"github.com/liuxd6825/dapr/pkg/operator"
	"github.com/liuxd6825/dapr/pkg/operator/monitoring"
	"github.com/liuxd6825/dapr/pkg/signals"
)

var log = logger.NewLogger("dapr.operator")

func main() {
	opts := options.New()

	// Apply options to all loggers.
	if err := logger.ApplyOptionsToLoggers(&opts.Logger); err != nil {
		log.Fatal(err)
	}

	log.Infof("Starting Dapr Operator -- version %s -- commit %s", buildinfo.Version(), buildinfo.Commit())
	log.Infof("Log level set to: %s", opts.Logger.OutputLevel)

	metricsExporter := metrics.NewExporterWithOptions(log, metrics.DefaultMetricNamespace, opts.Metrics)

	if err := monitoring.InitMetrics(); err != nil {
		log.Fatal(err)
	}

	ctx := signals.Context()
	op, err := operator.NewOperator(ctx, operator.Options{
		Config:                              opts.Config,
		TrustAnchorsFile:                    opts.TrustAnchorsFile,
		LeaderElection:                      !opts.DisableLeaderElection,
		WatchdogMaxRestartsPerMin:           opts.MaxPodRestartsPerMinute,
		WatchNamespace:                      opts.WatchNamespace,
		ServiceReconcilerEnabled:            !opts.DisableServiceReconciler,
		ArgoRolloutServiceReconcilerEnabled: opts.EnableArgoRolloutServiceReconciler,
		WatchdogEnabled:                     opts.WatchdogEnabled,
		WatchdogInterval:                    opts.WatchdogInterval,
		WatchdogCanPatchPodLabels:           opts.WatchdogCanPatchPodLabels,
	})
	if err != nil {
		log.Fatalf("error creating operator: %v", err)
	}

	err = concurrency.NewRunnerManager(
		metricsExporter.Run,
		op.Run,
	).Run(ctx)
	if err != nil {
		log.Fatalf("error running operator: %v", err)
	}

	log.Info("operator shut down gracefully")
}
