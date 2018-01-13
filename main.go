package main

import (
	"context"
	"fmt"

	"github.com/mpppk/go-scrapbox/scrapbox"
)

func main() {
	client := scrapbox.NewClient(nil)
	pages, _, err := client.Pages.ListByProject(context.Background(), "niboshi",
		&scrapbox.PageListByProjectOptions{Skip: 1, Limit: 5})
	if err != nil {
		panic(err)
	}

	for _, page := range pages {
		fmt.Println(*page)
	}
}
