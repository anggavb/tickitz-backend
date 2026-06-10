package controller

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/model"
	"github.com/tickitz-backend/internal/service"
)

type SeatController struct {
	seatService *service.SeatService
}

func NewSeatController(seatService *service.SeatService) *SeatController {
	return &SeatController{seatService: seatService}
}

// GetSeatMap godoc
// @Summary      Get seat map by movie cinema ID
// @Description  Retrieve movie, cinema, layout, and seat availability for one exact movie cinema schedule.
// @Tags         Seats
// @Accept       json
// @Produce      json
// @Param        movie_cinema_id  path      int  true  "Movie Cinema ID"
// @Success      200              {object}  dto.SeatMapWrappedResponse
// @Failure      400              {object}  dto.ErrorResponse
// @Failure      404              {object}  dto.ErrorResponse
// @Failure      500              {object}  dto.ErrorResponse
// @Router       /movie-cinemas/{movie_cinema_id}/seats [get]
func (c *SeatController) GetSeatMap(ctx *gin.Context) {
	movieCinemaID, err := strconv.ParseInt(ctx.Param("movie_cinema_id"), 10, 64)
	if err != nil || movieCinemaID <= 0 {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status:  "error",
			Message: "Invalid movie cinema ID",
		})
		return
	}

	seatMap, err := c.seatService.GetSeatMap(ctx.Request.Context(), movieCinemaID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Status:  "error",
				Message: "Movie cinema schedule not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Status:  "error",
			Message: "Failed to get seat map",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.SeatMapWrappedResponse{
		Status: "success",
		Data:   buildSeatMapResponse(seatMap),
	})
}

func buildSeatMapResponse(seatMap model.SeatMap) dto.SeatMapResponse {
	seats := make([]dto.SeatResponse, 0, len(seatMap.Seats))
	for _, seat := range seatMap.Seats {
		seats = append(seats, dto.SeatResponse{
			ID:     seat.ID,
			Code:   seat.Code,
			Row:    seat.Row,
			Number: seat.Number,
			Type:   seat.Type,
			Price:  seat.Price,
			Status: seat.Status,
		})
	}

	return dto.SeatMapResponse{
		MovieCinema: dto.MovieCinemaResponse{
			ID:         seatMap.MovieCinema.ID,
			Date:       seatMap.MovieCinema.ShowDate.Format("2006-01-02"),
			Time:       seatMap.MovieCinema.Showtime,
			Price:      seatMap.MovieCinema.Price,
			ShowtimeID: seatMap.MovieCinema.ShowtimeID,
		},
		Movie: dto.SeatMovie{
			ID:     seatMap.Movie.ID,
			Title:  seatMap.Movie.Title,
			Poster: seatMap.Movie.Poster,
			Genres: seatMap.Movie.Genres,
		},
		Cinema: dto.SeatCinema{
			ID:       seatMap.Cinema.ID,
			Name:     seatMap.Cinema.Name,
			Logo:     seatMap.Cinema.Logo,
			Location: seatMap.Cinema.Location,
		},
		Layout: dto.SeatLayout{
			Rows:    seatMap.Layout.Rows,
			Columns: seatMap.Layout.Columns,
		},
		Seats: seats,
	}
}
