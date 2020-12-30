package crawler

type Dispatcher interface {
  exec(url string) error
}


