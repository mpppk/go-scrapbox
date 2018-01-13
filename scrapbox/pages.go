package scrapbox

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
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
	endpoint := getPagesListByProjectEndpoint(project, opt)
	req, err := s.client.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create http request for PagesService.ListByProject")
	}

	var pagesRes pageListByProjectResponse
	resp, err := s.client.Do(ctx, req, &pagesRes)
	if err != nil {
		errMsg := fmt.Sprintf("failed to request to %s for PagesService.ListByProject", req.URL.String())
		return nil, resp, errors.Wrap(err, errMsg)
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
	req, err := s.client.NewRequest("GET", getPagesGetEndpoint(project, title), nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create http request for PagesService.Get")
	}

	var page Page
	resp, err := s.client.Do(ctx, req, &page)
	if err != nil {
		errMsg := fmt.Sprintf("failed to request to %s for PagesService.Get", req.URL.String())
		return nil, resp, errors.Wrap(err, errMsg)
	}

	return &page, resp, nil
}

// Get a single text of page.
func (s *PagesService) GetText(ctx context.Context, project, title string) (string, *http.Response, error) {
	endpoint := getPagesGetTextEndpoint(project, title)
	buffer, resp, err := s.requestAndDoWithBuffer(ctx, endpoint)
	if err != nil {
		errMsg := fmt.Sprintf("failed to request for PagesService.GetText. endpoint: %s", endpoint)
		return buffer.String(), resp, errors.Wrap(err, errMsg)
	}
	return buffer.String(), resp, nil
}

// Get a single icon of page.
func (s *PagesService) GetIcon(ctx context.Context, project, title string) (*image.Image, string, *http.Response, error) {
	endpoint := getPagesGetIconEndpoint(project, title)
	buffer, resp, err := s.requestAndDoWithBuffer(ctx, endpoint)
	if err != nil {
		errMsg := fmt.Sprintf("failed to request for PagesService.GetIcon. endpoint: %s", endpoint)
		return nil, "", nil, errors.Wrap(err, errMsg)
	}

	icon, ext, err := image.Decode(buffer)
	if err != nil {
		return &icon, ext, resp, errors.Wrap(err, "failed to decode icon")
	}
	return &icon, ext, resp, nil
}

func (s *PagesService) requestAndDoWithBuffer(ctx context.Context, endpoint string) (*bytes.Buffer, *http.Response, error) {
	req, err := s.client.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	buffer := new(bytes.Buffer)
	resp, err := s.client.Do(ctx, req, buffer)
	if err != nil {
		errMsg := fmt.Sprintf("failed to request to %s for PagesService", req.URL.String())
		return buffer, resp, errors.Wrap(err, errMsg)
	}

	return buffer, resp, nil
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

func getPagesListByProjectEndpoint(project string, opt *PageListByProjectOptions) string {
	query := generateListByProjectQuery(opt)
	escapedProject := url.PathEscape(project)
	return fmt.Sprintf("%s/%s?%s", apiEndpoint, escapedProject, query)
}

func getPagesGetEndpoint(project, title string) string {
	escapedProject := url.PathEscape(project)
	escapedTitle := url.PathEscape(title)
	return fmt.Sprintf("%s/%s/%s", apiEndpoint, escapedProject, escapedTitle)
}

func getPagesGetTextEndpoint(project, title string) string {
	return getPagesGetEndpoint(project, title) + "/text"
}

func getPagesGetIconEndpoint(project, title string) string {
	return getPagesGetEndpoint(project, title) + "/icon"
}
