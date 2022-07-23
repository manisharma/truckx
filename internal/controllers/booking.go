package controllers

import (
	"net/http"
	"strconv"
	"time"
	"truckx/internal/models"
	"truckx/internal/repository"
	"truckx/internal/services/util"

	"github.com/gin-gonic/gin"
)

type bookingCtrl struct {
	movieRepo   repository.MovieRepo
	bookingRepo repository.BookingRepo
}

func NewBookingCtrl(movieRepo repository.MovieRepo, bookingRepo repository.BookingRepo) *bookingCtrl {
	return &bookingCtrl{movieRepo, bookingRepo}
}

func (ctrl *bookingCtrl) Create(c *gin.Context) {
	var user *models.User
	if usr, ok := c.Get(string(models.UserCtxKey)); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not valid"})
		c.Abort()
		return
	} else {
		user = usr.(*models.User)
	}

	createBooking := &models.CreateBooking{}
	err := c.ShouldBindJSON(createBooking)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	createBooking.UserId = int(user.Id)
	createBooking.Date, err = time.Parse(util.DateFormat, createBooking.DateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	avalibility, err := ctrl.movieRepo.GetAvailability(c.Request.Context(), createBooking.MovieId, createBooking.State, createBooking.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	canBook := true
	for _, movieShowDetail := range avalibility[0].MovieShowDetail {
		if movieShowDetail.MovieShowId == createBooking.MovieShowId {
			if len(movieShowDetail.Seats) > len(createBooking.Seats) {
				canBook = false
				break
			}
			seatsToBeBooked := []int{}
			for _, seat := range createBooking.Seats {
				seatsToBeBooked = append(seatsToBeBooked, int(seat.SeatId))
			}
			for _, seat := range movieShowDetail.Seats {
				if ok, _ := util.Contains(seatsToBeBooked, seat); ok {
					canBook = false
					break
				}
			}
		}
	}

	if canBook {
		id, err := ctrl.bookingRepo.Create(c.Request.Context(), *createBooking)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{"bookingId": id})
	}
	c.JSON(http.StatusNotAcceptable, gin.H{"error": "booking can not be done due to unavalability"})
	c.Abort()
}

func (ctrl *bookingCtrl) Get(c *gin.Context) {
	bookingIdStr, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "booking id missing"})
		c.Abort()
		return
	}
	bookingId, err := strconv.Atoi(bookingIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		c.Abort()
		return
	}
	bookingDetail, err := ctrl.bookingRepo.Get(c.Request.Context(), bookingId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"bookingDetail": bookingDetail})
}

func (ctrl *bookingCtrl) List(c *gin.Context) {
	var user *models.User
	if usr, ok := c.Get(string(models.UserCtxKey)); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not valid"})
		c.Abort()
		return
	} else {
		user = usr.(*models.User)
	}
	bookingDetail, err := ctrl.bookingRepo.List(c.Request.Context(), int(user.Id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"bookingDetail": bookingDetail})
}
