package log

// Option is config option.
type Option func(*Options)

type Options struct {
	// logger options
	mode     string // dev or prod
	logLevel string // debug, info, warn, error, panic, panic, fatal
	encoding string // console or json

	// lumberjack options
	logFilename string
	maxSize     int
	maxBackups  int
	maxAge      int
	compress    bool
}

// DefaultOptions .
func DefaultOptions() *Options {
	return &Options{
		mode:     "dev",
		logLevel: "info",
		encoding: "console",

		logFilename: "logs.log",
		maxSize:     500, // megabytes
		maxBackups:  3,
		maxAge:      28, //days
		compress:    true,
	}
}

func Apply(opts ...Option) *Options {
	options := DefaultOptions()
	for _, o := range opts {
		o(options)
	}
	return options
}

func WithMode(mode string) Option {
	return func(o *Options) {
		o.mode = mode
	}
}

func WithLogLevel(level string) Option {
	return func(o *Options) {
		o.logLevel = level
	}
}

func WithEncoding(encoding string) Option {
	return func(o *Options) {
		o.encoding = encoding
	}
}

func WithFilename(filename string) Option {
	return func(o *Options) {
		o.logFilename = filename
	}
}

func WithMaxSize(maxSize int) Option {
	return func(o *Options) {
		o.maxSize = maxSize
	}
}

func WithMaxBackups(maxBackups int) Option {
	return func(o *Options) {
		o.maxBackups = maxBackups
	}
}

func WithMaxAge(maxAge int) Option {
	return func(o *Options) {
		o.maxAge = maxAge
	}
}

func WithCompress(compress bool) Option {
	return func(o *Options) {
		o.compress = compress
	}
}
