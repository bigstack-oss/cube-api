package log

const (
	defaultPath       = "/var/log/cube-api.log"
	defaultLevel      = 2
	defaultMaxSize    = 10
	defaultMaxBackups = 3
	defaultMaxAge     = 28
	defaultCompress   = true
)

type Option func(*Options)

type Options struct {
	File     string `json:"file"`
	Level    int    `json:"level"`
	Rotation `json:"rotation"`
}

type Rotation struct {
	Backups  int  `json:"backups"`
	Size     int  `json:"size"`
	TTL      int  `json:"ttl"`
	Compress bool `json:"compress"`
}

func File(file string) Option {
	return func(o *Options) {
		o.File = file
	}
}

func Level(level int) Option {
	return func(o *Options) {
		o.Level = level
	}
}

func Backups(backups int) Option {
	return func(o *Options) {
		o.Rotation.Backups = backups
	}
}

func Size(size int) Option {
	return func(o *Options) {
		o.Rotation.Size = size
	}
}

func TTL(ttl int) Option {
	return func(o *Options) {
		o.Rotation.TTL = ttl
	}
}

func Compress(compress bool) Option {
	return func(o *Options) {
		o.Rotation.Compress = compress
	}
}
