package crawler

import (
	"context"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
  "strings"
)

type config struct {
  base *url.URL 
  depth int
  workerLimit int
  jobLimit int
}

func NewConfig(ustr string,depth int,wlimit int,jlimit int) *config {
  p,err := url.Parse(ustr)
  if err != nil {
    panic("FATAL:Can't parse this url")
  }
  return &config{
    base:p,
    depth:depth,
    workerLimit:wlimit,
    jobLimit:jlimit,
  }
}

type worker struct {
  config *config
  wg sync.WaitGroup
  sem chan struct{}
  job chan link
  dis Dispatcher
  mux sync.RWMutex
  visited map[string] struct{}
}

func NewWorker(d Dispatcher,c *config) *worker {
  return &worker{
    config: c,
    sem: make(chan struct{},c.workerLimit),
    job: make(chan link,c.jobLimit),
    dis:d,
    visited:make(map[string]struct{}),
  }
}

func (w *worker) Run() {
  w.Add()
  ctx := context.Background()
  interrupt,cancel := context.WithCancel(ctx)
  sig := make(chan os.Signal,1)
  signal.Notify(sig,syscall.SIGINT,syscall.SIGKILL,syscall.SIGTERM)
  go func() {
    <-sig
    cancel()
    defer w.Done()
  }()
  go w.loop(interrupt)
}

func (w *worker) Add() {
  w.wg.Add(1)
}

func (w *worker) Send(l link) {
  w.visitedWrite(l)
  w.job <- l
}

func (w *worker) Wait() {
  w.wg.Wait()
}

func (w *worker) Done() {
  w.wg.Done()
}

func (w *worker) getLink(l link) ([]link,error) {
  next_depth := l.depth + 1
  //絶対パスに変換
  u,_ := ToAbsUrl(w.config.base,l.url)
  links,err := w.dis.exec(u)
  if err != nil {
    return nil,err
  }

  if links == nil {
    return nil,nil
  }

  var t []link
  for _,s := range links {
    tmp := link{
      depth: next_depth,
      url:s,
    }
    t = append(t,tmp)
  }
  return t,nil
}

func (w *worker) visitedWrite(l link) {
  w.mux.Lock()
  defer w.mux.Unlock()
  if _,ok :=w.visited[l.url]; ok {
    return
  }
  w.visited[l.url] = struct{}{}
}

func (w *worker) isVisited(l link) bool {
  w.mux.RLock()
  defer w.mux.RUnlock()
  _,ok := w.visited[l.url]
  return ok 
}

func (w *worker) isSend(l link) bool {
  if w.isVisited(l) {
    return false
  }
  if !strings.Contains(l.url,w.config.base.Host) {
    return false
  }
  if l.depth > w.config.depth {
    return false
  }
  return true
}

func (w *worker) loop(ctx context.Context) {
  var wg sync.WaitGroup
loop:
  for {
    select {
    case <-ctx.Done():
      break loop
    case job := <-w.job:
      w.sem <- struct{}{}
      wg.Add(1)
      go func(d link) {
        defer func() { 
          wg.Done() 
          <-w.sem
        }()
        data,err := w.getLink(job)
        if err != nil {
          log.Println(err)
          return
        }
        if data == nil {
          return
        }
        for _,s := range data {
          if !w.isSend(s) {
             continue
          }
            w.Send(s)
        }
      }(job)
    }
  }
  wg.Wait()
}
