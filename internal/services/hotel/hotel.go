package hotel

type Hotel struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Location       string `json:"location"`
	Date           string `json:"date"`
	HotelPrice     int64  `json:"hotel_price"` // Preço em centavos
	RoomsAvailable int    `json:"rooms_available"`
}
