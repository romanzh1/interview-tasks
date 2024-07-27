package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
)

var (
	re = regexp.MustCompile("[0-9]+")
)

func getStars(ctx context.Context, link string) (result int, err error) {
	b, err := getUrlContent(link)
	if err != nil {
		return 0, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
	if err != nil {
		return 0, fmt.Errorf("get %s: parse page: %w", link, err)
	}
	var (
		found      = false
		parseError error
	)
	doc.Find("#repo-stars-counter-star").Each(func(i int, selection *goquery.Selection) {
		if parseError != nil {
			return
		}
		v, ok := selection.Attr("aria-label")
		if !ok {
			parseError = fmt.Errorf("empty attr[aria-label] of #repo-stars-counter-star")
			return
		}
		m := re.FindAllString(v, 1)
		if len(m) == 0 {
			parseError = fmt.Errorf("can't extract number of star from value: %s", v)
			return
		}
		if result, err = strconv.Atoi(m[0]); err != nil {
			parseError = fmt.Errorf("can't parse %s: %w", m[0], err)
			return
		}
		found = true
	})
	if parseError != nil {
		return 0, parseError
	}
	if !found {
		return 0, fmt.Errorf("can't find #repo-stars-counter-star element: %s", link)
	}
	return
}
