package types

import (
	"strings"

	"gopkg.in/yaml.v2"
)

// String implements fmt.Stringer interface
func (d Symbol) String() string {
	out, _ := yaml.Marshal(d)
	return string(out)
}

// Equal implements equal interface
func (d Symbol) Equal(d1 *Symbol) bool {
	return d.Symbol == d1.Symbol
}

// DenomList is array of Denom
type DenomList []Symbol

// String implements fmt.Stringer interface
func (dl DenomList) String() (out string) {
	for _, d := range dl {
		out += d.String() + "\n"
	}
	return strings.TrimSpace(out)
}
