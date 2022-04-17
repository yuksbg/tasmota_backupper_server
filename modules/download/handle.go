package download

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func HandleDownload(c *gin.Context) {
	fname := c.Param("filename")
	mac := c.Param("mac")

	c.FileAttachment(fmt.Sprintf("backup_data/%s/%s", mac, fname), fname)
}
