package main

import (
	"context"
	"fmt"
	"image/jpeg"
	"os"

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

	page, _, err := client.Pages.Get(context.Background(), "niboshi", "go-scrapbox")
	if err != nil {
		panic(err)
	}
	fmt.Println(*page)

	text, _, err := client.Pages.GetText(context.Background(), "niboshi", "go-scrapbox")
	if err != nil {
		panic(err)
	}

	fmt.Println(text)

	icon, ext, _, err := client.Pages.GetIcon(context.Background(), "niboshi", "niboshi%2Fmpppk")
	if err != nil {
		panic(err)
	}

	file, _ := os.Create("icon." + ext)
	defer file.Close()

	if err := jpeg.Encode(file, *icon, &jpeg.Options{Quality: 100}); err != nil {
		panic(err)
	}
}
