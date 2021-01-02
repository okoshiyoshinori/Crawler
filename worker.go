package crawler

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
  url string
  workerLimit int
  jobLimit int
}

//type Job chan string

type Worker struct {
  wg sync.WaitGroup
  sem chan struct{}
  job chan urlStr
  dis Dispatcher
  domain string
  mux sync.RWMutex
  visited map[string] struct{}
}

func NewWorker(d Dispatcher,config Config) *Worker {
  return &Worker{
    sem: make(chan struct{},config.workerLimit),
    job: make(chan urlStr,config.jobLimit),
    dis:d,
  }
}

func (w *Worker) Run() {
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

func (w *Worker) Add() {
  w.wg.Add(1)
}

func (w *Worker) Send(url urlStr) {
  w.job <- url
}

func (w *Worker) Wait() {
  w.wg.Wait()
}
func (w *Worker) Done() {
  w.wg.Done()
}

func (w *Worker) loop(ctx context.Context) {
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
        if data == nil {
          return
        }
        if d,ok := data.([]urlStr); ok {
          for _,s := range d {
            w.Send(s)
          }
        }
      }(job)
    }
  }
  wg.Wait()
}
