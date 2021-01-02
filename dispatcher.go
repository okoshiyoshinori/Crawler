package crawler

type Dispatcher interface {
  exec(url urlStr) (interface{},error)
}


