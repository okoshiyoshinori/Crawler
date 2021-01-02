package crawler

type Dispatcher interface {
  exec(url string) (interface{},error)
}


