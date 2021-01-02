package crawler

import (
	"fmt"
	"testing"
)

type Display struct {}

func (d *Display) exec(url string) (interface{},error) {
  fmt.Println("url:",url)
  return nil, nil
}
/*
func TestMain(m *testing.M) {
  dis := &Display{}
  w := NewWorker(dis)
  w.Run()

  for i:=0;i<10; i++ {
    url := fmt.Sprintf("http://google.com/%d",i)
    w.Send(url)
  }
  w.Wait()
}
*/
