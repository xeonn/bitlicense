package db

// TODO: IMPLEMENT COUCHDB DATABASE FOR LICENSE MANAGEMENT

type DbConfig struct {
	Uri  string
	User string
	Pass string
}

func NewDbConfig(uri, user, pass string) *DbConfig {
	return &DbConfig{
		Uri:  uri,
		User: user,
		Pass: pass,
	}
}
