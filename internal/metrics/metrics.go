package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

type Metrics struct {
	WebhooksRegistered prometheus.Counter
	WebhooksCalled     *prometheus.CounterVec
	ActionsRegistered  *prometheus.CounterVec
	ActionsRun         *prometheus.CounterVec
}

func New() (*Metrics, error) {
	actionsRunCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "sre_operator_actions_run",
	}, []string{"action", "trigger", "status"})
	actionsRegisteredCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "sre_operator_actions_registered",
	}, []string{"type"})
	webhooksRegisteredCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sre_operator_webhooks_registered",
	})
	webhooksCalledCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "sre_operator_webhooks_called",
	}, []string{"id"})

	if err := metrics.Registry.Register(actionsRunCounter); err != nil {
		return nil, err
	}
	if err := metrics.Registry.Register(actionsRegisteredCounter); err != nil {
		return nil, err
	}
	if err := metrics.Registry.Register(webhooksRegisteredCounter); err != nil {
		return nil, err
	}
	if err := metrics.Registry.Register(webhooksCalledCounter); err != nil {
		return nil, err
	}

	return &Metrics{
		WebhooksRegistered: webhooksRegisteredCounter,
		WebhooksCalled:     webhooksCalledCounter,
		ActionsRegistered:  actionsRegisteredCounter,
		ActionsRun:         actionsRunCounter,
	}, nil
}
