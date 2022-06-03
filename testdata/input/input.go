package main

import (
	"context"
	"errors"
	"time"

	"github.com/calyptia/plugin"
	"github.com/calyptia/plugin/metric"
)

func init() {
	plugin.RegisterInput("go-test-input-plugin", "Golang input plugin for testing", &inputPlugin{})
}

type inputPlugin struct {
	foo            string
	collectCounter metric.Counter
}

func (plug *inputPlugin) Init(ctx context.Context, conf plugin.ConfigLoader, metrics plugin.Metrics) error {
	plug.foo = conf.String("foo")
	plug.collectCounter = metrics.NewCounter("collect_total", "Total number of collects", "go-test-input-plugin")
	return nil
}

func (plug inputPlugin) Collect(ctx context.Context, ch chan<- plugin.Message) error {
	tick := time.NewTicker(time.Second)

	for {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			if err != nil && !errors.Is(err, context.Canceled) {
				return err
			}

			return nil
		case <-tick.C:
			plug.collectCounter.Add(1)

			ch <- plugin.Message{
				Time: time.Now(),
				Record: map[string]string{
					"message": "hello from go-test-input-plugin",
					"foo":     plug.foo,
				},
			}
		}
	}
}

func main() {}
