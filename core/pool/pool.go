// Package pool 提供带并发限制的协程池。
//
// 基础用法（仅限制并发数）：
//
//	p := pool.New(10)
//	for _, item := range list {
//	    item := item
//	    p.Go(func() error {
//	        return process(item)
//	    })
//	}
//	if err := p.Wait(); err != nil {
//	    // 处理合并后的错误
//	}
//
// 进阶用法（context 取消 + 自定义 panic 处理）：
//
//	p := pool.New(10,
//	    pool.WithContext(ctx),
//	    pool.WithPanicHandler(func(r any) {
//	        log.Printf("panic: %v", r)
//	    }),
//	)
package pool

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Pool 限制并发协程数，收集任务错误，支持 context 取消。
type Pool struct {
	sem          chan struct{}
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
	panicHandler func(r any)

	mu   sync.Mutex
	errs []error
}

// Option 用于配置 Pool。
type Option func(*Pool)

// WithContext 使 Pool 感知外部 context。ctx 取消时：
//   - Go 立即返回 false（任务不提交）
//   - TryGo 立即返回 false
//
// 调用 Pool.Cancel() 同样会取消该 context。
func WithContext(ctx context.Context) Option {
	return func(p *Pool) {
		p.ctx, p.cancel = context.WithCancel(ctx)
	}
}

// WithPanicHandler 设置自定义 panic 处理函数。
// r 为 recover() 捕获的值。未设置时 panic 会被静默恢复。
func WithPanicHandler(h func(r any)) Option {
	return func(p *Pool) {
		p.panicHandler = h
	}
}

// New 创建最多允许 n 个协程同时运行的 Pool。n < 1 时 panic。
func New(n int, opts ...Option) *Pool {
	if n < 1 {
		panic("pool: n must be >= 1")
	}
	ctx, cancel := context.WithCancel(context.Background())
	p := &Pool{
		sem:    make(chan struct{}, n),
		ctx:    ctx,
		cancel: cancel,
	}
	for _, o := range opts {
		o(p)
	}
	return p
}

// Go 提交任务 f。若并发槽位已满则阻塞等待；若 context 已取消则立即返回 false。
// f 返回的非 nil error 会被收集，可通过 Wait 或 Errors 获取。
func (p *Pool) Go(f func() error) bool {
	select {
	case <-p.ctx.Done():
		return false
	case p.sem <- struct{}{}:
	}
	p.wg.Add(1)
	go p.run(f)
	return true
}

// TryGo 非阻塞提交任务 f。有空闲槽位且 context 未取消时返回 true，否则立即返回 false。
func (p *Pool) TryGo(f func() error) bool {
	select {
	case <-p.ctx.Done():
		return false
	case p.sem <- struct{}{}:
		p.wg.Add(1)
		go p.run(f)
		return true
	default:
		return false
	}
}

// Cancel 取消 Pool 的 context，后续 Go / TryGo 调用将返回 false。
// 已在运行的协程不会被强制中断，如需协作退出请在任务内部检查 ctx.Done()。
func (p *Pool) Cancel() {
	p.cancel()
}

// Wait 阻塞直到所有已提交的协程执行完毕，返回合并后的错误（含 panic 转化的错误）。
func (p *Pool) Wait() error {
	p.wg.Wait()
	p.mu.Lock()
	defer p.mu.Unlock()
	return errors.Join(p.errs...)
}

// Errors 返回当前已收集的错误列表，不阻塞。
// 建议在 Wait 之后调用以获得完整结果。
func (p *Pool) Errors() []error {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]error, len(p.errs))
	copy(out, p.errs)
	return out
}

func (p *Pool) run(f func() error) {
	defer func() {
		<-p.sem
		// recover 和 addErr 必须在 wg.Done 之前执行，
		// 确保 Wait 返回时所有错误已收集完毕。
		if r := recover(); r != nil {
			if p.panicHandler != nil {
				p.panicHandler(r)
			}
			p.addErr(fmt.Errorf("panic: %v", r))
		}
		p.wg.Done()
	}()
	if err := f(); err != nil {
		p.addErr(err)
	}
}

func (p *Pool) addErr(err error) {
	p.mu.Lock()
	p.errs = append(p.errs, err)
	p.mu.Unlock()
}
