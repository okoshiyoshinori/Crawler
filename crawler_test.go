package crawler

import (
	"fmt"
	"testing"
)

type Display struct {}

func (d *Display) exec(url string) (interface{},error) {
  fmt.Println("url:",url)
  return "", nil
}

func TestMain(m *testing.M) {
}

