package support

import (
	"context"
	"errors"
	"github.com/Snipa22/core-go-lib/milieu"
	lru "github.com/hashicorp/golang-lru/v2"
)

var worlds, _ = lru.New[int, string](256)

func GetWorldName(m milieu.Milieu, worldInternalID int) (string, error) {
	if res, ok := worlds.Get(worldInternalID); ok {
		return res, nil
	}
	row := m.GetRawPGXPool().QueryRow(context.Background(), `select display_name from sq_worlds where internal_id = $1`, worldInternalID)
	res := ""
	if err := row.Scan(&res); err != nil {
		return "", err
	}
	if len(res) == 0 {
		return "", errors.New("world not found")
	}
	worlds.Add(worldInternalID, res)
	return res, nil
}
