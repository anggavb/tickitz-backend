package dto

type DashboardSalesChartRequest struct {
	MovieName string `form:"movie_name"`
	Period    string `form:"period"`
}

type DashboardTicketSalesRequest struct {
	Category string `form:"category"`
	Location string `form:"location"`
	Period   string `form:"period"`
}

type DashboardSalesPoint struct {
	Period  string `json:"period"`
	Revenue int64  `json:"revenue"`
}

type DashboardSalesData struct {
	MovieName string                `json:"movie_name,omitempty"`
	Period    string                `json:"period"`
	Points    []DashboardSalesPoint `json:"points"`
}

type DashboardSalesResponse struct {
	Success bool               `json:"success"`
	Data    DashboardSalesData `json:"data"`
}

type DashboardTicketSalesPoint struct {
	Period      string `json:"period"`
	TicketCount int64  `json:"ticket_count"`
	Revenue     int64  `json:"revenue"`
}

type DashboardTicketSalesResponse struct {
	Success bool                        `json:"success"`
	Data    []DashboardTicketSalesPoint `json:"data"`
}
