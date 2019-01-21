package ethertest

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"golang.org/x/crypto/sha3"
	. "github.com/logrusorgru/aurora"
	"github.com/tokencard/ethertest/srcmap"
)

type sourceCodeCoverage struct {
	name     string
	ast      solcSource
	source   []byte
	coverage []byte
}

func newSourceCodeCoverage(name string, source []byte, ast solcSource) *sourceCodeCoverage {
	return &sourceCodeCoverage{
		name:     name,
		ast:      ast,
		source:   source,
		coverage: make([]byte, len(source)),
	}
}

func (s *sourceCodeCoverage) percentageCovered() float64 {
	green := 0
	red := 0
	for _, c := range s.coverage {
		switch c {
		case 'R':
			red++
		case 'G':
			green++
		}
	}
	if green == 0 && red == 0 {
		return 100.0
	}
	return float64(green) / (float64(red) + float64(green)) * 100.0
}

func (s *sourceCodeCoverage) paintGreen(from, length int) {
	for i := from; i < from+length; i++ {
		if s.coverage[i] == 'R' {
			s.coverage[i] = 'G'
		}
	}
}

func (s *sourceCodeCoverage) paintRed(from, length int) error {

	for i := from; i < from+length; i++ {
		if i >= len(s.coverage) {
			return fmt.Errorf("combined.json of %s seems to be out of date", s.name)
		}
		s.coverage[i] = 'R'
	}

	return nil
}

func (s *sourceCodeCoverage) Print() {
	for from := 0; from < len(s.source); {
		to := from + 1
		for to < len(s.source) && s.coverage[to] == s.coverage[from] {
			to++
		}
		text := string(s.source[from:to])
		if s.coverage[from] == 'R' {
			fmt.Print(Red(text))
		} else if s.coverage[from] == 'G' {
			fmt.Print(Green(text))
		} else {
			fmt.Print(text)
		}
		from = to
	}

}

type bytecodeWithMapping struct {
	name          string
	tracer        *tracer
	hash          common.Hash
	sourcemap     []srcmap.Entry
	pcToIndex     map[uint64]int
	skipCoverage  []bool
	coverages     []*sourceCodeCoverage
	binary        []byte
	isConstructor bool
}

func (b *bytecodeWithMapping) executed(pc uint64, contractAddress common.Address, contract *vm.Contract) bool {
	if contract.CodeHash != b.hash {
		return false
	}

	if len(b.binary) == 0 {
		return false
	}

	if b.isConstructor {
		c := contract.Code
		if len(c) < len(b.binary) {
			return false
		}
		if bytes.Compare(contract.Code[:len(b.binary)], b.binary) != 0 {
			return false
		}
	}

	idx, f := b.pcToIndex[pc]
	if !f {
		panic(fmt.Errorf("Could not find instruction index for pc %d of contract with hash %s", pc, b.hash.Hex()))
	} else {
		sm := b.sourcemap[idx]
		if !b.skipCoverage[idx] {
			if sm.F >= 0 {
				cov := b.coverages[sm.F]
				cov.paintGreen(sm.S, sm.L)
				b.tracer.executed(cov.name, string(cov.source), sm.S, sm.S+sm.L)
			}
		}
	}
	return true
}

func newBytecodeMapping(t *tracer, name, contractHex string, coverages []*sourceCodeCoverage, smap string, isConstructor bool) (*bytecodeWithMapping, error) {

	contractBinary := common.Hex2Bytes(contractHex)

	hash := common.Hash{}

	if !isConstructor {
		sha := sha3.NewLegacyKeccak256()
		_, err := sha.Write(contractBinary)
		if err != nil {
			return nil, err
		}
		s := sha.Sum(nil)
		copy(hash[:], s)
	}

	ptoi := pcToInstructionMapping(contractBinary)

	sm, err := srcmap.Uncompress(smap)
	if err != nil {
		return nil, err
	}
	skip := make([]bool, len(sm))

	for i, sme := range sm {

		if sme.F >= 0 {
			cov := coverages[sme.F]
			ast := cov.ast.Ast
			as, found := ast.findBySrcPrefix(fmt.Sprintf("%d:%d:", sme.S, sme.L))
			if found {
				switch as.Name {
				case
					"ContractDefinition",
					"IfStatement",
					"FunctionDefinition",
					"Block",
					"PragmaDirective",
					"SourceUnit":
					skip[i] = true
					continue
				}
			}

			for i := sme.S; i < sme.S+sme.L; i++ {
				if sme.F >= 0 {
					err = cov.paintRed(sme.S, sme.L)
					if err != nil {
						return nil, err
					}
				}
			}

		}

	}

	return &bytecodeWithMapping{
		name:          name,
		tracer:        t,
		binary:        contractBinary,
		hash:          hash,
		sourcemap:     sm,
		pcToIndex:     ptoi,
		skipCoverage:  skip,
		coverages:     coverages,
		isConstructor: isConstructor,
	}, nil
}

