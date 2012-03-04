package path

import (
	"github.com/afajl/assert"
	"testing"
)

func TestCleanPathname(t *testing.T) {
	type pathtest struct {
		path   []string
		expect string
	}
	tests := []pathtest{
		// evil stuff
		{[]string{`../foo`}, "..,foo"},
		{[]string{`/foo`}, ",foo"},
		{[]string{`.././foo`}, "..,.,foo"},
		{[]string{`-foo`}, "_foo"},

		// ok
		{[]string{`b-foo`}, "b-foo"},

		// ugly
		{[]string{`!`}, "."},
		{[]string{`|`}, "."},
		{[]string{`*`}, "."},
		{[]string{`?`}, "."},
		{[]string{`'`}, "."},
		{[]string{`"`}, "."},
		{[]string{`<`}, "."},
		{[]string{`>`}, "."},
		{[]string{`\foo`}, ".foo"},
	}
	for _, test := range tests {
		res := CleanPathname(test.path...)
		assert.Equal(t, res, test.expect)
	}
}
