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

// PagesService handles communication with the page related methods of the Scrapbox API.
type PagesService service

// pageListByProjectResponse is a Scrapbox API response of /api/pages/:projectName
type pageListByProjectResponse struct {
	Skip  int
	Limit int
	Count int
	Pages []*PageWithUserId
}

// User represents a Scrapbox user.
type User struct {
	Id          string
	Name        string
	DisplayName string
	Photo       string
}

// Page represents a Scrapbox page on a project.
// (Some API return only user ID instead of completely user information)
type PageWithUserId struct {
	Page
	UserId string `json:"user"`
}

// Page represents a Scrapbox page on a project.
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

// PageListByRepoOptions specifies the optional parameters to the PagesService.ListByProject method.
type PageListByProjectOptions struct {
	Skip  int
	Limit int
}

// ListByRepo lists the pages for the specified project.
func (s *PagesService) ListByProject(ctx context.Context, project string, opt *PageListByProjectOptions) ([]*Page, *http.Response, error) {
	query := generateListByProjectQuery(opt)
	req, err := s.client.NewRequest("GET", fmt.Sprintf("/api/pages/%s?%s", project, query), nil)
	if err != nil {
		return nil, nil, err
	}

	var pagesRes pageListByProjectResponse
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

// Get a single page.
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

// Get a single text of page.
func (s *PagesService) GetText(ctx context.Context, project, title string) (string, *http.Response, error) {
	buffer, resp, err := s.requestAndDoWithBuffer(ctx, fmt.Sprintf("%s/%s/%s/text", apiEndpoint, project, title))
	return buffer.String(), resp, err
}

// Get a single icon of page.
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
