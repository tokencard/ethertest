package ethertest

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	. "github.com/logrusorgru/aurora"
)

func newContract(source []byte, ss solcSource, con *solcContract) *contract {
	cov := make([]byte, len(source))

	type rr struct {
		from int
		to   int
	}

	sha := sha3.NewKeccak256()
	sha.Write([]byte(con.contractBinary()))
	hash := sha.Sum(nil)
	sm := con.runtimeScrmap()
	skip := make([]bool, len(sm))

	for i, sme := range sm {

		as, found := ss.Ast.findBySrcPrefix(fmt.Sprintf("%d:%d:", sme.s, sme.l))
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

		for i := sme.s; i < sme.s+sme.l; i++ {
			cov[i] = 'R'
		}
	}

	return &contract{
		source:       source,
		coverage:     cov,
		pcToIndex:    con.pcToInstructionMapping(),
		sourcemap:    sm,
		hash:         common.BytesToHash(hash),
		skipCoverage: skip,
	}
}

type contract struct {
	source   []byte
	coverage []byte
	// sc           *solcCombined
	pcToIndex    map[uint64]int
	sourcemap    []srcmap
	hash         common.Hash
	skipCoverage []bool
}

func (c *contract) executed(codeHash common.Hash, pc uint64) {
	if codeHash != c.hash {
		return
	}
	idx, f := c.pcToIndex[pc]
	if !f {
		panic(fmt.Errorf("Could not find instruction index for pc %d of contract with hash %s", pc, c.hash.Hex()))
	} else {
		if !c.skipCoverage[idx] {
			sm := c.sourcemap[idx]
			for i := sm.s; i < sm.s+sm.l; i++ {
				if c.coverage[i] == 'R' {
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
	Contracts map[string]*solcContract `json:"contracts"`
	Sources   map[string]solcSource    `json:"sources"`
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
	IsConstructor bool `json:"isConstructor"`
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

func nextMapElement(el string, prev srcmap) srcmap {

	if el == "" {
		return prev
	}

	parts := strings.Split(el, ":")

	r := srcmap{}

	if len(parts) >= 1 {
		if parts[0] == "" {
			r.s = prev.s
		} else {
			var err error
			r.s, err = strconv.Atoi(parts[0])
			if err != nil {
				panic(err)
			}
		}
	}

	if len(parts) >= 2 {
		if parts[1] == "" {
			r.l = prev.l
		} else {
			var err error
			r.l, err = strconv.Atoi(parts[1])
			if err != nil {
				panic(err)
			}
		}
	}
	if len(parts) >= 3 {
		if parts[2] == "" {
			r.f = prev.f
		} else {
			var err error
			r.f, err = strconv.Atoi(parts[2])
			if err != nil {
				panic(err)
			}
		}
	}

	if len(parts) >= 4 {
		if parts[3] == "" {
			r.j = prev.j
		} else {
			r.j = parts[3]
		}
	}

	return r

}

func (sc *solcContract) runtimeScrmap() []srcmap {
	els := strings.Split(sc.SrcmapRuntime, ";")
	prev := srcmap{}
	res := []srcmap{}
	for _, el := range els {
		sme := nextMapElement(el, prev)

		res = append(res, sme)
		prev = sme
	}

	return res
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

type srcmap struct {
	s int
	l int
	f int
	j string
}
