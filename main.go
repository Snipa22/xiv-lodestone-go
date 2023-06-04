package main

import (
	"context"
	"github.com/Snipa22/core-go-lib/milieu"
	"github.com/Snipa22/core-go-lib/milieu/middleware"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"os"
	"xiv-lodestone-go/routes/marketV1"
	"xiv-lodestone-go/routes/rssV1"
	"xiv-lodestone-go/support"
	"xiv-lodestone-go/tasks"
)

var ctx = context.Background()

func main() {
	psqlURL := os.Getenv("PSQL_SERVER")
	redisURI := os.Getenv("REDIS_SERVER")
	sentryURI := os.Getenv("SENTRY_SERVER")
	m, err := milieu.NewMilieu(&psqlURL, &redisURI, &sentryURI)
	if err != nil {
		m.CaptureException(err)
		log.Fatalf("Unable to setup Milieu, check sentry for details")
	}

	// We need to broadcast through a channel to each worker that
	// Enable Recurring Tasks
	c := cron.New(cron.WithSeconds())
	_, err = c.AddFunc("0 * * * * *", tasks.SetupGettersForTopics(m))
	_, err = c.AddFunc("0 * * * * *", tasks.SetupGetMaintenancePages(m))
	_, err = c.AddFunc("0 * * * * *", tasks.SetupGetNoticePages(m))
	_, err = c.AddFunc("0 * * * * *", tasks.SetupGetStatusPages(m))
	if err != nil {
		m.CaptureException(err)
	}

	go func() {
		for {
			tasks.UniversalisSocket(m)
		}
	}()

	go c.Run()

	r := gin.Default()
	r.Use(middleware.SetupMilieu(m))
	rss := r.Group("/rss")
	{
		maint := rss.Group("/maint")
		{
			maint.GET("/NA", rssV1.GetMaintForLang(support.NA))
			maint.GET("/EU", rssV1.GetMaintForLang(support.EU))
			maint.GET("/JP", rssV1.GetMaintForLang(support.JP))
			maint.GET("/FR", rssV1.GetMaintForLang(support.FR))
			maint.GET("/DE", rssV1.GetMaintForLang(support.DE))
		}
		status := rss.Group("/status")
		{
			status.GET("/NA", rssV1.GetStatusForLang(support.NA))
			status.GET("/EU", rssV1.GetStatusForLang(support.EU))
			status.GET("/JP", rssV1.GetStatusForLang(support.JP))
			status.GET("/FR", rssV1.GetStatusForLang(support.FR))
			status.GET("/DE", rssV1.GetStatusForLang(support.DE))
		}
		notices := rss.Group("/notices")
		{
			notices.GET("/NA", rssV1.GetNoticeForLang(support.NA))
			notices.GET("/EU", rssV1.GetNoticeForLang(support.EU))
			notices.GET("/JP", rssV1.GetNoticeForLang(support.JP))
			notices.GET("/FR", rssV1.GetNoticeForLang(support.FR))
			notices.GET("/DE", rssV1.GetNoticeForLang(support.DE))
		}
		topics := rss.Group("/topics")
		{
			topics.GET("/NA", rssV1.GetTopicsForLang(support.NA))
			topics.GET("/EU", rssV1.GetTopicsForLang(support.EU))
			topics.GET("/JP", rssV1.GetTopicsForLang(support.JP))
			topics.GET("/FR", rssV1.GetTopicsForLang(support.FR))
			topics.GET("/DE", rssV1.GetTopicsForLang(support.DE))
		}
	}

	r.GET("/market/:world/:itemID", marketV1.GetMarketData)
	r.GET("/ping", func(c *gin.Context) {
		m := middleware.MustGetMilieu(c)
		var i int
		err := m.GetRawPGXPool().QueryRow(ctx, "select 1").Scan(&i)
		if err != nil {
			m.CaptureException(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"value":   i,
		})
	})

	err = r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		m.CaptureException(err)
		log.Fatalf("Unable to initalize gin")
	}
}
