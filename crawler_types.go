package crawler

import "runtime"

const (
  TOPDEPTH = 0
)

var CPUS int = runtime.NumCPU()

type link struct {
  url string
  depth int
}
