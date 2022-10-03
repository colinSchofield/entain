package db

const (
	sportsList = "list"
	sportsGet  = "Get"
)

func getSportQueries() map[string]string {
	return map[string]string{
		sportsList: `
			SELECT 
				id, 
				meeting_id, 
				name, 
				number, 
				visible, 
				advertised_start_time 
			FROM sports
		`,
		sportsGet: `
			SELECT 
				id, 
				meeting_id, 
				name, 
				number, 
				visible, 
				advertised_start_time 
			FROM sports
			WHERE id = ?
		`,
	}
}
