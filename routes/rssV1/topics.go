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

func GetTopicsForLang(region support.Regions) func(c *gin.Context) {
	return func(c *gin.Context) {
		milieu := middleware.MustGetMilieu(c)
		rows, err := milieu.Pgx.Query(context.Background(), "Select id, title, uri, square_edit, topic_body, topic_image from ls_topics where region = $1 order by date_found desc limit 10", region)
		if err != nil {
			sentry.CaptureException(err)
		}
		feed := &feeds.Feed{
			Title:       "FFXIV Lodestone topics RSS feed for " + region.String(),
			Link:        &feeds.Link{Href: "https://xivrss.jagtech.io/rss/topics/" + region.URIExtensions()},
			Description: "Lodestone Topics Notifications",
			Author:      &feeds.Author{Name: "Impala#0059"},
			Created:     time.Now(),
		}
		var feedItems []*feeds.Item
		for rows.Next() {
			var id, title, uri, maintBody, topicImage string
			var edit time.Time
			if err = rows.Scan(&id, &title, &uri, &edit, &maintBody, &topicImage); err != nil {
				sentry.CaptureException(err)
			}
			feedItems = append(feedItems,
				&feeds.Item{
					Id:          topicImage,
					Title:       title,
					Link:        &feeds.Link{Href: uri},
					Description: maintBody,
					Created:     edit,
					Author:      &feeds.Author{Name: id, Email: topicImage},
				})
		}
		feed.Items = feedItems
		rssFeed := (&feeds.Rss{Feed: feed}).RssFeed()
		xmlRssFeeds := rssFeed.FeedXml()
		c.XML(200, xmlRssFeeds)
	}
}
