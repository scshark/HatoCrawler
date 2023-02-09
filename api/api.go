package api

import (
	"HatoCrawler/api/controllers"
	"HatoCrawler/config"
	"HatoCrawler/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RouterInit(r *gin.Engine) {
	setCors(r)
	api := r.Group("/api", auth)

	public := api.Group("/crawler")
	{
		public.GET("/getArticles/:site", controllers.GetArticles)
	}
}

func setCors(r *gin.Engine) {
	conf := cors.DefaultConfig()
	conf.AllowAllOrigins = true
	conf.AllowHeaders = append(conf.AllowHeaders, "Authorization")
	r.Use(cors.New(conf))
}

func auth(c *gin.Context) {
	key := c.GetHeader("Authorization")

	if key != config.Cfg.Api.Auth {
		utils.ErrorStrResp(c, utils.INVALID_AUTH_KEY, "Invalid auth key")
		return
	}
	c.Next()
}
