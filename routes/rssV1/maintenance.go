package rssV1

import (
	"context"
	"github.com/Snipa22/xiv-lodestone-go/support"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"time"
)

func GetMaintForLang(region support.Regions) func(c *gin.Context) {
	return func(c *gin.Context) {
		milieu := c.MustGet("MILIEU").(support.Milieu)
		rows, err := milieu.Pgx.Query(context.Background(), "Select id, title, uri, square_edit from ls_maint where region = $1 order by date_found desc limit 10", region)
		if err != nil {
			sentry.CaptureException(err)
		}
		feed := &feeds.Feed{
			Title:       "FFXIV Lodestone Maintenance RSS feed for " + region.String(),
			Link:        &feeds.Link{Href: "https://xivrss.jagtech.io/rss/maint/" + region.URIExtensions()},
			Description: "Lodestone Maintenance Notifications",
			Author:      &feeds.Author{Name: "Impala#0059"},
			Created:     time.Now(),
		}
		var feedItems []*feeds.Item
		for rows.Next() {
			var id, title, uri string
			var edit time.Time
			if err = rows.Scan(&id, &title, &uri, &edit); err != nil {
				sentry.CaptureException(err)
			}
			feedItems = append(feedItems,
				&feeds.Item{
					Id:          id,
					Title:       title,
					Link:        &feeds.Link{Href: uri},
					Description: "None",
					Created:     edit,
				})
		}
		feed.Items = feedItems
		rssFeed := (&feeds.Rss{Feed: feed}).RssFeed()
		xmlRssFeeds := rssFeed.FeedXml()
		c.XML(200, xmlRssFeeds)
	}
}
