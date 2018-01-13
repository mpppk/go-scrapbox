package scrapbox

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"net/http"
	"net/url"
)

const apiEndpoint = "/api/pages"

type service struct {
	client *Client
}

type PagesService service

type PageListByProjectResponse struct {
	Skip  int
	Limit int
	Count int
	Pages []*PageWithUserId
}

type User struct {
	Id          string
	Name        string
	DisplayName string
	Photo       string
}

type PageWithUserId struct {
	Page
	UserId string `json:"user"`
}

type Page struct {
	ID           string
	Title        string
	Image        string
	Descriptions []string
	User         *User
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

	var pages []*Page
	for _, pageSummary := range pagesRes.Pages {
		page := &pageSummary.Page
		page.User = &User{Id: pageSummary.UserId}
		pages = append(pages, page)
	}

	return pages, resp, nil
}

func (s *PagesService) Get(ctx context.Context, project, title string) (*Page, *http.Response, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("%s/%s/%s", apiEndpoint, project, title), nil)
	if err != nil {
		return nil, nil, err
	}

	var page Page
	resp, err := s.client.Do(ctx, req, &page)
	if err != nil {
		return nil, resp, err
	}

	return &page, resp, nil
}

func (s *PagesService) GetText(ctx context.Context, project, title string) (string, *http.Response, error) {
	buffer, resp, err := s.requestAndDoWithBuffer(ctx, fmt.Sprintf("%s/%s/%s/text", apiEndpoint, project, title))
	return buffer.String(), resp, err
}

func (s *PagesService) GetIcon(ctx context.Context, project, title string) (*image.Image, string, *http.Response, error) {
	buffer, resp, err := s.requestAndDoWithBuffer(ctx, fmt.Sprintf("%s/%s/%s/icon", apiEndpoint, project, title))
	if err != nil {
		return nil, "", nil, err
	}
	icon, ext, err := image.Decode(buffer)
	return &icon, ext, resp, err
}

func (s *PagesService) requestAndDoWithBuffer(ctx context.Context, endpoint string) (*bytes.Buffer, *http.Response, error) {
	req, err := s.client.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	buffer := new(bytes.Buffer)
	resp, err := s.client.Do(ctx, req, buffer)
	return buffer, resp, err
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
