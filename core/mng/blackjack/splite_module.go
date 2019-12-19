package blackjack

import (
	"encoding/json"
	"regexp"
	"sync"
)

type RuleModule struct {
	SplitCounter      map[string]int
	DoubleDownCounter map[string]int
	l                 *sync.Mutex
}

func NewRuleModule() *RuleModule {
	return &RuleModule{
		SplitCounter:      make(map[string]int),
		DoubleDownCounter: make(map[string]int),
		l:                 &sync.Mutex{},
	}
}

func ParseRuleModule(source string) *RuleModule {
	ruleModule := &RuleModule{}
	err := json.Unmarshal([]byte(source), ruleModule)
	if err == nil {
		ruleModule.l = &sync.Mutex{}
		return ruleModule
	}
	return NewRuleModule()
}

func clearInternalId(internalId string) string {
	re := regexp.MustCompile(`(.*)(\d)(s{1,})`)
	return re.ReplaceAllString(internalId, "$1$2")
}

func (r *RuleModule) CheckCanSplit(internalId string, rule *Rule) bool {
	r.l.Lock()
	defer r.l.Unlock()
	internalId = clearInternalId(internalId)
	return r.SplitCounter[internalId] <= rule.SplitLimit
}

func (r *RuleModule) IncreaseSplitCounter(internalId string) {
	r.l.Lock()
	defer r.l.Unlock()
	internalId = clearInternalId(internalId)
	r.SplitCounter[internalId]++
}

func (r *RuleModule) DecreaseSplitCounter(internalId string) {
	r.l.Lock()
	defer r.l.Unlock()
	internalId = clearInternalId(internalId)
	r.SplitCounter[internalId]--
	if r.SplitCounter[internalId] <= 0 {
		delete(r.SplitCounter, internalId)
	}
}

func (r *RuleModule) CheckCanDoubleDown(internalId string, rule *Rule) bool {
	r.l.Lock()
	defer r.l.Unlock()
	internalId = clearInternalId(internalId)
	return r.SplitCounter[internalId] <= rule.DoubleDownLimit
}

func (r *RuleModule) IncreaseDoubleDownCounter(internalId string) {
	r.l.Lock()
	defer r.l.Unlock()
	internalId = clearInternalId(internalId)
	r.SplitCounter[internalId]++
}

func (r *RuleModule) DecreaseDoubleDownCounter(internalId string) {
	r.l.Lock()
	defer r.l.Unlock()
	internalId = clearInternalId(internalId)
	r.SplitCounter[internalId]--
	if r.SplitCounter[internalId] <= 0 {
		delete(r.SplitCounter, internalId)
	}
}
