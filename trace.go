package ethertest

import (
	"encoding/json"
	"fmt"
)

type Contract struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

func (c Contract) lines(from, to int) (int, string) {

	var realFrom, realTo int

	for realFrom = from; realFrom > 0 && c.Source[realFrom] != '\n'; realFrom-- {
	}

	for realTo = to; realTo < len(c.Source) && c.Source[realTo] != '\n'; realTo++ {
	}

	lineNumber := 0
	for i := 0; i < realFrom; i++ {
		if c.Source[i] == '\n' {
			lineNumber++
		}
	}

	return lineNumber, c.Source[realFrom:realTo]

}

type Step struct {
	contractIndex int
	from          int
	to            int
}

func (s Step) MarshalJSON() ([]byte, error) {
	return json.Marshal([]int{s.contractIndex, s.from, s.to})
}

type Trace struct {
	Contracts []Contract `json:"contracts"`
	Steps     []Step     `json:"steps"`
}

func (t *Trace) LastStep() string {
	if len(t.Steps) == 0 {
		return "N/A"
	}
	step := t.Steps[len(t.Steps)-1]
	contract := t.Contracts[step.contractIndex]
	lineNr, source := contract.lines(step.from, step.to)
	return fmt.Sprintf("%s:%d\n%s\n", contract.Name, lineNr+1, source)
}

type tracer struct {
	trace *Trace
}

func newTracer() *tracer {
	return &tracer{
		trace: &Trace{},
	}
}

func (t *tracer) reset() {
	t.trace = &Trace{}
}

func (t *tracer) executed(name, source string, start, end int) {
	idx := -1
	for i, c := range t.trace.Contracts {
		if c.Name == name {
			idx = i
		}
	}
	if idx == -1 {
		idx = len(t.trace.Contracts)
		t.trace.Contracts = append(t.trace.Contracts, Contract{
			Name:   name,
			Source: source,
		})
	}

	newStep := Step{idx, start, end}

	if len(t.trace.Steps) > 0 {
		if t.trace.Steps[len(t.trace.Steps)-1] == newStep {
			return
		}
	}

	t.trace.Steps = append(t.trace.Steps, newStep)

}
