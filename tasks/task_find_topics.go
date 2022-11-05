package tasks

import (
	"context"
	"github.com/Snipa22/xiv-lodestone-go/support"
	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/html"
	"strconv"
	"strings"
	"time"
)

func SetupGettersForTopics(milieu support.Milieu) func() {
	return func() {
		// Loop through all lodestones, download all the data required.
		for _, v := range []support.Regions{0, 1, 2, 3, 4} {
			// HTTP request the /lodestone/news/category/2 path
			baseURI := support.GetLodestoneBaseURI(v)
			baseURI += "/lodestone/topics"
			page, err := support.GetHtmlPage(baseURI)
			if err != nil {
				sentry.CaptureException(err)
				continue
			}
			tkn := html.NewTokenizer(strings.NewReader(page))
			inTextParseZone := false
			inMaintLine := false
			maintLine := ""
			maintURL := ""
			var hash, topicBody string
			var val int
			for {
				tt := tkn.Next()
				if tt == html.StartTagToken {
					t := tkn.Token()
					if t.Data == "a" && t.Attr[0].Key == "href" && strings.Contains(t.Attr[0].Val, "/lodestone/topics/detail/") {
						maintURL = support.GetLodestoneBaseURI(v) + t.Attr[0].Val
					}
					if t.Data == "header" && len(t.Attr) == 1 && t.Attr[0].Key == "class" && t.Attr[0].Val == "news__list--header clearfix" {
						inMaintLine = true
					}
					if t.Data == "div" && len(t.Attr) == 1 && t.Attr[0].Key == "class" && t.Attr[0].Val == "news__list--banner" {
						inTextParseZone = true
					}
				}
				if inMaintLine && tt == html.TextToken {
					t := tkn.Token()
					if inMaintLine && len(maintLine) == 0 {
						maintLine += t.Data
					}
					if strings.Contains(t.Data, "ldst_strftime") {
						ts := extractStrftime.FindString(t.Data)[1:]
						hash = extractHash.FindString(maintURL)
						// Do the time conversion
						val, err = strconv.Atoi(ts)
						if err != nil {
							sentry.CaptureException(err)
						}
						inMaintLine = false
					}
				}
				if inTextParseZone && tt == html.TextToken && !inMaintLine {
					topicBody += tkn.Token().String()
				}
				if tt == html.EndTagToken {
					t := tkn.Token()
					if t.Data == "div" && inTextParseZone {
						inTextParseZone = false
						row := milieu.Pgx.QueryRow(context.Background(), "select id from ls_topics where id = $1 and region = $2", hash, v)

						var bid string
						if err := row.Scan(&bid); err != nil && err == pgx.ErrNoRows {
							// Do the SQL insert if appropriate
							_, err = milieu.Pgx.Exec(context.Background(), "insert into ls_topics (id, region, title, uri, square_edit, topic_body)"+
								"values ($1, $2, $3, $4, $5, $6) on conflict do nothing", hash, v, maintLine, maintURL, time.Unix(int64(val), 0), topicBody)
							if err != nil {
								sentry.CaptureException(err)
							}
						}
						inTextParseZone = false
						inMaintLine = false
						maintLine = ""
						maintURL = ""
						hash = ""
						topicBody = ""
					}
				}
				if tt == html.ErrorToken {
					break
				}
			}
		}
	}
}
