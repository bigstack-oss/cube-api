package tuning

import (
	"context"
	"fmt"
	"time"

	"github.com/cnf/structhash"
	"go-micro.dev/v5/cache"
	log "go-micro.dev/v5/logger"
	"go-micro.dev/v5/registry"
)

var (
	logCache = cache.NewCache(cache.Expiration(time.Second * 10))
)

func convertAction(action string) string {
	switch action {
	case "create":
		return "joined"
	case "delete":
		return "left"
	}

	return action
}

func logWithThrottling(event *registry.Result) {
	key, err := structhash.Hash(event, 1)
	if err != nil {
		return
	}

	_, _, err = logCache.Get(context.Background(), key)
	if err == nil {
		return
	}
	if len(event.Service.Nodes) == 0 {
		return
	}

	logCache.Put(context.Background(), key, []byte{}, time.Second*10)
	log.Infof(
		"Node resynced: %s %s %s",
		event.Service.Name,
		fmt.Sprintf("%s(%s)", event.Service.Metadata["hostname"], event.Service.Nodes[0].Address),
		convertAction(event.Action),
	)
}
