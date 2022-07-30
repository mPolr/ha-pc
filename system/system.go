package system

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
	"syscall"
)

func PostSystem(c *gin.Context) {
	action := c.Param("action")
	switch action {
	case "shutdown":
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "ok",
			"message": "shutdown",
		})
		//system.DoShutdown(c)
	case "reboot":
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "ok",
			"message": "reboot",
		})
		//system.DoReboot(c)
	default:
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Not found",
		})
		return
	}
}

func DoShutdown(c *gin.Context) {
	switch runtime.GOOS {
	case "windows":
		fmt.Print("windows")
	case "linux":
		syscall.Sync()
		err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
	default:
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Unsupported OS",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Shutting down " + runtime.GOOS,
	})
}

func DoReboot(c *gin.Context) {
	switch runtime.GOOS {
	case "windows":
		fmt.Print("windows")
	case "linux":
		syscall.Sync()
		err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
	default:
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Unsupported OS",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Rebooting " + runtime.GOOS,
	})
}
