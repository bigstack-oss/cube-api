package runtime

import (
	"fmt"
	"os"

	"github.com/bigstack-oss/cube-api/internal/api"
	cubeConf "github.com/bigstack-oss/cube-api/internal/config"
	"github.com/bigstack-oss/cube-api/internal/cubeos"
	definition "github.com/bigstack-oss/cube-api/internal/definition/v1"
	"github.com/bigstack-oss/cube-api/internal/helpers/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/micro/plugins/v5/server/http"
	"go-micro.dev/v5/config"
	"go-micro.dev/v5/logger"
	"go-micro.dev/v5/server"
)

func RegisterHandlersToRolesByNodeRole(router *gin.Engine) error {
	var err error
	definition.CurrentRole, err = cubeos.GetNodeRole()
	if err != nil {
		logger.Errorf("failed to get node role: %s", err.Error())
		return err
	}

	groupHandlers := api.GetGroupHandlersByRole(definition.CurrentRole)
	if len(groupHandlers) == 0 {
		return fmt.Errorf("no handlers found for role: %s", definition.CurrentRole)
	}

	for _, handlers := range groupHandlers {
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
	return router
}

func serviceDiscoveryAddr() string {
	return fmt.Sprintf(
		"%s:%d",
		cubeConf.Conf.Spec.Access.Address.Advertise,
		cubeConf.Conf.Spec.Access.Port,
	)
}

func localAddr() string {
	return fmt.Sprintf(
		"%s:%d",
		cubeConf.Conf.Spec.Access.Local,
		cubeConf.Conf.Spec.Access.Port,
	)
}

func initNodeInfo() {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	definition.HostID = definition.GenerateNodeHashByMacAddr()
	definition.Hostname = hostname
	definition.IsGPUEnabled = cubeos.IsGPUEnabled()
}

func initLogger() (logger.Logger, error) {
	return log.NewCentralLogger(
		log.File(cubeConf.Conf.Spec.Log.File),
		log.Level(cubeConf.Conf.Spec.Log.Level),
		log.Backups(cubeConf.Conf.Spec.Log.Rotation.Backups),
		log.Size(cubeConf.Conf.Spec.Log.Rotation.Size),
		log.TTL(cubeConf.Conf.Spec.Log.Rotation.TTL),
		log.Compress(cubeConf.Conf.Spec.Log.Rotation.Compress),
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
	err := RegisterHandlersToRolesByNodeRole(router)
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
	err := conf.Get().Scan(&cubeConf.Conf)
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
