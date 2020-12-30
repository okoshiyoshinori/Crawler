package crawler

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
  workerLimit = 3
  jobLimit = 1000
)

type Job chan string

type Worker struct {
  wg sync.WaitGroup
  sem chan struct{}
  job Job
  dis Dispatcher
}

func NewWorker(d Dispatcher) *Worker {
  return &Worker{
    sem: make(chan struct{},workerLimit),
    job: make(Job,jobLimit),
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
  w.Wait()
}

func (w *Worker) Add() {
  w.wg.Add(1)
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
      go func(d string) {
        defer func() { 
          wg.Done() 
          <-w.sem
        }()
        if err := w.dis.exec(d); err!= nil {
          log.Printf("ERROR:%s",err)
        }
      }(job)
    }
  }
  wg.Wait()
}
