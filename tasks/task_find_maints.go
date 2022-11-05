package tasks

import (
	"context"
	"fmt"
	"github.com/Snipa22/xiv-lodestone-go/support"
	"github.com/getsentry/sentry-go"
	"golang.org/x/net/html"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var extractStrftime, _ = regexp.Compile(`\(\d+`)
var extractHash, _ = regexp.Compile(`\w{40}`)

func SetupGetMaintencePages(milieu support.Milieu) func() {
	return func() {
		// Loop through all lodestones, download all the data required.
		for _, v := range []support.Regions{0, 1, 2, 3, 4} {
			fmt.Println(v)
			// HTTP request the /lodestone/news/category/2 path
			baseURI := support.GetLodestoneBaseURI(v)
			baseURI += "/lodestone/news/category/2"
			page, err := support.GetHtmlPage(baseURI)
			if err != nil {
				sentry.CaptureException(err)
				continue
			}
			tkn := html.NewTokenizer(strings.NewReader(page))
			inMaintLine := false
			SummedLines := 0
			maintLine := ""
			maintURL := ""
			for {
				tt := tkn.Next()
				if tt == html.StartTagToken {
					t := tkn.Token()
					if t.Data == "a" && len(t.Attr) == 2 && t.Attr[1].Key == "class" && t.Attr[1].Val == "news__list--link ic__maintenance--list" {
						maintURL = support.GetLodestoneBaseURI(v) + t.Attr[0].Val
					}
					if t.Data == "div" && len(t.Attr) == 1 && t.Attr[0].Key == "class" && t.Attr[0].Val == "clearfix" {
						inMaintLine = true
					}
				}
				if inMaintLine && tt == html.TextToken {
					t := tkn.Token()
					if SummedLines < 2 {
						maintLine += t.Data
						SummedLines += 1
					}
					if strings.Contains(t.Data, "ldst_strftime") {
						ts := extractStrftime.FindString(t.Data)[1:]
						hash := extractHash.FindString(maintURL)
						// Do the time conversion
						val, err := strconv.Atoi(ts)
						if err != nil {
							sentry.CaptureException(err)
						}
						// Do the SQL insert if appropriate
						_, err = milieu.Pgx.Exec(context.Background(), "insert into ls_maint (id, region, title, uri, square_edit)"+
							"values ($1, $2, $3, $4, $5) on conflict do nothing", hash, v, maintLine, maintURL, time.Unix(int64(val), 0))
						if err != nil {
							sentry.CaptureException(err)
						}
						inMaintLine = false
						SummedLines = 0
						maintLine = ""
					}
				}
				if tt == html.ErrorToken {
					break
				}
			}
		}
	}
}
