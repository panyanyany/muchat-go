package gin_util

import (
	"net/http"

	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Validate(c *gin.Context, obj interface{}) bool {
	var err error
	validate := validator.New()

	if err = c.ShouldBind(obj); err != nil {
		seelog.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	err = validate.Struct(obj)
	if err != nil {
		seelog.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	return true
}
