package tuning

import (
	"context"

	definition "github.com/bigstack-oss/cube-api/internal/definition/v1"
	"github.com/bigstack-oss/cube-api/internal/service"
	log "go-micro.dev/v5/logger"
	"go-micro.dev/v5/registry"
)

var (
	module = "node"
)

func init() {
	service.RegisterController(module, NewController())
}

type Controller struct {
	ctx             context.Context
	cancel          context.CancelFunc
	isFirstTimeSync bool
}

func NewController() *Controller {
	ctx, cancel := context.WithCancel(context.Background())
	return &Controller{
		ctx:             ctx,
		cancel:          cancel,
		isFirstTimeSync: true,
	}
}

func (c *Controller) Name() string {
	return module
}

func (c *Controller) Sync() {
	watcher, err := registry.Watch()
	if err != nil {
		log.Errorf("Failed to create watcher (%s)", err.Error())
		return
	}

	defer watcher.Stop()
	select {
	case <-c.ctx.Done():
		return
	default:
		c.watchAndSyncNodeRoles(&watcher)
	}
}

func (c *Controller) Stop() {
	c.cancel()
}

func (c *Controller) watchAndSyncNodeRoles(watcher *registry.Watcher) {
	event, err := (*watcher).Next()
	if err == nil {
		definition.SyncNodesOfRole()
		logWithThrottling(event)
		return
	}

	log.Errorf(
		"Failed to get service discovery event",
		err.Error(),
	)
}
