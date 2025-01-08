package runtime

import (
	"fmt"
	"os"

	"github.com/bigstack-oss/cube-cos-api/internal/api"
	"github.com/bigstack-oss/cube-cos-api/internal/auth"
	apiConf "github.com/bigstack-oss/cube-cos-api/internal/config"
	"github.com/bigstack-oss/cube-cos-api/internal/cubecos"
	definition "github.com/bigstack-oss/cube-cos-api/internal/definition/v1"
	"github.com/bigstack-oss/cube-cos-api/internal/helpers/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/micro/plugins/v5/server/http"
	"go-micro.dev/v5/config"
	"go-micro.dev/v5/logger"
	"go-micro.dev/v5/server"
)

func setGroupHandlersToRouter(router *gin.Engine, handlers []api.Handler) {
	for _, h := range handlers {
		if h.Version == "" {
			logger.Warnf("Skip invalid API registration: %s %s (no version provided)", h.Method, h.Path)
			continue
		}

		routerGroup := router.Group(h.Version)
		routerGroup.Handle(h.Method, h.Path, h.Func)
		logger.Infof("Register API: %s %s", h.Method, fmt.Sprintf("%s%s", h.Version, h.Path))
	}
}

func RegisterHandlersByRole(router *gin.Engine) error {
	groupHandlers := api.GetGroupHandlersByRole(definition.CurrentRole)
	if len(groupHandlers) == 0 {
		return fmt.Errorf("no handlers found for role(%s)", definition.CurrentRole)
	}

	for _, handlers := range groupHandlers {
		setGroupHandlersToRouter(router, handlers)
	}

	return nil
}

func initReqInfo(c *gin.Context) {
	uuidV4 := uuid.New()
	c.Set("requestID", uuidV4)
	logger.Infof("Request(%s): %s %s", uuidV4, c.Request.Method, c.Request.URL.Path)
	c.Next()
}

func newRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(initReqInfo)
	router.Use(auth.VerifyReq())
	return router
}

func serviceDiscoveryAddr() string {
	return fmt.Sprintf(
		"%s:%d",
		apiConf.Data.Spec.Listen.Address.Advertise,
		apiConf.Data.Spec.Listen.Port,
	)
}

func localAddr() string {
	return fmt.Sprintf(
		"%s:%d",
		apiConf.Data.Spec.Listen.Local,
		apiConf.Data.Spec.Listen.Port,
	)
}

func initNodeInfo() {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	hostID, err := definition.GenerateNodeHashByMacAddr()
	if err != nil {
		panic(err)
	}

	role, err := cubecos.GetNodeRole()
	if err != nil {
		panic(err)
	}

	definition.HostID = hostID
	definition.Hostname = hostname
	definition.CurrentRole = role
	definition.IsGPUEnabled = cubecos.IsGPUEnabled()
}

func initLogger() (logger.Logger, error) {
	return log.NewCentralLogger(
		log.File(apiConf.Data.Spec.Log.File),
		log.Level(apiConf.Data.Spec.Log.Level),
		log.Backups(apiConf.Data.Spec.Log.Rotation.Backups),
		log.Size(apiConf.Data.Spec.Log.Rotation.Size),
		log.TTL(apiConf.Data.Spec.Log.Rotation.TTL),
		log.Compress(apiConf.Data.Spec.Log.Rotation.Compress),
	)
}

func genMetadata() map[string]string {
	return map[string]string{
		"hostname":     definition.Hostname,
		"nodeID":       definition.HostID,
		"isGPUEnabled": fmt.Sprintf("%t", definition.IsGPUEnabled),
	}
}

func newHttpServer() (*server.Server, error) {
	router := newRouter()
	err := RegisterHandlersByRole(router)
	if err != nil {
		logger.Errorf("failed to register handlers: %s", err.Error())
		return nil, err
	}

	srv := http.NewServer(
		server.Name(definition.CurrentRole),
		server.Metadata(genMetadata()),
		server.WithLogger(logger.DefaultLogger),
		server.Address(localAddr()),
		server.Advertise(serviceDiscoveryAddr()),
	)

	err = srv.Handle(srv.NewHandler(router))
	if err != nil {
		logger.Errorf("failed to new handler: %s", err.Error())
		return nil, err
	}

	return &srv, nil
}

func NewRuntime(conf config.Config) (*server.Server, error) {
	err := conf.Get().Scan(&apiConf.Data)
	if err != nil {
		logger.Errorf("failed to scan config: %s", err.Error())
		return nil, err
	}

	logger.DefaultLogger, err = initLogger()
	if err != nil {
		logger.Errorf("failed to init logger: %s", err.Error())
		return nil, err
	}

	showPromptMessage()
	initNodeInfo()
	return newHttpServer()
}
