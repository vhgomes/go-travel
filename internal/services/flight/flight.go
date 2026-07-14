package flight

type Flight struct {
	ID             string `json:"id"`
	Origin         string `json:"origin"`
	Destination    string `json:"destination"`
	Date           string `json:"date"`
	FlightPrice    int64  `json:"flight_price"` // Preço em centavos
	AvailableSeats int    `json:"available_seats"`
}
