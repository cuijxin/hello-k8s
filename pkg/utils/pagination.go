package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/unknwon/com"
)

func GetPage(c *gin.Context) int {
	result := 0
	page, _ := com.StrTo(c.Query("page")).Int()
	if page > 0 {
		result = (page - 1) * viper.GetInt("app.page_size")
	}
	return result
}
