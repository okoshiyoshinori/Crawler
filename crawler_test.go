package crawler

import (
	"fmt"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

type Display struct {}

func (d *Display) exec(url string) ([]string,error) {
  doc,err := goquery.NewDocument(url)
  if err != nil {
    return nil,err
  }
  var links []string 
  doc.Find("a").Each(func(_ int,s *goquery.Selection){
    u,_ := s.Attr("href")
    fmt.Println(u)
    links = append(links,u)
  })
  return links,nil
}

func TestMain(m *testing.M) {
  d := &Display{}
  conf := NewConfig("https://www.omecci.jp",0,3,1000)
  worker := NewWorker(d,conf)
  worker.Run()
  worker.Send(link{url:"https://www.omecci.jp",depth:0})
  worker.Wait()
}


