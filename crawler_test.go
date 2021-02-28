package crawler

import (
	"fmt"
	"testing"
	"time"
  "sync"
  "strconv"

	"github.com/PuerkitoBio/goquery"
)

type Display struct {}
var count int = 0
var mux sync.RWMutex

func (d *Display) exec(url string) ([]string,error) {
  time.Sleep(1000 * time.Millisecond)
  //fmt.Printf("start:%s\n",url)
  doc,err := goquery.NewDocument(url)
  if err != nil {
    return nil,err
  }
  //get data
  section := doc.Find("div.spot-text")
  section.Each(func(index int,s *goquery.Selection){
    mux.Lock()
    count += 1
    mux.Unlock()
    var data string = strconv.Itoa(count) 
   insec := s.Find("dt.spot-name")
   insec.Each(func(index int,s *goquery.Selection){
     data += "," + s.Text()
   })
   in2 := s.Find("dd.spot-detail-value > span.spot-detail-value-text")
   in2.Each(func(index int,s *goquery.Selection){
     data += "," + s.Text()
   })
   fmt.Println(data)
  })

  var links []string 
  pageSec := doc.Find("ul.paging-section > li")
  pageSec.Each(func(index int,s *goquery.Selection){
    link,_ := s.Find("a").Attr("href")
    links = append(links,link)
  })
  return links,nil
}

func TestMain(m *testing.M) {
  d := &Display{}
  conf := NewConfig("https://www.navitime.co.jp/category/0102001002/",999,CPUS,1000)
  worker := NewWorker(d,conf)
  worker.Run()
  worker.Wait()
}


