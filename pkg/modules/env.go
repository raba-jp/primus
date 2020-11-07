package modules

import "os"

var (
	_          Env = (*osEnv)(nil)
	_          Env = (*FakeEnv)(nil)
	defaultEnv Env = osEnv{}
)

func Getenv(key string) string {
	return defaultEnv.Getenv(key)
}

type Env interface {
	Getenv(string) string
}

type osEnv struct {
	Env
}

func (osEnv) Getenv(key string) string {
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

func (f FakeEnv) Getenv(key string) string {
	return f.values[key]
}
