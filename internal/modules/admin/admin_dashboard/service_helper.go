package admindashboard

import "time"

func getRange(filter string)(time.Time,time.Time){

	now := time.Now()

	switch filter {

	case "week":
		return now.AddDate(0,0,-7), now

	case "month":
		return now.AddDate(0,-1,0), now

	case "year":
		return now.AddDate(-1,0,0), now

	default:
		return now.AddDate(0,-1,0), now
	}
}