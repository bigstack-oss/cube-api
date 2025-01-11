package mongo

var (
	Opts *Options
)

type Option func(*Options)

type Options struct {
	Uri string `validate:"required"`
	Auth
	ReplicaSet string
	Connect    string

	Database    string
	Collection  string
	Databases   map[string]string
	Collections map[string]string
}

type Auth struct {
	Enable   bool
	Source   string
	Username string
	Password string
}

func Uri(uri string) Option {
	return func(o *Options) {
		o.Uri = uri
	}
}

func ReplicaSet(replicaSet string) Option {
	return func(o *Options) {
		o.ReplicaSet = replicaSet
	}
}

func Connect(connect string) Option {
	return func(o *Options) {
		o.Connect = connect
	}
}

func Database(database string) Option {
	return func(o *Options) {
		o.Database = database
	}
}

func Collection(collection string) Option {
	return func(o *Options) {
		o.Collection = collection
	}
}

func Databases(databases map[string]string) Option {
	return func(o *Options) {
		o.Databases = databases
	}
}

func Collections(collections map[string]string) Option {
	return func(o *Options) {
		o.Collections = collections
	}
}

func AuthEnable(enable bool) Option {
	return func(o *Options) {
		o.Auth.Enable = enable
	}
}

func AuthSource(source string) Option {
	return func(o *Options) {
		o.Auth.Source = source
	}
}

func AuthUsername(username string) Option {
	return func(o *Options) {
		o.Auth.Username = username
	}
}

func AuthPassword(password string) Option {
	return func(o *Options) {
		o.Auth.Password = password
	}
}
