package dto

type MovieRequest struct {
	Name             string   `json:"name" form:"name" binding:"required"`
	ReleaseDate      string   `json:"release_date" form:"release_date" binding:"required"`
	DurationInMinute int      `json:"duration_in_minute" form:"duration_in_minute" binding:"required,gt=0"`
	DirectorName     string   `json:"director_name,omitempty" form:"director_name"`
	Synopsis         string   `json:"synopsis,omitempty" form:"synopsis"`
	Image            string   `json:"image,omitempty" form:"image"`
	Categories       []string `json:"categories,omitempty" form:"categories"`
	Casts            []string `json:"casts,omitempty" form:"cast"`
}

type MovieResponse struct {
	ID               int64    `json:"id"`
	Name             string   `json:"name"`
	ReleaseDate      string   `json:"release_date"`
	DurationInMinute int      `json:"duration_in_minute"`
	DirectorName     string   `json:"director_name,omitempty"`
	Synopsis         string   `json:"synopsis,omitempty"`
	Image            string   `json:"image,omitempty"`
	Categories       []string `json:"categories,omitempty"`
	Casts            []string `json:"casts,omitempty"`
	CreatedAt        string   `json:"created_at,omitempty"`
	UpdatedAt        string   `json:"updated_at,omitempty"`
}

// MovieListResponse is the wrapper for list responses
type MovieListResponse struct {
	Success    bool            `json:"success"`
	Data       []MovieResponse `json:"data"`
	Pagination interface{}     `json:"pagination"`
}

// MovieSingleResponse is the wrapper for single movie responses
type MovieSingleResponse struct {
	Success bool          `json:"success"`
	Data    MovieResponse `json:"data"`
}
