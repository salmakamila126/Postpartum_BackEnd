package gin

import "github.com/gin-gonic/gin"

func Start() *gin.Engine {
	r := gin.Default()
	return r
}
