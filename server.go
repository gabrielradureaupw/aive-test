package main

import "github.com/gin-gonic/gin"

func Serve() {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	appointments := r.Group("appointments")
	{
		appointments.GET("", ListAvailableSlots)
		appointments.POST("", MakeAppointment)
		appointments.PATCH("", ConfirmAppointment)
		appointments.GET("daily", gin.BasicAuth(gin.Accounts{
			"aive": "test",
		}), ListVaccinationCenterDailyAppointments)
	}

	r.Run()
}
