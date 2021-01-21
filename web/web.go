package web

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/rokoga/filas-backend/infra"
	"github.com/rokoga/filas-backend/service"
	"github.com/rokoga/filas-backend/vo"
)

// Run implements the main function of web API
func Run(done chan string) {
	const PORT = ":8080"
	router := gin.Default()
	router.Use(cors.Default())

	dbClient, dbCollection, err := infra.GetConnection("config/dev/.env")
	if err != nil {
		panic(err)
	}
	defer infra.CloseConnection(dbClient)

	svc := service.NewStoreServiceImpl(dbCollection)

	router.PUT("/store", func(c *gin.Context) {
		createRequest := vo.CreateRequest{}
		c.BindJSON(&createRequest)

		store, err := svc.Create(createRequest.Name)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, store)
	})

	router.DELETE("/store/:storeid", func(c *gin.Context) {
		id := c.Param("storeid")

		err := svc.RemoveStore(id)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, nil)
	})

	router.GET("/stores", func(c *gin.Context) {

		stores, err := svc.GetAllStores()
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, stores)
	})

	router.GET("/store/name/:name", func(c *gin.Context) {
		name := c.Param("name")

		domainStore, err := svc.GetStore(name)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, domainStore)
	})

	router.GET("/store/id/:id", func(c *gin.Context) {
		id := c.Param("id")

		domainStore, err := svc.GetStoreByID(id)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, domainStore)
	})

	router.PUT("/consumer", func(c *gin.Context) {
		addConsumerRequest := vo.AddConsumerRequest{}
		c.BindJSON(&addConsumerRequest)

		accessURL, err := svc.AddConsumer(addConsumerRequest.StoreID, addConsumerRequest.Name, addConsumerRequest.Phone)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, accessURL)
	})

	router.DELETE("/consumer/:storeid/:number", func(c *gin.Context) {
		storeid := c.Param("storeid")
		number := c.Param("number")

		err := svc.RemoveConsumer(storeid, number)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, nil)
	})

	router.GET("/consumer/:storeid/:number", func(c *gin.Context) {
		storeid := c.Param("storeid")
		number := c.Param("number")

		position, consumer, err := svc.GetConsumer(storeid, number)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response := gin.H{
			"position":  position,
			"name":      consumer.Name,
			"phone":     consumer.Phone,
			"accessKey": consumer.Name,
		}

		c.JSON(200, response)
	})

	router.GET("/consumers/:storeid", func(c *gin.Context) {
		storeid := c.Param("storeid")

		allConsumers, err := svc.GetAllConsumers(storeid)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, allConsumers)
	})

	fmt.Printf("Server is listening at %s", PORT)
	router.Run(PORT)

	done <- "Server shutdown"
}
