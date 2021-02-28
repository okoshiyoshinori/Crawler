package crawler

type Dispatcher interface {
  Exec(_url string) ([]string,error)
}


