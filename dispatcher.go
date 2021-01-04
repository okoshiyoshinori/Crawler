package crawler

type Dispatcher interface {
  exec(_url string) ([]string,error)
}


