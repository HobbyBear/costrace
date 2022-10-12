package costrace

import (
	"fmt"
)

const fmtStr = "%s%s (%dms %d%%)\n"

func levelPrint(level int, span *Span, prefix string, ret string) string {
	var (
		lastTabs   string
		noLastTabs string
	)
	noLastTabs = prefix + "├─"
	lastTabs = prefix + "└─"
	span.lock.Lock()
	for i, child := range span.child {
		tabs := noLastTabs
		if i == len(span.child)-1 {
			tabs = lastTabs
		}
		childCostMs := child.cost().Milliseconds()
		fatherCostMs := span.cost().Milliseconds()
		radio := int64(0)
		if fatherCostMs > 0 {
			radio = childCostMs * 100 / fatherCostMs
		}
		ret += fmt.Sprintf(fmtStr, tabs, child.title, childCostMs, radio)
		if len(child.child) > 0 {
			if i == len(span.child)-1 {
				ret = levelPrint(level+1, child, prefix+"  ", ret)
			} else {
				ret = levelPrint(level+1, child, prefix+"│  ", ret)
			}
		}
	}
	span.lock.Unlock()
	return ret
}

// ToString format a tracer to string
func (s *Span) String() (ret string) {
	ret += fmt.Sprintf(fmtStr, "", s.title, s.cost().Milliseconds(), 100)
	ret = levelPrint(0, s, "", ret)
	return
}
