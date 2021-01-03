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
  tmp,err := url.Parse(ustr)
  if err != nil {
    panic("FATAL:not string")
  }
  return &config{
    base:tmp,
    depth:depth,
    workerLimit:wlimit,
    jobLimit:jlimit,
  }
}

type worker struct {
  config *config
  wg sync.WaitGroup
  sem chan struct{}
  job chan urlStr
  dis Dispatcher
  mux sync.RWMutex
  visited map[urlStr] struct{}
}

func NewWorker(d Dispatcher,config config) *worker {
  return &worker{
    sem: make(chan struct{},config.workerLimit),
    job: make(chan urlStr,config.jobLimit),
    dis:d,
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

func (w *worker) Send(url urlStr) {
  w.visitedWrite(url)
  w.job <- url
}

func (w *worker) Wait() {
  w.wg.Wait()
}

func (w *worker) Done() {
  w.wg.Done()
}

func (w *worker) visitedWrite(url urlStr) {
  w.mux.Lock()
  defer w.mux.Unlock()
  if _,ok :=w.visited[url]; ok {
    return
  }
  w.visited[url] = struct{}{}
}

func (w *worker) isVisited(url urlStr) bool {
  w.mux.RLock()
  defer w.mux.RUnlock()
  _,ok := w.visited[url]
  return ok 
}

func (w *worker) isSend(u urlStr) bool {
  if w.isVisited(u) {
    return false
  }
  if !strings.Contains(string(u),w.config.base.Host) {
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
      go func(d urlStr) {
        defer func() { 
          wg.Done() 
          <-w.sem
        }()

        data,err := w.dis.exec(job)
        if err != nil {
          log.Println(err)
          return
        }
        if data == "" {
          return
        }
        if d,ok := data.([]urlStr); ok {
          for _,s := range d {
            if !w.isSend(s) {
              continue
            }
            w.Send(s)
          }
        } else {
          s,_ := data.(urlStr)
          if !w.isSend(s) {
            return
          }
          w.Send(s)
        }
      }(job)
    }
  }
  wg.Wait()
}
