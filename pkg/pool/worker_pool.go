package pool

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"sync"
	"time"
)

type Task func(ctx context.Context) error

type WorkerPool struct {
	workerCount int       //goroutines number
	queueSize   int       //task queue's size
	taskQueue   chan Task //task queue
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	logger      *zap.Logger
}

type Option func(*WorkerPool)

func WithWorkerCount(count int) Option {
	return func(wp *WorkerPool) {
		if count > 0 {
			wp.workerCount = count
		}
	}
}

func WithQueueSize(size int) Option {
	return func(wp *WorkerPool) {
		if size > 0 {
			wp.queueSize = size
		}
	}
}

func NewWorkerPool(logger *zap.Logger, opts ...Option) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &WorkerPool{
		workerCount: 10,
		queueSize:   1000,
		ctx:         ctx,
		cancel:      cancel,
		logger:      logger,
	}

	for _, opt := range opts {
		opt(pool)
	}

	pool.taskQueue = make(chan Task, pool.queueSize)
	pool.startWorkers()

	return pool
}

func (p *WorkerPool) startWorkers() {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go func(workerID int) {
			defer p.wg.Done()
			p.workerLoop(workerID)
		}(i)
	}
	p.logger.Info("Worker pool started",
		zap.Int("workerCount", p.workerCount),
		zap.Int("queueSize", p.queueSize))
}

func (p *WorkerPool) workerLoop(workerID int) {
	p.logger.Debug("Worker pool started", zap.Int("workerId", workerID))
	defer p.logger.Debug("Worker pool exited", zap.Int("workerId", workerID))

	for {
		select {
		case <-p.ctx.Done():
			return
		case task, ok := <-p.taskQueue:
			if !ok {
				return
			}
			p.executeTask(workerID, task)
		}
	}
}

func (p *WorkerPool) executeTask(workerID int, task Task) {
	startTime := time.Now()
	defer func() {
		if r := recover(); r != nil {
			p.logger.Error("Worker task panic",
				zap.Int("workerId", workerID),
				zap.Any("panic", r),
				zap.Stack("stack"),
				zap.Duration("time", time.Since(startTime)))
		}
	}()

	if err := task(p.ctx); err != nil {
		p.logger.Warn("Worker task execution failed",
			zap.Int("workerId", workerID),
			zap.Error(err),
			zap.Stack("stack"),
			zap.Duration("time", time.Since(startTime)),
		)
	} else {
		p.logger.Debug("Worker task execution succeeded",
			zap.Int("workerId", workerID),
			zap.Duration("time", time.Since(startTime)),
		)
	}
}

func (p *WorkerPool) Submit(task Task) error {
	select {
	case p.taskQueue <- task:
		return p.ctx.Err()
	case p.taskQueue <- task:
		return nil
	default:
		return ErrQueueFull
	}
}

func (p *WorkerPool) Shutdown() {
	close(p.taskQueue)
	p.cancel()
	p.wg.Wait()
	p.logger.Info("Worker pool shutdown completed")
}

// for error
var (
	ErrQueueFull = errorf("task queue is full")
)

func errorf(format string, v ...interface{}) error {
	return &poolError{format: format, args: v}
}

type poolError struct {
	format string
	args   []interface{}
}

func (e *poolError) Error() string {
	return fmt.Sprintf(e.format, e.args...)
}
