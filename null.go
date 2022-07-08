package timings

import (
       "io"
       "context"
)

type NullMonitor struct {
     Monitor
}

func NewNullMonitor(ctx context.Context) (Monitor, error) {
     nm := &NullMonitor{}
     return nm, nil
}

func (nm *NullMonitor) Start(ctx context.Context, wr io.Writer) error {
     return nil
}

func (nm *NullMonitor) Stop(ctx context.Context) error {
     return nil
}

func (nm *NullMonitor) Signal(ctx context.Context, args ...interface{}) error {
     return nil
}
