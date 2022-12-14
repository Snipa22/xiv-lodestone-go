package main

import (
	"context"
	"github.com/Snipa22/xiv-lodestone-go/routes/rssV1"
	"github.com/Snipa22/xiv-lodestone-go/support"
	"github.com/Snipa22/xiv-lodestone-go/support/middleware"
	"github.com/Snipa22/xiv-lodestone-go/tasks"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"os"
)

var ctx = context.Background()

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_SERVER"),
		// Specify a fixed sample rate:
		// We recommend adjusting this value in production
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_SERVER"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pgpool, err := pgxpool.Connect(ctx, os.Getenv("PSQL_SERVER"))
	if err != nil {
		sentry.CaptureException(err)
		log.Fatalf("Unable to connect to PSQL: %s", os.Getenv("PSQL_SERVER"))
	}

	// We need to broadcast through a channel to each worker that
	// Prep-rep Milieu
	milieu := support.Milieu{
		Pgx:   pgpool,
		Redis: rdb,
	}
	// Enable Recurring Tasks
	c := cron.New(cron.WithSeconds())
	_, err = c.AddFunc("0 * * * * *", tasks.SetupGettersForTopics(milieu))
	_, err = c.AddFunc("0 * * * * *", tasks.SetupGetMaintencePages(milieu))
	_, err = c.AddFunc("0 * * * * *", tasks.SetupGetNoticePages(milieu))
	_, err = c.AddFunc("0 * * * * *", tasks.SetupGetStatusPages(milieu))
	if err != nil {
		sentry.CaptureException(err)
	}
	go c.Run()

	r := gin.Default()
	r.Use(middleware.SetupMilieu(milieu))
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
	r.GET("/ping", func(c *gin.Context) {
		milieu := c.MustGet("MILIEU").(support.Milieu)
		var i int
		err := milieu.Pgx.QueryRow(ctx, "select 1").Scan(&i)
		if err != nil {
			sentry.CaptureException(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"value":   i,
		})
	})

	err = r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		sentry.CaptureException(err)
		log.Fatalf("Unable to initalize gin")
	}
}
