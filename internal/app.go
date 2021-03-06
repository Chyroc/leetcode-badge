package internal

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func APP() *gin.Engine {
	fmt.Printf("config %#v\n", Conf)
	if Conf.Release {
		gin.SetMode(gin.ReleaseMode)
	}
	app := gin.Default()
	app.Use(cors.Default())

	app.GET("", func(c *gin.Context) {
		name := c.Query("name")
		if name == "" {
			c.String(400, "name is empty")
			return
		}

		logVisitor(name)

		leetcodeData, err := fetchLeetcodeData(name)
		if err != nil {
			c.String(400, err.Error())
			return
		} else if leetcodeData == nil {
			c.String(400, "fetch leetcode data of account: %s, result: nil", name)
			return
		}

		shieldsData, err := fetchShieldsData(c.Query("leetcode_badge_style"), leetcodeData)
		if err != nil {
			c.String(400, err.Error())
			return
		}

		c.Writer.WriteHeader(200)
		c.Writer.Header().Add("Content-Type", "image/svg+xml; charset=utf-8")
		c.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Type, Cache-Control, Expires")
		c.Writer.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0")
		c.Writer.Header().Add("Expires", "0")
		c.Writer.Header().Add("Pragma", "no-cache")
		c.Writer.WriteString(shieldsData)

		return
	})

	return app
}
