package tasks

import (
	"context"
	"git.jagtech.io/ffxiv/xiv-lodestone-go/support"
	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/html"
	"strconv"
	"strings"
	"time"
)

func SetupGetNoticePages(milieu support.Milieu) func() {
	return func() {
		// Loop through all lodestones, download all the data required.
		for _, v := range []support.Regions{0, 1, 2, 3, 4} {
			// HTTP request the /lodestone/news/category/2 path
			baseURI := support.GetLodestoneBaseURI(v)
			baseURI += "/lodestone/news/category/1"
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
					if t.Data == "a" && len(t.Attr) == 2 && t.Attr[1].Key == "class" && t.Attr[1].Val == "news__list--link ic__info--list" {
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
						row := milieu.Pgx.QueryRow(context.Background(), "select id from ls_notices where id = $1 and region = $2", hash, v)
						var bid string
						if err := row.Scan(&bid); err != nil && err == pgx.ErrNoRows {
							// Get the full data set
							internalPage, err := support.GetHtmlPage(maintURL)
							if err != nil {
								sentry.CaptureException(err)
								continue
							}
							intTkn := html.NewTokenizer(strings.NewReader(internalPage))
							inMaintBody := false
							maintBody := ""
							for {
								tokenLoop := intTkn.Next()
								if tokenLoop == html.StartTagToken {
									t := intTkn.Token()
									if t.Data == "div" && len(t.Attr) == 1 && t.Attr[0].Key == "class" && t.Attr[0].Val == "news__detail__wrapper" {
										inMaintBody = true
									}
									if t.Data == "div" && len(t.Attr) == 1 && t.Attr[0].Key == "class" && t.Attr[0].Val == "news__footer" {
										inMaintBody = false
									}
								}
								if inMaintBody && tokenLoop == html.TextToken {
									t := intTkn.Token()
									maintBody += t.String()
								}
								if tokenLoop == html.ErrorToken {
									break
								}
							}
							// Do the SQL insert if appropriate
							_, err = milieu.Pgx.Exec(context.Background(), "insert into ls_notices (id, region, title, uri, square_edit, notice_body)"+
								"values ($1, $2, $3, $4, $5, $6) on conflict do nothing", hash, v, maintLine, maintURL, time.Unix(int64(val), 0), maintBody)
							if err != nil {
								sentry.CaptureException(err)
							}
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
