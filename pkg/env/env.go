package env

import "os"

var (
	_          Env = (*osEnv)(nil)
	_          Env = (*FakeEnv)(nil)
	defaultEnv Env = osEnv{}
)

func Get(key string) string {
	return defaultEnv.Get(key)
}

type Env interface {
	Get(string) string
}

type osEnv struct {
	Env
}

func (osEnv) Get(key string) string {
	return os.Getenv(key)
}

type FakeEnv struct {
	Env
	values map[string]string
}

func SetFakeEnv(kv map[string]string) {
	f := FakeEnv{values: kv}
	defaultEnv = f
}

func (f FakeEnv) Get(key string) string {
	return f.values[key]
}
