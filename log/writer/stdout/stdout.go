package stdout

import (
	"fmt"
	"os"
	"time"
)

type stdout struct {
}

func New() (*stdout, error) {
	return &stdout{}, nil
}

func (n *stdout) Write(content []byte) error {
	fmt.Fprint(os.Stdout, time.Now(), string(content))
	return nil
}

func (n *stdout) Close() error {
	return nil
}
