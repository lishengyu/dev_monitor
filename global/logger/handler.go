package logger

import (
	"context"
	"log/slog"

	"github.com/petermattis/goid" // 协程ID库
)

// 实现 slog.Handler 接口
type CustomHandler struct {
	slog.Handler
	opts *Options
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	// 添加协程ID
	r.AddAttrs(slog.Int("goroutine", int(goid.Get())))

	/*
		// 添加源代码位置
		if r.PC != 0 {
			fs := runtime.CallersFrames([]uintptr{r.PC})
			frame, _ := fs.Next()
			r.AddAttrs(
				slog.String("file", filepath.Base(frame.File)),
				slog.Int("line", frame.Line),
			)
		}
	*/

	return h.Handler.Handle(ctx, r)
}

func NewCustomHandler(h slog.Handler, opts *Options) *CustomHandler {
	return &CustomHandler{
		Handler: h,
		opts:    opts,
	}
}
