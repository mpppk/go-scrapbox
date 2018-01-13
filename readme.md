# go-scrapbox
golang Client for Scrapbox

[![GoDoc](https://godoc.org/github.com/mpppk/go-scrapbox?status.svg)](http://godoc.org/github.com/mpppk/go-scrapbox)
[![CircleCI](https://circleci.com/gh/mpppk/go-scrapbox.svg?style=svg)](https://circleci.com/gh/mpppk/go-scrapbox)
[![codebeat badge](https://codebeat.co/badges/8cdde201-641d-4055-90f6-228a867a51b3)](https://codebeat.co/projects/github-com-mpppk-go-scrapbox-initial)

# Usage

## import

```golang
import github.com/mpppk/go-scrapbox/scrapbox
```

## List pages

```golang
client := scrapbox.NewClient(nil)
pages, httpResponse, err := client.Pages.ListByProject(context.Background(), "project name",
    &scrapbox.PageListByProjectOptions{Skip: 1, Limit: 5})
fmt.Println(pages[0].Title) // => Display first page title
```

## Get page

```golang
client := scrapbox.NewClient(nil)
page, httpResponse, err := client.Pages.Get(context.Background(), "project name", "page title")
fmt.Println(page.Title) // => Display page title
```

### Get page as plain text

```golang
client := scrapbox.NewClient(nil)
text, httpResponse, err := client.Pages.GetText(context.Background(), "project name", "page title")
fmt.Println(text) // => Display page as plain text
```

### Get page icon

```golang
client := scrapbox.NewClient(nil)
icon, ext, httpResponse, err := client.Pages.GetIcon(context.Background(), "project name", "page title")
file, err := os.Create("icon." + ext)
jpeg.Encode(file, *icon, &jpeg.Options{Quality: 100}) // => Save icon as "icon.jpg/png/gif"
```
