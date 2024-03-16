package config

import (
	"context"
	"os"
	"test-mkp/exception"

	"github.com/joho/godotenv"
)

type Config interface {
	Get(key string) string
}

type configImpl struct {
}

func (config *configImpl) Get(key string) string {
	return os.Getenv(key)
}

func New(filenames ...string) Config {
	err := godotenv.Load(filenames...)
	exception.PanicIfNeeded(err)
	return &configImpl{}
}

func NewDefaultContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
