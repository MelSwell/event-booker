package controllers

import (
	"net/http"
	"strconv"

	"example.com/event-booker/apperrors"
	"example.com/event-booker/middlewares"
	"example.com/event-booker/models"
	"github.com/gin-gonic/gin"
)

func GetEvents(c *gin.Context) {
	e, err := models.GetEvents()

	if err != nil {
		middlewares.SetError(c, apperrors.Validation{Message: err.Error()})
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
		middlewares.SetError(c, apperrors.Validation{Message: err.Error()})
		return
	}

	e.UserID = c.GetInt64("userId")
	id, err := models.Create(e)
	if err != nil {
		middlewares.SetError(c, apperrors.Validation{Message: err.Error()})
		return
	}

	// fetch created event back from DB in order to reflect default values in resp
	if err = models.GetByID(&e, id); err != nil {
		middlewares.SetError(c, apperrors.Internal{Message: "something went wrong"})
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   e,
	})
}

func GetEvent(c *gin.Context) {
	e, err := getEventByParam(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   e,
	})
}

func UpdateEvent(c *gin.Context) {
	e, err := getEventByParam(c)
	if err != nil {
		return
	}

	if err = c.ShouldBindJSON(&e); err != nil {
		middlewares.SetError(c, apperrors.Validation{Message: err.Error()})
		return
	}

	if e.UserID != c.GetInt64("userId") {
		middlewares.SetError(c, apperrors.Unauthorized{Message: "not authorized to update this event"})
		return
	}

	if err = models.Update(*e, e.ID); err != nil {
		middlewares.SetError(c, apperrors.Internal{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   e,
	})
}

func DeleteEvent(c *gin.Context) {
	e, err := getEventByParam(c)
	if err != nil {
		return
	}

	if e.UserID != c.GetInt64("userId") {
		middlewares.SetError(c, apperrors.Unauthorized{Message: "not authorized to delete this event"})
		return
	}

	if err = models.Delete(e, e.ID); err != nil {
		middlewares.SetError(c, apperrors.Internal{Message: err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"status": "success",
		"data":   nil,
	})
}

func getEventByParam(c *gin.Context) (*models.Event, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		middlewares.SetError(c, apperrors.Validation{Message: "Invalid Id"})
		return nil, err
	}

	var e models.Event
	if err = models.GetByID(&e, id); err != nil {
		middlewares.SetError(c, apperrors.NotFound{Message: "Could not find event with that ID"})
		return nil, err
	}
	return &e, nil
}