func newContract(name string, t *tracer, source []byte, ss solcSource, con *solcContract, coverages []*sourceCodeCoverage) (*contract, error) {
	functions := map[[4]byte]*Function{}

	var err error

	ss.Ast.visit(func(n solcASTNode) bool {
		if n.Name == "FunctionDefinition" {
			if !n.Attributes.IsConstructor {

				argTypes := []string{}

				n.visit(func(n solcASTNode) bool {
					if n.Name == "ParameterList" {
						for _, plc := range n.Children {
							argTypes = append(argTypes, plc.Attributes.Type)
						}
						return false
					}
					return true
				})

				n := n.Attributes.Name
				h := sha3.NewLegacyKeccak256()
				fn := fmt.Sprintf("%s(%s)", n, strings.Join(argTypes, ","))
				_, err = h.Write([]byte(fn))
				if err != nil {
					return false
				}
				sum := h.Sum(nil)

				key := [4]byte{}
				copy(key[:], sum)

				functions[key] = &Function{
					name: fn,
				}
				return false
			}
		}
		return true
	})
	if err != nil {
		return nil, err
	}

	// sourceCodeCoverage := newSourceCodeCoverage(name, source, sourceIndex)

	runtimeMapping, err := newBytecodeMapping(t, name, con.BinRuntime, coverages, con.SrcmapRuntime, false)
	if err != nil {
		return nil, err
	}

	constructorMapping, err := newBytecodeMapping(t, name, con.Bin, coverages, con.Srcmap, true)
	if err != nil {
		return nil, err
	}

	return &contract{
		name:      name,
		coverages: coverages,
		mappings:  []*bytecodeWithMapping{runtimeMapping, constructorMapping},
		functions: functions,
		addresses: map[common.Address]struct{}{},
	}, nil
}

type contract struct {
	name      string
	coverages []*sourceCodeCoverage
	mappings  []*bytecodeWithMapping
	functions map[[4]byte]*Function
	addresses map[common.Address]struct{}
}

type Function struct {
	name    string
	gasUsed []uint64
}

func (c *contract) hasAnyGasInformation() bool {
	for _, f := range c.functions {
		if len(f.gasUsed) > 0 {
			return true
		}
	}
	return false
}

func (c *contract) transactionCommited(to common.Address, data []byte, gasUsed uint64) {

	_, found := c.addresses[to]
	if !found {
		return
	}

	if len(data) < 4 {
		return
	}

	prefix := [4]byte{}
	copy(prefix[:], data)
	f, found := c.functions[prefix]
	if found {
		f.gasUsed = append(f.gasUsed, gasUsed)
	}

}

func (c *contract) executed(pc uint64, contractAddress common.Address, contract *vm.Contract) {

	for _, m := range c.mappings {
		matched := m.executed(pc, contractAddress, contract)
		if matched {
			c.addresses[contractAddress] = struct{}{}
		}
	}

}

type solcCombined struct {
	Contracts  map[string]*solcContract `json:"contracts"`
	SourceList []string                 `json:"sourceList"`
	Sources    map[string]solcSource    `json:"sources"`
}

func (s solcCombined) findSourceIndex(name string) int {

	for i, n := range s.SourceList {
		if n == name {
			return i
		}
	}
	return -1

}

type solcContract struct {
	BinRuntime    string  `json:"bin-runtime"`
	SrcmapRuntime string  `json:"srcmap-runtime"`
	Bin           string  `json:"bin"`
	Srcmap        string  `json:"srcmap"`
	Asm           solcAsm `json:"asm"`
}

type solcAsm struct {
	Code []solcAsmCodeEntry `json:".code"`
	Data map[string]solcAsm `json:".data"`
}

type solcAsmCodeEntry struct {
	Begin int    `json:"begin"`
	End   int    `json:"end"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type solcSource struct {
	Ast solcASTNode `json:"AST"`
}

type solcAttributes struct {
	IsConstructor bool   `json:"isConstructor"`
	Name          string `json:"name"`
	Type          string `json:"type"`
}

type solcASTNode struct {
	Attributes solcAttributes `json:"attributes"`
	Children   []solcASTNode  `json:"children"`
	ID         int            `json:"id"`
	Src        string         `json:"src"`
	Name       string         `json:"name"`
}

func (n solcASTNode) findBySrcPrefix(srcPrefix string) (solcASTNode, bool) {
	if strings.HasPrefix(n.Src, srcPrefix) {
		return n, true
	}
	for _, c := range n.Children {
		fc, found := c.findBySrcPrefix(srcPrefix)
		if found {
			return fc, true
		}
	}
	return solcASTNode{}, false
}

func (n solcASTNode) findByName(name string) (solcASTNode, bool) {
	if n.Name == name {
		return n, true
	}
	for _, c := range n.Children {
		fc, found := c.findByName(name)
		if found {
			return fc, true
		}
	}
	return solcASTNode{}, false
}

func (n solcASTNode) allChildren() []solcASTNode {
	r := []solcASTNode{}
	for _, c := range n.Children {
		r = append(r, c)
		cc := c.allChildren()
		r = append(r, cc...)
	}
	return r
}

func (n solcASTNode) visit(f func(n solcASTNode) bool) {
	if f(n) {
		for _, c := range n.Children {
			c.visit(f)
		}
	}
}

func pcToInstructionMapping(b []byte) map[uint64]int {
	mapping := map[uint64]int{}
	cnt := 0
	for i := uint64(0); i < uint64(len(b)); i++ {
		mapping[i] = cnt
		opcode := b[i]
		if opcode >= 0x60 && opcode <= 0x7f {
			i += uint64(opcode) - 0x60 + 1
		}
		cnt++
	}
	return mapping
}
