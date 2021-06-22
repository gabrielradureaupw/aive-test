package main

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	st AppointmentStore
}

func NewHandler(st AppointmentStore) Handler {
	return Handler{
		st: st,
	}
}

func (h Handler) ListAvailableSlots(c *gin.Context) {
	centersWithSlots, err := ListAvailableSlots(c, h.st)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"vaccinationCenters": centersWithSlots})
}

func (h Handler) MakeAppointment(c *gin.Context) {
	req := MakeAppointmentRequest{}
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if _, err := govalidator.ValidateStruct(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err := MakeAppointment(c, h.st, req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusAccepted)
}

func (h Handler) ConfirmAppointment(c *gin.Context) {
	req := ConfirmAppointmentRequest{}
	if err := c.BindQuery(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if _, err := govalidator.ValidateStruct(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err := ConfirmAppointment(c, h.st, req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "RDV confirmé"})
}

func (h Handler) ListDailyAppointments(c *gin.Context) {
	resp, err := ListDailyAppointments(c, h.st)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
