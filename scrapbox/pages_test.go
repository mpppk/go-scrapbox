package scrapbox

import (
	"fmt"
	"strings"
	"testing"
)

func TestGenerateListByProjectQuery(t *testing.T) {
	cases := []struct {
		opt            *PageListByProjectOptions
		expectedQuerys []string
	}{
		{
			opt:            &PageListByProjectOptions{},
			expectedQuerys: []string{""},
		},
		{
			opt: &PageListByProjectOptions{
				Skip:  0,
				Limit: 100,
			},
			expectedQuerys: []string{"limit=100"},
		},
		{
			opt: &PageListByProjectOptions{
				Skip:  1,
				Limit: 0,
			},
			expectedQuerys: []string{"skip=1"},
		},
		{
			opt: &PageListByProjectOptions{
				Skip:  1,
				Limit: 100,
			},
			expectedQuerys: []string{"skip=1", "limit=100"},
		},
	}

	findQuery := func(targetQuery string, queries []string) bool {
		for _, query := range queries {
			if targetQuery == query {
				return true
			}
		}
		return false
	}

	for _, c := range cases {
		queriesStr := generateListByProjectQuery(c.opt)
		optStr := fmt.Sprintf("%#v", c.opt)
		queries := strings.Split(queriesStr, "&")

		for _, expectedQuery := range c.expectedQuerys {
			if !findQuery(expectedQuery, queries) {
				t.Errorf("generateListByProjectQuery(%q) == %q, want %q", optStr, queriesStr, c.expectedQuerys)
			}
		}
	}
}
