package rssV1

import (
	"context"
	"git.jagtech.io/Impala/corelib/middleware"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"time"
	"xiv-lodestone-go/support"
)

func GetNoticeForLang(region support.Regions) func(c *gin.Context) {
	return func(c *gin.Context) {
		milieu := middleware.MustGetMilieu(c)
		rows, err := milieu.Pgx.Query(context.Background(), "Select id, title, uri, square_edit, notice_body from ls_notices where region = $1 order by date_found desc limit 10", region)
		if err != nil {
			sentry.CaptureException(err)
		}
		feed := &feeds.Feed{
			Title:       "FFXIV Lodestone Notice RSS feed for " + region.String(),
			Link:        &feeds.Link{Href: "https://xivrss.jagtech.io/rss/notices/" + region.URIExtensions()},
			Description: "Lodestone Notice Notifications",
			Author:      &feeds.Author{Name: "Impala#0059"},
			Created:     time.Now(),
		}
		var feedItems []*feeds.Item
		for rows.Next() {
			var id, title, uri, maint_body string
			var edit time.Time
			if err = rows.Scan(&id, &title, &uri, &edit, &maint_body); err != nil {
				sentry.CaptureException(err)
			}
			feedItems = append(feedItems,
				&feeds.Item{
					Id:          id,
					Title:       title,
					Link:        &feeds.Link{Href: uri},
					Description: maint_body,
					Created:     edit,
					Author:      &feeds.Author{Name: id, Email: uri},
				})
		}
		feed.Items = feedItems
		rssFeed := (&feeds.Rss{Feed: feed}).RssFeed()
		xmlRssFeeds := rssFeed.FeedXml()
		c.XML(200, xmlRssFeeds)
	}
}
