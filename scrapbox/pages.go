package scrapbox

import (
	"context"
	"fmt"
	"net/http"
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
}

type Icon struct {
}

//func (s *PagesService) Create(ctx context.Context, project string) (*Page, *http.Response, error) {
//
//}

func (s *PagesService) ListByProject(ctx context.Context, project string, opt *PageListByProjectOptions) ([]*Page, *http.Response, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("/api/pages/%s", project), nil) // TODO fix urlStr
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

//func (s *PagesService) Get(ctx context.Context, project, title string) (*Page, *http.Response, error) {
//
//}
//
//func (s *PagesService) GetIcon(ctx context.Context, project, title string) (*Icon, *http.Response, error) {
//
//}
