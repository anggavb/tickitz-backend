package dto

type SeatMapWrappedResponse struct {
	Status string          `json:"status" example:"success"`
	Data   SeatMapResponse `json:"data"`
}

type SeatMapResponse struct {
	MovieCinema MovieCinemaResponse `json:"movie_cinema"`
	Movie       SeatMovie           `json:"movie"`
	Cinema      SeatCinema          `json:"cinema"`
	Layout      SeatLayout          `json:"layout"`
	Seats       []SeatResponse      `json:"seats"`
}

type MovieCinemaResponse struct {
	ID         int64  `json:"id"`
	Date       string `json:"date"`
	Time       string `json:"time"`
	Price      int    `json:"price"`
	ShowtimeID int64  `json:"showtime_id"`
}

type SeatMovie struct {
	ID     int64    `json:"id"`
	Title  string   `json:"title"`
	Poster string   `json:"poster"`
	Genres []string `json:"genres"`
}

type SeatCinema struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Logo     string `json:"logo"`
	Location string `json:"location"`
}

type SeatLayout struct {
	Rows    []string `json:"rows"`
	Columns []int    `json:"columns"`
}

type SeatResponse struct {
	ID     int64  `json:"id"`
	Code   string `json:"code"`
	Row    string `json:"row"`
	Number int    `json:"number"`
	Type   string `json:"type"`
	Price  int    `json:"price"`
	Status string `json:"status"`
}
