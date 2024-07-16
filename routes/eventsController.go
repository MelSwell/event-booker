package routes

import (
	"net/http"
	"strconv"

	"example.com/event-booker/models"
	"github.com/gin-gonic/gin"
)

func getEvents(c *gin.Context) {
	e, err := models.GetAllEvents()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   e,
	})
}

func createEvent(c *gin.Context) {
	var e models.Event
	err := c.ShouldBindJSON(&e)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	err = e.Save()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   e,
	})
}

func getEvent(c *gin.Context) {
	e, err := getterHelper(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   e,
	})
}

func updateEvent(c *gin.Context) {
	e, err := getterHelper(c)
	if err != nil {
		return
	}

	err = c.ShouldBindJSON(&e)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	err = e.Update()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   e,
	})
}

func deleteEvent(c *gin.Context) {
	e, err := getterHelper(c)
	if err != nil {
		return
	}

	err = e.Delete()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"status": "success",
		"data":   nil,
	})
}

func getterHelper(c *gin.Context) (*models.Event, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Invalid id",
		})
		return nil, err
	}

	e, err := models.GetEventByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "Could not find an event with that ID",
		})
		return nil, err
	}

	return e, nil
}
