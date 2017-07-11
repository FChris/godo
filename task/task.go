package task

import "time"


type DayList []Day

func (t DayList) Len() int {
	return len(t)
}

func (t DayList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t DayList) Less(i, j int) bool {
	return t[i].Date.Before(t[j].Date)
}

func (t DayList) hasDate(date time.Time) bool {
	for _, d := range t {
		if date == d.Date {
			return true
		}
	}
	return false
}

//getDay returns the day for the give date from the DayList or a newly initialized day
//for the date if the list does not contain it yet
func (t DayList) GetDay(date time.Time) Day {
	for _, d := range t {
		if date == d.Date {
			return d
		}
	}

	return Day{date, make([]Todo, 0, 1)}
}


type Day struct {
	Date  time.Time
	Todos []Todo
}

type Todo struct {
	Description string
	Complete    bool
}

