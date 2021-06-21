package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListAvailableSlots(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
func MakeAppointment(c *gin.Context) {
	c.Status(http.StatusAccepted)
}
func ConfirmAppointment(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
func ListVaccinationCenterDailyAppointments(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
