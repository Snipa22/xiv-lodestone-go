package marketV1

import (
	"context"
	"git.jagtech.io/Impala/corelib/middleware"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/montanaflynn/stats"
	"strconv"
	"time"
	"xiv-lodestone-go/support"
)

var bg = context.Background()

type kupoMarketShell struct {
	ItemID                int                  `json:"itemID"`
	WorldID               int                  `json:"worldID"`
	LastUploadTime        int64                `json:"lastUploadTime"`
	Listings              []kupoMarketListings `json:"listings"`
	WorldName             string               `json:"worldName"`
	MinPrice              int                  `json:"minPrice"`
	MinPriceHQ            int                  `json:"minPriceHQ"`
	CurrentAveragePrice   int                  `json:"currentAveragePrice"`
	AveragePrice          int                  `json:"averagePrice"`
	CurrentAveragePriceNQ int                  `json:"currentAveragePriceNQ"`
	AveragePriceNQ        int                  `json:"averagePriceNQ"`
	CurrentAveragePriceHQ int                  `json:"currentAveragePriceHQ"`
	AveragePriceHQ        int                  `json:"averagePriceHQ"`
}

type kupoMarketListings struct {
	PricePerUnit int    `json:"pricePerUnit"`
	Total        int    `json:"total"`
	Quantity     int    `json:"quantity"`
	HQ           bool   `json:"hq"`
	Timestamp    int64  `json:"timestamp"`
	WorldName    string `json:"worldName"`
}

func getMean(in []int) int {
	if len(in) == 0 {
		return 0
	}
	data := stats.LoadRawData(in)
	outliers, err := data.QuartileOutliers()
	if err != nil {
		sentry.CaptureException(err)
		return 0
	}
	var newParse []float64
	for _, v := range data {
		found := false
		for _, e := range outliers.Extreme {
			if v == e {
				found = true
			}
		}
		if !found {
			newParse = append(newParse, v)
		}
	}
	data = stats.LoadRawData(newParse)
	retVal, _ := data.Mean()
	return int(retVal)
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
	defer rows.Close()
	returnData := kupoMarketShell{
		ItemID:                itemID,
		WorldID:               0,
		LastUploadTime:        lastUploadTime,
		Listings:              make([]kupoMarketListings, 0),
		WorldName:             "",
		MinPrice:              0, // Done
		MinPriceHQ:            0, // Done
		CurrentAveragePrice:   0,
		AveragePrice:          0,
		CurrentAveragePriceNQ: 0,
		AveragePriceNQ:        0,
		CurrentAveragePriceHQ: 0,
		AveragePriceHQ:        0,
	}
	curPrice := make([]int, 0)
	curPriceHQ := make([]int, 0)
	curPriceNQ := make([]int, 0)
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
			returnData.LastUploadTime = res.Timestamp
		}
		returnData.Listings = append(returnData.Listings, res)
		if res.PricePerUnit < returnData.MinPrice || returnData.MinPrice == 0 {
			returnData.MinPrice = res.PricePerUnit
		}
		if res.HQ && (res.PricePerUnit < returnData.MinPriceHQ || returnData.MinPriceHQ == 0) {
			returnData.MinPriceHQ = res.PricePerUnit
		}
		curPrice = append(curPrice, res.PricePerUnit)
		if res.HQ {
			curPriceHQ = append(curPriceHQ, res.PricePerUnit)
		} else {
			curPriceNQ = append(curPriceNQ, res.PricePerUnit)
		}
	}
	returnData.CurrentAveragePrice = getMean(curPrice)
	returnData.CurrentAveragePriceHQ = getMean(curPriceHQ)
	returnData.CurrentAveragePriceNQ = getMean(curPriceNQ)
	// Time to get the historical/sales data.
	query = `select sales.price, sales.hq from sales where item_id = $2 and world_id in (select sqw.internal_id from sq_worlds as sqw
		join sq_logical_datacenters sld on sld.id = sqw.sq_logical_datacenter_id
		join sq_physical_datacenters spd on spd.id = sld.physical_dc_id
		where lower(spd.display_name) = lower($1)
		or lower(spd.internal_name) = lower($1)
		or lower(sld.display_name) = lower($1)
		or lower(sld.internal_name) = lower($1)
		or lower(sqw.display_name) = lower($1)
		or lower(sqw.internal_name) = lower($1))`
	rows, err = milieu.Pgx.Query(bg, query, world, itemID)
	if err != nil && err != pgx.ErrNoRows {
		sentry.CaptureException(err)
		support.Err400(c, "No rows found")
		return
	} else if err != nil && err == pgx.ErrNoRows {
		c.JSON(200, returnData)
		return
	}
	defer rows.Close()
	histPrice := make([]int, 0)
	histPriceHQ := make([]int, 0)
	histPriceNQ := make([]int, 0)
	for rows.Next() {
		hq := false
		price := 0
		if err = rows.Scan(&price, &hq); err != nil {
			sentry.CaptureException(err)
			support.Err500(c)
			return
		}
		histPrice = append(histPrice, price)
		if hq {
			histPriceHQ = append(histPriceHQ, price)
		} else {
			histPriceNQ = append(histPriceNQ, price)
		}
	}
	returnData.AveragePrice = getMean(histPrice)
	returnData.AveragePriceHQ = getMean(histPriceHQ)
	returnData.AveragePriceNQ = getMean(histPriceNQ)
	c.JSON(200, returnData)
	return
}
