package scrapbox

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type service struct {
	client *Client
}

type PagesService service

type PageListByProjectResponse struct {
	Skip  int
	Limit int
	Count int
	Pages []*Page
}

type Page struct {
	ID           string
	Title        string
	Image        string
	Descriptions []string
	User         string
	Pin          int64
	Views        int
	Point        int
	Linked       int
	CommitID     string
	Created      int
	Updated      int
	Accessed     int
}

type PageListByProjectOptions struct {
	Skip  int
	Limit int
}

type Icon struct {
}

//func (s *PagesService) Create(ctx context.Context, project string) (*Page, *http.Response, error) {
//
//}

func (s *PagesService) ListByProject(ctx context.Context, project string, opt *PageListByProjectOptions) ([]*Page, *http.Response, error) {
	query := generateListByProjectQuery(opt)
	req, err := s.client.NewRequest("GET", fmt.Sprintf("/api/pages/%s?%s", project, query), nil)
	if err != nil {
		return nil, nil, err
	}

	var pagesRes PageListByProjectResponse
	resp, err := s.client.Do(ctx, req, &pagesRes)
	if err != nil {
		return nil, resp, err
	}

	return pagesRes.Pages, resp, nil
}

func generateListByProjectQuery(opt *PageListByProjectOptions) string {
	values := url.Values{}
	if opt != nil {
		if opt.Skip != 0 {
			values.Add("skip", fmt.Sprint(opt.Skip))
		}

		if opt.Limit != 0 {
			values.Add("limit", fmt.Sprint(opt.Limit))
		}
	}
	return values.Encode()
}

//func (s *PagesService) Get(ctx context.Context, project, title string) (*Page, *http.Response, error) {
//
//}
//
//func (s *PagesService) GetIcon(ctx context.Context, project, title string) (*Icon, *http.Response, error) {
//
//}
