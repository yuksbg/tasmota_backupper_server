package files

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

func HandleDelete(c *gin.Context) {

	os.Remove(fmt.Sprintf("backup_data/%s/%s", c.Param("mac"), c.Param("filename")))

	c.JSON(200, gin.H{"result": true})
}
