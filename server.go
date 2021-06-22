package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Serve() {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	appointments := r.Group("appointments")
	{
		h := NewHandler(NewStore())
		appointments.GET("", h.ListAvailableSlots)
		appointments.POST("", h.MakeAppointment)
		appointments.GET("confirm", h.ConfirmAppointment)
		appointments.GET("daily", // use GET http method for simple use with email link
			func(c *gin.Context) {
				fmt.Println(c.GetHeader("Authorization"))
			},
			gin.BasicAuth(gin.Accounts{
				"aive": "test",
			}),
			h.ListDailyAppointments)
	}

	r.Run()
}
