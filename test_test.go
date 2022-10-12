package costrace_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/mmyj/costrace"
)

func a(ctx context.Context) {
	sp, ctxx := costrace.StartSpanFromContext(ctx, "a2")
	defer func() {
		sp.FinishWithPrint()
	}()
	time.Sleep(time.Millisecond * 10)
	a1(ctxx)
	a3(ctxx)
}
func a1(ctx context.Context) {
	sp, _ := costrace.StartSpanFromContext(ctx, "a1")
	defer sp.Finish()
	time.Sleep(time.Second * 1)

}

func a3(ctx context.Context) {
	sp, _ := costrace.StartSpanFromContext(ctx, "a3")
	defer sp.Finish()
	time.Sleep(time.Second * 1)
}

func TestParallel(t *testing.T) {
	sp, ctx := costrace.StartSpanFromContext(context.Background(), "trace the function")
	defer func() {
		sp.FinishWithPrint()
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	psp, pctx := costrace.ParallelCtx(ctx)
	go func() {
		defer wg.Done()
		sp1, _ := costrace.StartSpanFromContext(pctx, "func1")
		defer sp1.Finish()
		time.Sleep(time.Second)
	}()

	go func() {
		defer wg.Done()
		sp1, _ := costrace.StartSpanFromContext(pctx, "func2")
		defer sp1.Finish()
		time.Sleep(time.Second)
	}()
	wg.Wait()
	psp.Finish()
	// nest
	sp3, _ := costrace.StartSpanFromContext(ctx, "nest-func1")
	defer sp3.Finish()
	time.Sleep(time.Second)
}

func TestSegmentTrace(t *testing.T) {
	a(context.TODO())
}
