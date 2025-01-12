package tuning

import (
	"time"

	definition "github.com/bigstack-oss/cube-cos-api/internal/definition/v1"
	"github.com/bigstack-oss/cube-cos-api/internal/helpers/mongo"
	"github.com/bigstack-oss/cube-cos-api/internal/service"
	"k8s.io/client-go/util/workqueue"
)

var (
	ReqQueue workqueue.Interface
	module   = "tuning"
)

func init() {
	ReqQueue = workqueue.New()
	service.RegisterController(module, NewController())
}

type Controller struct {
	mongo *mongo.Helper
}

func NewController() *Controller {
	m := mongo.NewHelper(mongo.NewDefaultConf("tuning"))
	return &Controller{mongo: m}
}

func (c *Controller) Name() string {
	return module
}

func (c *Controller) Sync() {
	req, shutdown := ReqQueue.Get()
	if shutdown {
		return
	}

	tuning := req.(definition.Tuning)
	err := c.syncByDesiredAction(tuning)

	c.handleExit(tuning, err)
	ReqQueue.Done(req)
}

func (c *Controller) Stop() {
	ReqQueue.ShutDown()
	c.waitForLastTask()
	c.mongo.Close()
}

func (c *Controller) waitForLastTask() {
	for ReqQueue.Len() >= 1 {
		time.Sleep(time.Second * 1)
	}
}
