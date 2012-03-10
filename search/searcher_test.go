package search
import (
	"github.com/afajl/assert"
	"github.com/afajl/ctrl/remote"
	/*"net/url"*/
	"testing"
)

type testSearcher struct {
	got_id  []string
	got_tag []string
	ret []*remote.Host
}

func (t *testSearcher) Id(s ...string) ([]*remote.Host, error) {
	t.got_id = s
	return nil, nil
}

func (t *testSearcher) Tags(s ...string) ([]*remote.Host, error) {
	t.got_tag = s
	return nil, nil
}

func (t *testSearcher) String() string {
	return "testSearcher"
}


func TestGroupMatches(t *testing.T) {
	a1 := &remote.Host{Id: "a"}
	a2 := &remote.Host{Id: "a"}
	b1 := &remote.Host{Id: "b"}
	b2 := &remote.Host{Id: "b"}

    type grouptest struct {
		matches [][]*remote.Host
		ok bool
	}

	tests := []grouptest{
		{ {{a1, a2}, {b1, b2}}, false}
		{ {{a1, b2}, {b1, a2}}, true}
	}

	for i, test := range tests {
		res, err := groupMatches(test.matches)
		if err != nil {
			t.Fatal(i, err)
		}
		for _, group := range res {
			var group_id string
			for j := 0; j < len(group); j++ {
				if j == 0 {
					group_id = group[0].Id
				}
				assert.Equal(t, group[j].Id, group_id)
			}
		}
	}
}


