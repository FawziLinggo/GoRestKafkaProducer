package main

import (
	"github.com/gin-gonic/gin/binding"
	"net/http"
	model "restAPIKafkaProducer/model"
	"restAPIKafkaProducer/producer"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	err := r.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		panic(err)
	}

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"admin": "password",
	}))
	authorized.POST("/post", func(c *gin.Context) {
		var data model.Data
		if err := c.ShouldBindWith(&data, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result := data
		c.JSON(http.StatusOK, gin.H{"result": result, "response": "success", "status": http.StatusOK})
		producer.Producer(data)
	})

	r.RunTLS("localhost:443", "ssl/server.crt", "ssl/server.key")
}
