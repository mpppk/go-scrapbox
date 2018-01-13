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

func TestGetPagesListByProjectEndpoint(t *testing.T) {
	opt := &PageListByProjectOptions{
		Skip:  1,
		Limit: 5,
	}

	cases := []struct {
		opt              *PageListByProjectOptions
		project          string
		expectedEndpoint string
	}{
		{
			opt:              opt,
			project:          "test-project",
			expectedEndpoint: fmt.Sprintf("%s/test-project?limit=5&skip=1", apiEndpoint),
		},
		{
			opt:              opt,
			project:          "test te/s;t",
			expectedEndpoint: apiEndpoint + "/test%20te%2Fs%3Bt?limit=5&skip=1",
		},
	}

	for _, c := range cases {
		optStr := fmt.Sprintf("%#v", c.opt)
		actual := getPagesListByProjectEndpoint(c.project, c.opt)
		if actual != c.expectedEndpoint {
			t.Errorf("getPagesListByProjectEndpoint(%q, %q) == %q, want %q",
				c.project, optStr, actual, c.expectedEndpoint)
		}
	}
}

func TestGetEndpoint(t *testing.T) {
	cases := []struct {
		project          string
		title            string
		expectedEndpoint string
	}{
		{
			project:          "test-project",
			title:            "test-title",
			expectedEndpoint: fmt.Sprintf("%s/test-project/test-title", apiEndpoint),
		},
		{
			project:          "test te/s;t",
			title:            "title ti/t;le",
			expectedEndpoint: apiEndpoint + "/test%20te%2Fs%3Bt/title%20ti%2Ft%3Ble",
		},
	}

	for _, c := range cases {
		actual := getPagesGetEndpoint(c.project, c.title)
		if actual != c.expectedEndpoint {
			t.Errorf("getPagesListByProjectEndpoint(%q, %q) == %q, want %q",
				c.project, c.title, actual, c.expectedEndpoint)
		}
	}
}
