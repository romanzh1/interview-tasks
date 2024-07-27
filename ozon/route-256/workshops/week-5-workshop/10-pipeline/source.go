package main

import (
	"bytes"
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
	"strings"
)

func source(ctx context.Context, categories []string) <-chan *link {
	content, err := getUrlContent("https://github.com/avelino/awesome-go")
	if err != nil {
		panic(err)
	}
	var (
		categoriesDictionary = make(map[string]struct{})
		out                  = make(chan *link)
	)

	for _, c := range categories {
		categoriesDictionary[strings.ToLower(c)] = struct{}{}
	}

	go func() {
		defer close(out)

		var (
			doc               *goquery.Document
			readMEHtmlContent string
		)

		if doc, err = goquery.NewDocumentFromReader(bytes.NewBuffer(content)); err != nil {
			panic(err)
		}

		doc.Find(`script[type="application/json"]`).Each(func(i int, selection *goquery.Selection) {
			value := gjson.Get(selection.Text(), "props.initialPayload.overview.overviewFiles")
			if !value.Exists() {
				return
			}
			files := value.Value().([]any)
			for _, f := range files {
				file := f.(map[string]any)
				name := file["displayName"].(string)
				if name == "README.md" {
					readMEHtmlContent = file["richText"].(string)
					break
				}
			}
		})
		if len(readMEHtmlContent) == 0 {
			panic("rich content is empty")
		}
		doc, err = goquery.NewDocumentFromReader(strings.NewReader(readMEHtmlContent))
		if err != nil {
			panic(err)
		}

		doc.Find(`.markdown-heading > h2, .markdown-heading > h3`).Each(func(i int, selection *goquery.Selection) {
			if i == 0 {
				return
			}
			category := selection.Text()
			_, suitable := categoriesDictionary[strings.ToLower(category)]
			if !suitable {
				return
			}
			selection.Parent().Next().Next().Find("li").Each(func(i int, selection *goquery.Selection) {
				var (
					l       = selection.Find("a")
					href, _ = l.Attr("href")
					name    = l.Text()
				)
				select {
				case out <- &link{
					name:     name,
					link:     href,
					category: category,
				}:
				case <-ctx.Done():
					return
				}
			})
		})

	}()

	return out
}
