package new_relic

import (
	"context"
	"github.com/erkanzileli/rate-limiter/configs"
	"github.com/newrelic/go-agent/v3/newrelic"
	"runtime"
)

var (
	agent     *newrelic.Application
	emptyFunc = func() {}
)

func CreateAgent() (err error) {
	if configs.Config.Tracing.Enabled {
		newRelicConfig := configs.Config.Tracing.NewRelic
		agent, err = newrelic.NewApplication(
			newrelic.ConfigAppName(newRelicConfig.AppName),
			newrelic.ConfigLicense(newRelicConfig.LicenseKey),
			newrelic.ConfigDistributedTracerEnabled(newRelicConfig.DistributedTracerEnabled),
		)
	}
	return
}

func StartTransaction(ctx context.Context, name string) (context.Context, func()) {
	if agent == nil {
		return ctx, emptyFunc
	}

	txn := agent.StartTransaction(name)
	ctx = newrelic.NewContext(ctx, txn)
	return ctx, txn.End
}

func StartSegment(ctx context.Context) func() {
	txn := newrelic.FromContext(ctx)
	if txn == nil {
		return emptyFunc
	}

	segment := newrelic.Segment{
		Name:      getCallerName(),
		StartTime: txn.StartSegmentNow(),
	}
	return segment.End
}

func StartSegmentWithName(ctx context.Context, name string) func() {
	txn := newrelic.FromContext(ctx)
	if txn == nil {
		return emptyFunc
	}

	segment := newrelic.Segment{
		Name:      name,
		StartTime: txn.StartSegmentNow(),
	}
	return segment.End
}

func getCallerName() string {
	pc, _, _, ok := runtime.Caller(2)
	details := runtime.FuncForPC(pc)

	if ok && details != nil {
		return details.Name()
	}

	return "unknown"
}
