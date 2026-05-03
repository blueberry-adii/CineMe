package booking

type movieResponse struct {
	ID       string     `json:"id"`
	Title    string     `json:"title"`
	Genre    string     `json:"genre"`
	Duration string     `json:"duration"`
	Rating   string     `json:"rating"`
	Seats    [][]string `json:"seats"`
}

var Movies = []movieResponse{
	{
		ID:       "inception",
		Title:    "Inception",
		Genre:    "Sci-Fi • Action • Adventure",
		Duration: "148 min",
		Rating:   "8.8",
		Seats: [][]string{
			{"A1", "A2", "A3", "A4", "A5", "A6"},
			{"B1", "B2", "B3", "B4", "B5", "B6"},
			{"C1", "C2", "C3", "C4", "C5", "C6"},
			{"spacer"},
			{"D1", "D2", "D3", "D4", "D5", "D6"},
			{"E1", "E2", "E3", "E4", "E5", "E6", "E7"},
			{"F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8"},
		},
	},
	{
		ID:       "interstellar",
		Title:    "Interstellar",
		Genre:    "Sci-Fi • Adventure • Drama",
		Duration: "169 min",
		Rating:   "8.7",
		Seats: [][]string{
			{"A1", "A2", "A3", "A4", "A5", "A6"},
			{"B1", "B2", "B3", "B4", "B5", "B6"},
			{"C1", "C2", "C3", "C4", "C5", "C6"},
			{"spacer"},
			{"D1", "D2", "D3", "D4", "D5", "D6"},
			{"E1", "E2", "E3", "E4", "E5", "E6", "E7"},
			{"F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8"},
		},
	},
}
