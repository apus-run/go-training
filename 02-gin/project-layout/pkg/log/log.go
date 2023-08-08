package log

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/gin-gonic/gin"
)

const LOGGER_KEY = "zapLogger"

type Logger struct {
	logger *zap.Logger
}

func NewLogger(opts ...Option) *Logger {
	options := Apply(opts...)

	// 日志文件切割归档
	// writerSyncer := getLogWriter(opts...)
	writerSyncer := getLogConsoleWriter(options)

	// 编码器配置
	encoder := getEncoder(options.encoding)

	// 日志级别
	level := getLogLevel(options.logLevel)

	core := zapcore.NewCore(encoder, writerSyncer, level)
	if options.mode != "prod" {
		return &Logger{
			zap.New(core, zap.Development(), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)),
		}
	}
	return &Logger{
		zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)),
	}
}

func getEncoder(encoding string) zapcore.Encoder {
	if encoding == "console" {
		return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "Logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		})
	} else {
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		return zapcore.NewJSONEncoder(encoderConfig)
	}
}

// 自定义时间编码器
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	//enc.AppendString(t.Format("2006-01-02 15:04:05"))
	enc.AppendString(t.Format("2006-01-02 15:04:05.000000000"))
}

func getLogWriter(opts *Options) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   opts.logFilename,
		MaxSize:    opts.maxSize, // megabytes
		MaxBackups: opts.maxBackups,
		MaxAge:     opts.maxAge, //days
		Compress:   opts.compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func getLogConsoleWriter(opts *Options) zapcore.WriteSyncer {
	// 日志文件切割归档
	lumberJackLogger := &lumberjack.Logger{
		Filename:   opts.logFilename,
		MaxSize:    opts.maxSize, // megabytes
		MaxBackups: opts.maxBackups,
		MaxAge:     opts.maxAge, //days
		Compress:   opts.compress,
	}

	// 打印到控制台和文件
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
}

func getLogLevel(logLevel string) zapcore.Level {
	level := new(zapcore.Level)
	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		return zap.ErrorLevel
	}

	return *level
}

func (l *Logger) Info(msg string, tags ...zap.Field) {
	l.logger.Info(msg, tags...)
}

func (l *Logger) Error(msg string, tags ...zap.Field) {
	l.logger.Error(msg, tags...)
}

func (l *Logger) Debug(msg string, tags ...zap.Field) {
	l.logger.Debug(msg, tags...)
}

func (l *Logger) Warn(msg string, tags ...zap.Field) {
	l.logger.Warn(msg, tags...)
}

func (l *Logger) Fatal(msg string, tags ...zap.Field) {
	l.logger.Fatal(msg, tags...)
}

func (l *Logger) Panic(msg string, tags ...zap.Field) {
	l.logger.Panic(msg, tags...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.logger.Debug(msg, zap.Any("args", args))
}

func (l *Logger) Infof(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.logger.Info(msg, zap.Any("args", args))
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.logger.Warn(msg, zap.Any("args", args))
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.logger.Error(msg, zap.Any("args", args))
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.logger.Fatal(msg, zap.Any("args", args))
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.logger.Panic(msg, zap.Any("args", args))
}

// NewContext 给指定的context添加字段
func (l *Logger) NewContext(ctx *gin.Context, fields ...zapcore.Field) {
	ctx.Set(LOGGER_KEY, l.WithContext(ctx).logger.With(fields...))
}

// WithContext 从指定的context返回一个zap实例
func (l *Logger) WithContext(ctx *gin.Context) *Logger {
	if ctx == nil {
		return l
	}
	zl, _ := ctx.Get(LOGGER_KEY)
	ctxLogger, ok := zl.(*zap.Logger)
	if ok {
		return &Logger{ctxLogger}
	}
	return l
}

func (l *Logger) Close() {
	_ = l.logger.Sync()
}

func (l *Logger) Sync() {
	_ = l.logger.Sync()
}
