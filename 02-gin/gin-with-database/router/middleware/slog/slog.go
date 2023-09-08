package slog

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDCtx = "slog-gin.request-id"

// Builder 使用 Slog
type Builder struct {
	logger *slog.Logger
	level  slog.Level

	withRequestID bool
}

func NewBuilder(logger *slog.Logger) *Builder {
	return &Builder{
		logger: logger,
		level:  slog.LevelInfo,

		withRequestID: true,
	}
}

func (b *Builder) RequestID(enable bool) *Builder {
	b.withRequestID = enable
	return b
}

func (b *Builder) SetLevel(level slog.Level) *Builder {
	b.level = level
	return b
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		requestID := uuid.New().String()
		if b.withRequestID {
			c.Set(requestIDCtx, requestID)
			c.Header("X-Request-ID", requestID)
		}

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		attributes := []slog.Attr{
			slog.Int("status", c.Writer.Status()),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("ip", c.ClientIP()),
			slog.Duration("latency", latency),
			slog.String("user-agent", c.Request.UserAgent()),
			slog.Time("time", end),
		}

		if b.withRequestID {
			attributes = append(attributes, slog.String("request-id", requestID))
		}

		switch {
		case c.Writer.Status() >= http.StatusBadRequest && c.Writer.Status() < http.StatusInternalServerError:
			b.logger.LogAttrs(context.Background(), b.level, c.Errors.String(), attributes...)
		case c.Writer.Status() >= http.StatusInternalServerError:
			b.logger.LogAttrs(context.Background(), b.level, c.Errors.String(), attributes...)
		default:
			b.logger.LogAttrs(context.Background(), b.level, "Incoming request", attributes...)
		}
	}
}
