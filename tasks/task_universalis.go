package tasks

import (
	"context"
	"github.com/Snipa22/core-go-lib/milieu"
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

type universalisSale struct {
	BuyerName    *string
	HQ           bool
	OnMannequin  bool
	Timestamp    int
	PricePerUnit int
	Quantity     int
	Total        int
	WorldID      *int
	WorldName    *string
}

type universalisReceive struct {
	Event    string
	Item     int
	Listings []universalisListing
	Sales    []universalisSale
	World    int
}

var bg = context.Background()

func UniversalisSocket(m *milieu.Milieu) {
	defer func() {
		if r := recover(); r != nil {
			UniversalisSocket(m)
		}
	}()
	c, _, err := websocket.Dial(bg, "wss://universalis.app/api/ws", nil)
	c.SetReadLimit(32768 * 4)
	if err != nil {
		m.CaptureException(err)
		m.Error(err.Error())
		return
	}
	defer func() {
		m.CaptureException(c.Close(websocket.StatusInternalError, "Websocket Closed"))
	}()
	msg, _ := bson.Marshal(universalisSub{"subscribe", "listings/add"})
	if err = c.Write(bg, websocket.MessageBinary, msg); err != nil {
		m.CaptureException(err)
		return
	}
	msg, _ = bson.Marshal(universalisSub{"subscribe", "listings/remove"})
	if err = c.Write(bg, websocket.MessageBinary, msg); err != nil {
		m.CaptureException(err)
		return
	}
	msg, _ = bson.Marshal(universalisSub{"subscribe", "sales/add"})
	if err = c.Write(bg, websocket.MessageBinary, msg); err != nil {
		m.CaptureException(err)
		return
	}
	msg, _ = bson.Marshal(universalisSub{"subscribe", "sales/remove"})
	if err = c.Write(bg, websocket.MessageBinary, msg); err != nil {
		m.CaptureException(err)
		return
	}
	for {
		_, n, err := c.Read(bg)
		if err != nil {
			m.CaptureException(err)
			return
		}
		data := universalisReceive{}
		if err = bson.Unmarshal(n, &data); err != nil {
			m.CaptureException(err)
		} else {
			curTime := time.Now()
			if data.Event == "listings/add" {
				if _, err = m.GetRawPGXPool().Exec(bg, "delete from items where item_id = $1 and world_id = $2", data.Item, data.World); err != nil {
					m.CaptureException(err)
				}
				for _, v := range data.Listings {
					_, _ = m.GetRawPGXPool().Exec(bg, "insert into items (id, world_id, item_id, price, total, hq, date_updated, quantity) values ($1, $2, $3, $4, $5, $6, $8, $7) ON CONFLICT (id) DO UPDATE SET price = $4, total = $5, date_updated = $8, quantity = $7", v.ListingID, data.World, data.Item, v.PricePerUnit, v.Total, v.HQ, v.Quantity, curTime)
				}
			} else if data.Event == "listings/remove" {
				for _, v := range data.Listings {
					if _, err = m.GetRawPGXPool().Exec(bg, "delete from items where id = $1", v.ListingID); err != nil {
						m.CaptureException(err)
					}
				}
			} else if data.Event == "sales/add" {
				for _, v := range data.Sales {
					_, _ = m.GetRawPGXPool().Exec(bg, "insert into sales (world_id, item_id, price, total, hq, quantity, date_loaded) values ($1, $2, $3, $4, $5, $6, $7)", data.World, data.Item, v.PricePerUnit, v.Total, v.HQ, v.Quantity, time.Unix(int64(v.Timestamp), 0))
				}
			}
		}
	}
}
