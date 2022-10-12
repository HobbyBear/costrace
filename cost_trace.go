package costrace

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"
)

type (
	costRaceKey int
)

var (
	key costRaceKey
)

type Span struct {
	startTime time.Time
	endTime   time.Time
	title     string
	lock      sync.Mutex
	child     []*Span
}

func (s *Span) cost() time.Duration {
	return s.endTime.Sub(s.startTime)
}

func StartSpan(title string) (*Span, context.Context) {
	this := &Span{title: title, startTime: time.Now()}
	return this, context.WithValue(context.Background(), key, this)
}

func StartSpanFromContext(ctx context.Context, title string) (*Span, context.Context) {
	this := &Span{title: title, startTime: time.Now()}
	father, ok := ctx.Value(key).(*Span)
	if !ok {
		return this, context.WithValue(ctx, key, this)
	}
	father.lock.Lock()
	father.child = append(father.child, this)
	father.lock.Unlock()

	return this, context.WithValue(ctx, key, this)
}

// Done stop trace timer
func (s *Span) Finish() {
	s.endTime = time.Now()
}

func (s *Span) FinishWithPrint() {
	s.endTime = time.Now()
	s.Print()
}

func (s *Span) Print() {
	m := make(map[string]string)
	m["desc"] = s.String()
	desc, _ := json.Marshal(m)
	log.Println(string(desc))
}
