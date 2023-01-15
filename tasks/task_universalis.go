package tasks

import (
	"context"
	"git.jagtech.io/Impala/corelib"
	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson"
	"nhooyr.io/websocket"
	"time"
)

type universalisSub struct {
	Event   string
	Channel string
}

type universalisMateria struct {
	SlotID    int
	MateriaID int
}

type universalisListing struct {
	CreatorID      *string
	CreatorName    *string
	HQ             bool
	IsCrafted      bool
	LastReviewTime int
	ListingID      string
	Materia        []universalisMateria
	OnMannequin    bool
	PricePerUnit   int
	Quantity       int
	RetainerCity   int
	RetainerID     string
	RetainerName   string
	SellerID       string
	StainID        int
	Total          int
	WorldID        *int
	WorldName      *string
}

type universalisReceive struct {
	Event    string
	Item     int
	Listings []universalisListing
	Sales    []string
	World    int
}

var bg = context.Background()

func UniversalisSocket(milieu corelib.Milieu) {
	defer func() {
		if r := recover(); r != nil {
			UniversalisSocket(milieu)
		}
	}()
	c, _, err := websocket.Dial(bg, "wss://universalis.app/api/ws", nil)
	c.SetReadLimit(32768 * 4)
	if err != nil {
		sentry.CaptureException(err)
		milieu.Error(err.Error())
		return
	}
	defer func() {
		sentry.CaptureException(c.Close(websocket.StatusInternalError, "Websocket Closed"))
	}()
	msg, _ := bson.Marshal(universalisSub{"subscribe", "listings/add"})
	if err = c.Write(bg, websocket.MessageBinary, msg); err != nil {
		sentry.CaptureException(err)
		return
	}
	msg, _ = bson.Marshal(universalisSub{"subscribe", "listings/remove"})
	if err = c.Write(bg, websocket.MessageBinary, msg); err != nil {
		sentry.CaptureException(err)
		return
	}
	/*
		Disabling Sales as it's likely not needed for what I'm doing
		msg, _ = bson.Marshal(universalisSub{"subscribe", "sales/add"})
		if err = c.Write(bg, websocket.MessageBinary, msg); err != nil {
			sentry.CaptureException(err)
			return
		}
		msg, _ = bson.Marshal(universalisSub{"subscribe", "sales/remove"})
		if err = c.Write(bg, websocket.MessageBinary, msg); err != nil {
			sentry.CaptureException(err)
			return
		}
	*/
	for {
		_, n, err := c.Read(bg)
		if err != nil {
			sentry.CaptureException(err)
			return
		}
		data := universalisReceive{}
		if err = bson.Unmarshal(n, &data); err != nil {
			sentry.CaptureException(err)
		} else {
			curTime := time.Now()
			if data.Event == "listings/add" {
				for _, v := range data.Listings {
					_, _ = milieu.Pgx.Exec(bg, "insert into items (id, world_id, item_id, price, total, hq, date_updated, quantity) values ($1, $2, $3, $4, $5, $6, $8, $7) ON CONFLICT (id) DO UPDATE SET price = $4, total = $5, date_updated = $8, quantity = $7", v.ListingID, data.World, data.Item, v.PricePerUnit, v.Total, v.HQ, v.Quantity, curTime)
				}
			} else if data.Event == "listings/remove" {
				for _, v := range data.Listings {
					if _, err = milieu.Pgx.Exec(bg, "delete from items where id = $1", v.ListingID); err != nil {
						sentry.CaptureException(err)
					}
				}
			}
		}
	}
}
