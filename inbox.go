package insteon

import "context"

type inbox struct {
	C chan *packet

	ctx    context.Context
	cancel func()
}

func newInbox(ctx context.Context) *inbox {
	ctx, cancel := context.WithCancel(ctx)

	return &inbox{
		C:      make(chan *packet),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (i *inbox) Done() <-chan struct{} {
	return i.ctx.Done()
}

func (i *inbox) close() {
	i.cancel()
}
