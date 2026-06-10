package model

import "time"

type SeatMap struct {
	MovieCinema MovieCinemaInfo `json:"movie_cinema"`
	Movie       SeatMapMovie    `json:"movie"`
	Cinema      SeatMapCinema   `json:"cinema"`
	Layout      SeatLayout      `json:"layout"`
	Seats       []SeatItem      `json:"seats"`
}

type MovieCinemaInfo struct {
	ID         int64     `json:"id"`
	ShowDate   time.Time `json:"show_date"`
	Showtime   string    `json:"showtime"`
	Price      int       `json:"price"`
	ShowtimeID int64     `json:"showtime_id"`
}

type SeatMapMovie struct {
	ID     int64    `json:"id"`
	Title  string   `json:"title"`
	Poster string   `json:"poster"`
	Genres []string `json:"genres"`
}

type SeatMapCinema struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Logo     string `json:"logo"`
	Location string `json:"location"`
}

type SeatLayout struct {
	Rows    []string `json:"rows"`
	Columns []int    `json:"columns"`
}

type SeatItem struct {
	ID     int64  `json:"id"`
	Code   string `json:"code"`
	Row    string `json:"row"`
	Number int    `json:"number"`
	Type   string `json:"type"`
	Price  int    `json:"price"`
	Status string `json:"status"`
}
