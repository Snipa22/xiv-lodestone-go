package marketV1

import (
	"context"
	"git.jagtech.io/Impala/corelib/middleware"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"xiv-lodestone-go/support"
)

var bg = context.Background()

type kupoMarketShell struct {
	ItemID         int                  `json:"itemID"`
	WorldID        int                  `json:"worldID"`
	LastUploadTime int64                `json:"lastUploadTime"`
	Listings       []kupoMarketListings `json:"listings"`
	WorldName      string               `json:"worldName"`
}

type kupoMarketListings struct {
	PricePerUnit int    `json:"pricePerUnit"`
	Total        int    `json:"total"`
	Quantity     int    `json:"quantity"`
	HQ           bool   `json:"hq"`
	Timestamp    int64  `json:"timestamp"`
	WorldName    string `json:"worldName"`
}

func GetMarketData(c *gin.Context) {
	milieu := middleware.MustGetMilieu(c)
	world := c.Param("world")
	itemID, err := strconv.Atoi(c.Param("itemID"))
	if err != nil {
		sentry.CaptureException(err)
		support.Err400(c, "Unable to parse ItemID")
		return
	}
	var lastUploadTime int64
	query := `select items.world_id, items.price, items.total, items.hq, items.quantity, items.date_updated from items where item_id = $2 and world_id in (select sqw.internal_id from sq_worlds as sqw
		join sq_logical_datacenters sld on sld.id = sqw.sq_logical_datacenter_id
		join sq_physical_datacenters spd on spd.id = sld.physical_dc_id
		where lower(spd.display_name) = lower($1)
		or lower(spd.internal_name) = lower($1)
		or lower(sld.display_name) = lower($1)
		or lower(sld.internal_name) = lower($1)
		or lower(sqw.display_name) = lower($1)
		or lower(sqw.internal_name) = lower($1)) order by items.date_updated desc`
	rows, err := milieu.Pgx.Query(bg, query, world, itemID)
	if err != nil {
		sentry.CaptureException(err)
		support.Err400(c, "No rows found")
		return
	}
	results := make([]kupoMarketListings, 0)
	for rows.Next() {
		res := kupoMarketListings{}
		tData := time.Now()
		worldID := 0
		if err = rows.Scan(&worldID, &res.PricePerUnit, &res.Total, &res.HQ, &res.Quantity, &tData); err != nil {
			sentry.CaptureException(err)
			support.Err500(c)
			return
		}
		if res.WorldName, err = support.GetWorldName(milieu, worldID); err != nil {
			sentry.CaptureException(err)
			support.Err500(c)
			return
		}
		res.Timestamp = tData.Unix()
		if res.Timestamp > lastUploadTime {
			lastUploadTime = res.Timestamp
		}
		results = append(results, res)
	}
	c.JSON(200, kupoMarketShell{
		ItemID:         itemID,
		WorldID:        0,
		LastUploadTime: lastUploadTime,
		Listings:       results,
		WorldName:      "",
	})
	return
}
