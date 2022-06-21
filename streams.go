package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Stream struct {
	Id        int    `json:"streamId"`
	IsValid   int    `json:"isValidStream"`
	StreamKey string `json:"streamKey"`
}

type StreamController struct {
	DB *DatabaseUtility
}

func (sc *StreamController) VerifyStream(c *gin.Context) {
	key := c.DefaultPostForm("name", "")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			SUCCESS_KEY:          false,
			RESPONSE_MESSAGE_KEY: "Something went wrong getting the stream key!",
		})
		return
	}
	result, err := sc.DB.ToggleStream(1, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{SUCCESS_KEY: false, RESPONSE_MESSAGE_KEY: err.Error()})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	isValidKey := rowsAffected == 1
	switch {
	case isValidKey:
		c.JSON(http.StatusOK, gin.H{
			SUCCESS_KEY:          true,
			RESPONSE_MESSAGE_KEY: "Stream key is good!",
		})
	case !isValidKey:
		c.JSON(http.StatusNotFound, gin.H{
			SUCCESS_KEY:          false,
			RESPONSE_MESSAGE_KEY: "No stream key here!",
		})
	}
}

func (sc *StreamController) GetStreams(c *gin.Context) {
	streams := make([]Stream, 0)
	rows, err := sc.DB.GetStreams()
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			c.JSON(http.StatusNotFound, gin.H{SUCCESS_KEY: false, RESPONSE_MESSAGE_KEY: err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{SUCCESS_KEY: false, RESPONSE_MESSAGE_KEY: err.Error()})
		}
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id, isValid int
			streamKey   string
		)
		err := rows.Scan(&id, &isValid, &streamKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{SUCCESS_KEY: false, RESPONSE_MESSAGE_KEY: err.Error()})
			return
		}
		stream := Stream{id, isValid, streamKey}
		streams = append(streams, stream)
	}
	c.JSON(http.StatusOK, gin.H{SUCCESS_KEY: true, "streams": streams})
}

func (sc *StreamController) CreateKey(c *gin.Context) {
	guid := generateGUID()
	_, err := sc.DB.CreateStream(guid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			SUCCESS_KEY:          false,
			RESPONSE_MESSAGE_KEY: fmt.Sprintf("Failed to insert new record, %s", err.Error()),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		SUCCESS_KEY:          false,
		RESPONSE_MESSAGE_KEY: "Insert new stream ok!",
	})
}

func (sc *StreamController) DeleteStream(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			SUCCESS_KEY:          false,
			RESPONSE_MESSAGE_KEY: "Something went wrong getting the stream id!",
		})
		return
	}
	_, err := sc.DB.DeleteStream(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			SUCCESS_KEY:          false,
			RESPONSE_MESSAGE_KEY: fmt.Sprintf("Failed to delete existing record, %s", err.Error()),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		SUCCESS_KEY:          true,
		RESPONSE_MESSAGE_KEY: "Deleted stream!",
	})
}

func (sc *StreamController) EndStream(c *gin.Context) {
	key := c.DefaultPostForm("name", "")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			SUCCESS_KEY:          false,
			RESPONSE_MESSAGE_KEY: "Something went wrong getting the stream key!",
		})
		return
	}
	result, _ := sc.DB.ToggleStream(0, key)
	rowsAffected, _ := result.RowsAffected()
	isValidKey := rowsAffected == 1
	switch {
	case isValidKey:
		c.JSON(http.StatusOK, gin.H{
			SUCCESS_KEY:          true,
			RESPONSE_MESSAGE_KEY: "Stream key is good!",
		})
	case !isValidKey:
		c.JSON(http.StatusNotFound, gin.H{
			SUCCESS_KEY:          false,
			RESPONSE_MESSAGE_KEY: "No stream key here!",
		})
	}
}
