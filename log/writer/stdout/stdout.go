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
	fmt.Fprint(os.Stdout, fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), string(content)))
	return nil
}

func (n *stdout) Close() error {
	return nil
}
