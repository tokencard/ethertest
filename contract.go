package ethertest

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	. "github.com/logrusorgru/aurora"
	"github.com/tokencard/ethertest/srcmap"
)

func newContract(name string, source []byte, ss solcSource, con *solcContract, sourceIndex int) (*contract, error) {
	cov := make([]byte, len(source))

	type rr struct {
		from int
		to   int
	}

	sha := sha3.NewKeccak256()
	sha.Write([]byte(con.contractBinary()))
	hash := sha.Sum(nil)
	sm, err := con.runtimeScrmap()
	if err != nil {
		return nil, err
	}
	skip := make([]bool, len(sm))

	for i, sme := range sm {

		as, found := ss.Ast.findBySrcPrefix(fmt.Sprintf("%d:%d:", sme.S, sme.L))
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
			if sme.F == sourceIndex {
				if i >= len(cov) {
					return nil, fmt.Errorf("combined.json of %s seems to be out of date", name)
				}
				cov[i] = 'R'
			}
		}
	}

	functions := map[[4]byte]*Function{}

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
				h := sha3.NewKeccak256()
				fn := fmt.Sprintf("%s(%s)", n, strings.Join(argTypes, ","))
				_, err = h.Write([]byte(fn))
				if err != nil {
					panic(err)
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

	return &contract{
		name:         name,
		source:       source,
		coverage:     cov,
		pcToIndex:    con.pcToInstructionMapping(),
		sourcemap:    sm,
		hash:         common.BytesToHash(hash),
		skipCoverage: skip,
		sourceIndex:  sourceIndex,
		functions:    functions,
		addresses:    map[common.Address]struct{}{},
	}, nil
}

type contract struct {
	name     string
	source   []byte
	coverage []byte
	// sc           *solcCombined
	pcToIndex    map[uint64]int
	sourcemap    []srcmap.Entry
	hash         common.Hash
	skipCoverage []bool
	sourceIndex  int
	functions    map[[4]byte]*Function
	addresses    map[common.Address]struct{}
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

func (c *contract) executed(codeHash common.Hash, pc uint64, contractAddress common.Address) {
	if codeHash != c.hash {
		return
	}

	c.addresses[contractAddress] = struct{}{}

	idx, f := c.pcToIndex[pc]
	if !f {
		panic(fmt.Errorf("Could not find instruction index for pc %d of contract with hash %s", pc, c.hash.Hex()))
	} else {
		if !c.skipCoverage[idx] {
			sm := c.sourcemap[idx]
			for i := sm.S; i < sm.S+sm.L; i++ {
				if sm.F == c.sourceIndex && c.coverage[i] == 'R' {
					c.coverage[i] = 'G'
				}
			}
		}
	}

}

func (c *contract) Print() {
	for from := 0; from < len(c.source); {
		to := from + 1
		for to < len(c.source) && c.coverage[to] == c.coverage[from] {
			to++
		}
		text := string(c.source[from:to])
		if c.coverage[from] == 'R' {
			fmt.Print(Red(text))
		} else if c.coverage[from] == 'G' {
			fmt.Print(Green(text))
		} else {
			fmt.Print(text)
		}
		from = to
	}

}

func (c *contract) percentageCovered() float64 {
	green := 0
	red := 0
	for _, c := range c.coverage {
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

func (sc *solcContract) runtimeScrmap() (srcmap.Map, error) {
	return srcmap.Uncompress(sc.SrcmapRuntime)
}

func (sc *solcContract) contractBinary() []byte {
	return common.Hex2Bytes(sc.BinRuntime)
}

func (sc *solcContract) pcToInstructionMapping() map[uint64]int {
	mapping := map[uint64]int{}
	b := sc.contractBinary()
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

func (sc *solcContract) countInstructions() int {
	b := sc.contractBinary()
	cnt := 0
	for i := 0; i < len(b); i++ {
		cnt++
		opcode := b[i]
		if opcode <= 0x7f && opcode >= 0x60 {
			i += int(opcode) - 0x60 + 1
		}
	}
	return cnt
}
