package bus

import (
	"context"
	"time"
)

type Response struct {
	Response any
	Err      error
}

type dispatchableWithContext struct {
	ctx context.Context
	d   Dispatchable

	rsChan chan Response
}

// ConcurrentBus dispatches concurrently
type ConcurrentBus struct {
	h        map[string]Handler
	timeout  time.Duration
	poolSize chan struct{}
	in       chan dispatchableWithContext
}

// NewConcurrentBus is a constructor
func NewConcurrentBus(timeout time.Duration, concurrencyLimit int) ConcurrentBus {
	return ConcurrentBus{
		h:        make(map[string]Handler),
		timeout:  timeout,
		poolSize: make(chan struct{}, concurrencyLimit),
		in:       make(chan dispatchableWithContext),
	}
}

// CurrentSize is a getter
func (b ConcurrentBus) CurrentSize() int {
	return len(b.poolSize)
}

// Register adds a new handler to the bus
func (b ConcurrentBus) Register(n string, h Handler) {
	b.h[n] = h
}

// Dispatch dispatches a new dispatchable item
func (b ConcurrentBus) Dispatch(ctx context.Context, d Dispatchable) <-chan Response {
	rsChan := make(chan Response, 1)
	b.in <- dispatchableWithContext{
		ctx:    ctx,
		d:      d,
		rsChan: rsChan,
	}
	return rsChan
}

// Run start running the bus
func (b ConcurrentBus) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case dwc := <-b.in:
			h, ok := b.h[dwc.d.Name()]
			if !ok {
				dwc.rsChan <- Response{
					Err: ErrNotDispatchable,
				}
				continue
			}
			b.poolSize <- struct{}{}

			go func() {
				withTimeoutCtx, cancel := context.WithTimeout(dwc.ctx, b.timeout)
				defer cancel()
				defer func() {
					<-b.poolSize
				}()
				rs, err := h(withTimeoutCtx, dwc.d)
				dwc.rsChan <- Response{
					Response: rs,
					Err:      err,
				}
			}()
		}
	}
}
