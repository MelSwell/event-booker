package controllers

import (
	"net/http"
	"strconv"

	"example.com/event-booker/models"
	"github.com/gin-gonic/gin"
)

func GetEvents(c *gin.Context) {
	e, err := models.GetEvents()

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

func CreateEvent(c *gin.Context) {
	var e models.Event

	if err := c.ShouldBindJSON(&e); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	e.UserID = c.GetInt64("userId")
	id, err := models.Create(e)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}
	e.ID = id

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   e,
	})
}

func GetEvent(c *gin.Context) {
	e, err := getEventByID(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   e,
	})
}

func UpdateEvent(c *gin.Context) {
	e, err := getEventByID(c)
	if err != nil {
		return
	}

	if err = c.ShouldBindJSON(&e); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	if e.UserID != c.GetInt64("userId") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "not authorized to update this event",
		})
		return
	}

	if err = models.Update(*e, e.ID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   e,
	})
}

func DeleteEvent(c *gin.Context) {
	e, err := getEventByID(c)
	if err != nil {
		return
	}

	if e.UserID != c.GetInt64("userId") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "not authorized to delete this event",
		})
		return
	}

	if err = models.Delete(e, e.ID); err != nil {
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

func getEventByID(c *gin.Context) (*models.Event, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Invalid id",
		})
		return nil, err
	}

	var e models.Event
	if err = models.GetByID(&e, id); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "Could not find an event with that ID",
		})
		return nil, err
	}
	return &e, nil
}
