package models

import (
	"time"

	"example.com/event-booker/db"
)

type Event struct {
	ID          int64     `json:"id"`
	Name        string    `binding:"required" json:"name"`
	Description string    `binding:"required" json:"description"`
	Location    string    `binding:"required" json:"location"`
	DateTime    time.Time `binding:"required" json:"dateTime"`
	CreatedAt   time.Time `json:"createdAt"`
	UserID      int64     `json:"userId"`
}

func (Event) tableName() string {
	return "events"
}

func (e Event) columnNames() []string {
	return getColumnNames(e)
}

func GetEvents() ([]Event, error) {
	query := "SELECT * FROM events"
	r, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var events = []Event{}

	for r.Next() {
		var e Event
		err := r.Scan(
			&e.ID,
			&e.Name,
			&e.Description,
			&e.Location,
			&e.DateTime,
			&e.CreatedAt,
			&e.UserID,
		)

		if err != nil {
			return nil, err
		}

		events = append(events, e)
	}

	return events, nil
}
