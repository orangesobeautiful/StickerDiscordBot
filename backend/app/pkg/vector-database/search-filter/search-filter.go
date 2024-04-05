package searchfilter

import "golang.org/x/exp/constraints"

type FilterInstance struct {
	t filterType

	content oneOfFilter
}

func (c *FilterInstance) GetType() filterType {
	return c.t
}

func (c *FilterInstance) GetAnd() *FilterAnd {
	return c.content.(*FilterAnd)
}

type filterType int

const (
	FilterTypeAnd filterType = iota + 1
)

type oneOfFilter interface {
	isOneOfFilter()
}

type FilterAnd struct {
	conditions []*ConditionInstance
}

func (c *FilterAnd) GetConditions() []*ConditionInstance {
	return c.conditions
}

func (c *FilterAnd) isOneOfFilter() {}

func AndFilter(conds ...*ConditionInstance) *FilterInstance {
	return &FilterInstance{
		t: FilterTypeAnd,
		content: &FilterAnd{
			conditions: conds,
		},
	}
}

type ConditionInstance struct {
	t conditionType

	content oneOfCondition
}

func (c *ConditionInstance) GetType() conditionType {
	return c.t
}

func (c *ConditionInstance) GetField() *ConditionFieldInstance {
	return c.content.(*ConditionFieldInstance)
}

type conditionType int

const (
	ConditionTypeField conditionType = iota + 1
)

type oneOfCondition interface {
	isOneOfCondition()
}

type ConditionFieldInstance struct {
	key string `validate:"required"`

	match *MatchInstance
}

func (c *ConditionFieldInstance) isOneOfCondition() {}

func (c *ConditionFieldInstance) GetKey() string {
	return c.key
}

func (c *ConditionFieldInstance) GetMatch() *MatchInstance {
	return c.match
}

func ConditionField(key string, match *MatchInstance) *ConditionInstance {
	return &ConditionInstance{
		t: ConditionTypeField,
		content: &ConditionFieldInstance{
			key:   key,
			match: match,
		},
	}
}

type MatchInstance struct {
	t matchType

	matchValue oneOfMatchValue `validate:"required"`
}

func (c *MatchInstance) GetType() matchType {
	return c.t
}

func (c *MatchInstance) GetInIntegers() *MatchInIntegersInstance {
	return c.matchValue.(*MatchInIntegersInstance)
}

type matchType int

const (
	MatchTypeUnknow matchType = iota
	MatchTypeIntegers
)

type oneOfMatchValue interface {
	isOneOfMatchValue()
}

type MatchInIntegersInstance struct {
	value []int64 `validate:"required,dive,min=1"`
}

func (c *MatchInIntegersInstance) isOneOfMatchValue() {}

func (c *MatchInIntegersInstance) GetValue() []int64 {
	return c.value
}

func MatchInIntegers[T constraints.Integer](ints []T) *MatchInstance {
	int64s := make([]int64, len(ints))
	for i, v := range ints {
		int64s[i] = int64(v)
	}

	return &MatchInstance{
		t: MatchTypeIntegers,
		matchValue: &MatchInIntegersInstance{
			value: int64s,
		},
	}
}
