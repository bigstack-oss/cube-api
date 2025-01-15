package runtime

import (
	"fmt"
	"os"

	apihttp "github.com/bigstack-oss/bigstack-dependency-go/pkg/http"
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/log"
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/mongo"
	"github.com/bigstack-oss/cube-cos-api/internal/api"
	apitunings "github.com/bigstack-oss/cube-cos-api/internal/api/v1/tunings"
	apiConf "github.com/bigstack-oss/cube-cos-api/internal/config"
	"github.com/bigstack-oss/cube-cos-api/internal/controllers/v1/node"
	"github.com/bigstack-oss/cube-cos-api/internal/cubecos"
	definition "github.com/bigstack-oss/cube-cos-api/internal/definition/v1"
	"github.com/bigstack-oss/cube-cos-api/internal/keycloak"
	"github.com/bigstack-oss/cube-cos-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/micro/plugins/v5/server/http"
	"go-micro.dev/v5/config"
	"go-micro.dev/v5/logger"
	"go-micro.dev/v5/server"
)

func NewRuntime(conf config.Config) (*server.Server, error) {
	err := conf.Get().Scan(&apiConf.Data)
	if err != nil {
		logger.Errorf("failed to scan config: %s", err.Error())
		return nil, err
	}

	err = newGlobalLogHelper(apiConf.Data.Spec.Log)
	if err != nil {
		logger.Errorf("failed to init logger: %s", err.Error())
		return nil, err
	}

	err = newGlobalHttpHelper()
	if err != nil {
		logger.Errorf("failed to init http helper: %s", err.Error())
		return nil, err
	}

	err = newGlobalKeycloakAuth()
	if err != nil {
		logger.Errorf("failed to init keycloak auth: %s", err.Error())
		return nil, err
	}

	err = newGlobalMongoHelper(apiConf.Data.Spec.Store.MongoDB)
	if err != nil {
		logger.Errorf("failed to init mongo helper: %s", err.Error())
		return nil, err
	}

	initNodeIdentities()
	initNodeMemberSyncer()
	initNodeApiHandler()
	showPromptMessage()

	return newHttpServer()
}

func newGlobalLogHelper(opts log.Options) error {
	return log.NewGlobalHelper(
		log.File(opts.File),
		log.Level(opts.Level),
		log.Backups(opts.Rotation.Backups),
		log.Size(opts.Rotation.Size),
		log.TTL(opts.Rotation.TTL),
		log.Compress(opts.Rotation.Compress),
	)
}

func newGlobalMongoHelper(opts mongo.Options) error {
	return mongo.NewGlobalHelper(
		mongo.Uri(opts.Uri),
		mongo.AuthEnable(opts.Auth.Enable),
		mongo.AuthSource(opts.Auth.Source),
		mongo.AuthUsername(opts.Auth.Username),
		mongo.AuthPassword(opts.Auth.Password),
		mongo.ReplicaSet(opts.ReplicaSet),
	)
}

func newGlobalHttpHelper() error {
	return apihttp.NewGlobalHelper()
}

func newGlobalKeycloakAuth() error {
	return keycloak.NewGlobalSamlAuth(keycloak.Saml{
		IdentityProvider: keycloak.Provider{
			Host: keycloak.Host{
				Scheme:      "https",
				VirtualIp:   definition.ControllerVip,
				Port:        10443,
				InsecureTls: true,
			},
			MetadataPath: definition.DefaultIdpSamlMetadataPath,
		},
		ServiceProvider: keycloak.Provider{
			Host: keycloak.Host{
				Scheme:    "https",
				VirtualIp: definition.ControllerVip,
				Port:      8000,
				Auth: keycloak.Auth{
					Cert: definition.DefaultApiServerCert,
					Key:  definition.DefaultApiServerKey,
				},
			},
			MetadataPath: definition.DefaultSpSamlMetadataPath,
		},
	})
}

func initNodeIdentities() {
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

	vip, err := cubecos.GetControllerVirtualIp()
	if err != nil {
		panic(err)
	}

	isHaEnabled, err := cubecos.IsHaEnabled()
	if err != nil {
		panic(err)
	}

	definition.HostID = hostID
	definition.Hostname = hostname
	definition.CurrentRole = role
	definition.ControllerVip = vip
	definition.ListenAddr = localAddr()
	definition.AdvertiseAddr = serviceDiscoveryAddr()
	definition.IsHaEnabled = isHaEnabled
	definition.IsGpuEnabled = cubecos.IsGpuEnabled()
}

func initNodeMemberSyncer() {
	service.RegisterController(node.Name(), node.NewController())
}

func initNodeApiHandler() {
	api.RegisterHandlersToRoles(
		definition.Tunings,
		apitunings.Handlers,
		definition.RoleControl,
		definition.RoleCompute,
	)

	// Register other handlers here
	// ...
}

func newHttpServer() (*server.Server, error) {
	router := newRouter()
	err := registerHandlersByRole(router)
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

func newRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Any("/api/v1/saml/*action", gin.WrapH(keycloak.SamlAuth))
	router.Use(gin.Recovery())
	router.Use(initReqInfo)
	router.Use(adapter.Wrap(keycloak.SamlAuth.RequireAccount))
	return router
}

func initReqInfo(c *gin.Context) {
	uuidV4 := uuid.New()
	c.Set("requestID", uuidV4)
	logger.Infof("request(%s): %s %s", uuidV4, c.Request.Method, c.Request.URL.Path)
	c.Next()
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

func genMetadata() map[string]string {
	return map[string]string{
		"hostname":     definition.Hostname,
		"nodeID":       definition.HostID,
		"isGpuEnabled": fmt.Sprintf("%t", definition.IsGpuEnabled),
	}
}

func registerHandlersByRole(router *gin.Engine) error {
	groupHandlers := api.GetGroupHandlersByRole(definition.CurrentRole)
	if len(groupHandlers) == 0 {
		return fmt.Errorf("no handlers found for role(%s)", definition.CurrentRole)
	}

	for _, handlers := range groupHandlers {
		setGroupHandlersToRouter(router, handlers)
	}

	return nil
}

func setGroupHandlersToRouter(router *gin.Engine, handlers []api.Handler) {
	for _, h := range handlers {
		if h.Version == "" {
			logger.Warnf("skip invalid API registration: %s %s (no version provided)", h.Method, h.Path)
			continue
		}

		routerGroup := router.Group(h.Version)
		routerGroup.Handle(h.Method, h.Path, h.Func)
		logger.Infof("register API: %s %s", h.Method, fmt.Sprintf("%s%s", h.Version, h.Path))
	}
}
