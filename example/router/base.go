package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tuya/tuya-connector-go/example/service"
)

func NewGinEngin() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	initRouter(engine)
	return engine
}

func initRouter(r *gin.Engine) {
	deviceGroup := r.Group("/devices")
	deviceGroup.GET("/:device_id", service.GetDevice)
	deviceGroup.PUT("/:device_id", service.PutDevice)
	deviceGroup.GET("/:device_id/functions", service.GetDeviceFunc)
	deviceGroup.POST("/:device_id/commands", service.PostDeviceCmd)
}
