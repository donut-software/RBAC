package envvar

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Provider interface {
	Get(key string) (string, error)
}

type Configuration struct {
	provider Provider
}

func Load(filename string) error {
	if err := godotenv.Load(filename); err != nil {
		return fmt.Errorf("loading env var file: %w", err)
	}
	return nil
}

func New(provider Provider) *Configuration {
	return &Configuration{
		provider: provider,
	}
}

// Get returns the value from environment variable `<key>`. When an environment variable `<key>_SECURE` exists
// the provider is used for getting the value.
func (c *Configuration) Get(key string) (string, error) {
	res := os.Getenv(key)
	valSecret := os.Getenv(fmt.Sprintf("%s_SECURE", key))

	if valSecret != "" {
		valSecretRes, err := c.provider.Get(valSecret)
		if err != nil {
			return "", fmt.Errorf("provider get: %w", err)
		}

		res = valSecretRes
	}

	return res, nil
}
