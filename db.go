package main

import (
	"encoding/json"
	"strings"

	"github.com/cockroachdb/pebble"
)

type RefererData struct {
	Credits int
	Uses    int
}

func checkRequest(referer string) bool {
	key := []byte(strings.ToLower(referer))
	log := log.With().Str("referer", referer).Logger()

	batch := db.NewBatch()

	var data RefererData

	if b, closer, err := batch.Get(key); err == nil {
		if err := json.Unmarshal(b, &data); err != nil {
			log.Warn().Err(err).Msg("failed to parse JSON from db, allowing")
			closer.Close()
			return true
		}

		closer.Close()
	} else {
		data = RefererData{}
	}

	data.Uses++
	dataj, _ := json.Marshal(data)
	log.Debug().Str("data", string(dataj)).Msg("referer data updated")

	if err := batch.Set(key, dataj, pebble.Sync); err != nil {
		log.Error().Err(err).Msg("error saving referer data")
	}

	if err := db.Apply(batch, pebble.Sync); err != nil {
		log.Error().Err(err).Msg("error applying batch")
	}

	return true
	// return data.Uses < data.Credits
}
