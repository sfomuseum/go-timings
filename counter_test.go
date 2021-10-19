package timings

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestCounterMonitor(t *testing.T) {

	ctx := context.Background()

	d := time.Second * 1

	m, err := NewCounterMonitor(ctx, d)

	if err != nil {
		t.Fatalf("Failed to create monitor, %v", err)
	}

	err = m.Start(ctx, os.Stdout)

	if err != nil {
		t.Fatalf("Failed to start monitor, %v", err)
	}

	done_ch := make(chan bool)

	ticker := time.NewTicker(time.Second * 2)
	after := time.After(10 * time.Second)

	var signal_err error

	go func() {

		for {
			select {
			case <-ticker.C:
				err = m.Signal(ctx)

				if err != nil {
					signal_err = err
					done_ch <- true
					return
				}

			case <-after:
				done_ch <- true
				return
			}
		}
	}()

	<-done_ch

	if signal_err != nil {
		t.Fatalf("There was a problem signaling the monitor, %v", signal_err)
	}

	err = m.Stop(ctx)

	if err != nil {
		t.Fatalf("Failed to stop monitor, %v", err)
	}
}
