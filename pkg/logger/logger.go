package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

// formatTimeAttr 格式化时间属性
func formatTimeAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		// 格式化时间为 yyyy-MM-dd HH:mm:ss:SSS
		return slog.Attr{
			Key:   slog.TimeKey,
			Value: slog.StringValue(a.Value.Time().Format("2006-01-02 15:04:05.000")),
		}
	}
	return a
}

func Init(env string) {
	if env == "" || (env != "dev" && env != "prod") {
		env = "prod"
	}

	var baseHandler slog.Handler
	handlerOpts := &slog.HandlerOptions{
		ReplaceAttr: formatTimeAttr,
	}

	if env == "dev" {
		handlerOpts.Level = slog.LevelDebug
		baseHandler = slog.NewTextHandler(os.Stdout, handlerOpts)
	} else {
		handlerOpts.Level = slog.LevelInfo
		baseHandler = slog.NewJSONHandler(os.Stdout, handlerOpts)
	}

	Log = slog.New(baseHandler).With(
		"service", "ruoyi-go",
		"version", "v1.0.0",
	)

	// 设置为全局默认 logger（可选但推荐）
	slog.SetDefault(Log)
}
