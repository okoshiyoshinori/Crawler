package crawler

import (
	"fmt"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Display struct {}

func (d *Display) exec(url string) ([]string,error) {
  time.Sleep(50 * time.Millisecond)
  fmt.Printf("start:%s\n",url)
  doc,err := goquery.NewDocument(url)
  if err != nil {
    return nil,err
  }
  var links []string 
  doc.Find("a").Each(func(_ int,s *goquery.Selection){
    u,_ := s.Attr("href")
    links = append(links,u)
  })
  return links,nil
}

func TestMain(m *testing.M) {
  d := &Display{}
  conf := NewConfig("https://www.omecci.jp",1,3,1000)
  worker := NewWorker(d,conf)
  worker.Run()
  worker.Send(link{url:"https://www.omecci.jp",depth:0})
  worker.Wait()
}


