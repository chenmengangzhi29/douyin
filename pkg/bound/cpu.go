package bound

import (
	"context"
	"fmt"
	"net"

	"github.com/chenmengangzhi29/douyin/pkg/constants"
	"github.com/chenmengangzhi29/douyin/pkg/errno"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/remote"
	"github.com/shirou/gopsutil/cpu"
)

var _ remote.InboundHandler = &cpuLimitHandler{}

type cpuLimitHandler struct{}

func NewCpuLimitHandler() remote.InboundHandler {
	return &cpuLimitHandler{}
}

// OnActive implements the remote.InboundHandler interface.
func (c *cpuLimitHandler) OnActive(ctx context.Context, conn net.Conn) (context.Context, error) {
	return ctx, nil
}

// OnRead implements the remote.InboundHandler interface.
func (c *cpuLimitHandler) OnRead(ctx context.Context, conn net.Conn) (context.Context, error) {
	p := cpuPercent()
	klog.CtxInfof(ctx, "current cpu is %.2g", p)
	if p > constants.CPURateLimit {
		return ctx, errno.ServiceErr.WithMessage(fmt.Sprintf("cpu = %.2g", c))
	}
	return ctx, nil
}

//OnInactive implements the remote.InboundHandler interface.
func (c *cpuLimitHandler) OnInactive(ctx context.Context, conn net.Conn) context.Context {
	return ctx
}

// OnMessage implements the remote.InboundHandler interface.
func (c *cpuLimitHandler) OnMessage(ctx context.Context, args, result remote.Message) (context.Context, error) {
	return ctx, nil
}

func cpuPercent() float64 {
	percent, _ := cpu.Percent(0, false)
	return percent[0]
}
