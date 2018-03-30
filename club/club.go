package club

//Club is the main component of this project
type Club struct {
	ClubID       string `json:"club_id"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	OpeningHours string `json:"openinghours"`
	Price        string `json:"price"`
}
