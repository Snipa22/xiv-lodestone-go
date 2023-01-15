package support

import (
	"context"
	"errors"
	"git.jagtech.io/Impala/corelib"
	lru "github.com/hashicorp/golang-lru/v2"
)

var worlds, _ = lru.New[int, string](256)

func GetWorldName(milieu corelib.Milieu, worldInternalID int) (string, error) {
	if res, ok := worlds.Get(worldInternalID); ok {
		return res, nil
	}
	row := milieu.Pgx.QueryRow(context.Background(), `select display_name from sq_worlds where internal_id = $1`, worldInternalID)
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
