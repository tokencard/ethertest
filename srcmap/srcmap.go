package srcmap

import (
	"fmt"
	"strconv"
	"strings"
)

type Entry struct {
	S int
	L int
	F int
	J string
}

func (e Entry) String() string {
	return fmt.Sprintf("%d:%d:%d:%s", e.S, e.L, e.F, e.J)
}

type Map []Entry

func (m Map) String() string {
	parts := make([]string, len(m))
	for i, e := range m {
		parts[i] = e.String()
	}
	return strings.Join(parts, ";")
}

func nextMapEntry(el string, prev Entry) (Entry, error) {

	if el == "" {
		return prev, nil
	}

	parts := strings.Split(el, ":")

	r := prev

	if len(parts) >= 1 {
		if parts[0] != "" {
			var err error
			r.S, err = strconv.Atoi(parts[0])
			if err != nil {
				panic(err)
			}
		}
	}

	if len(parts) >= 2 {
		if parts[1] != "" {
			var err error
			r.L, err = strconv.Atoi(parts[1])
			if err != nil {
				return Entry{}, err
			}
		}
	}
	if len(parts) >= 3 {
		if parts[2] != "" {
			var err error
			r.F, err = strconv.Atoi(parts[2])
			if err != nil {
				return Entry{}, err
			}
		}
	}

	if len(parts) >= 4 {
		if parts[3] != "" {
			r.J = parts[3]
		}
	}

	return r, nil

}

// Uncompress will convert srcmap string into slice of srcmap entries.
func Uncompress(compressed string) (Map, error) {
	els := strings.Split(compressed, ";")
	prev := Entry{}
	res := []Entry{}
	for _, el := range els {
		sme, err := nextMapEntry(el, prev)
		if err != nil {
			return nil, err
		}

		res = append(res, sme)
		prev = sme

	}
	return res, nil
}
