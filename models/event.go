package models

import (
	"time"

	"example.com/event-booker/db"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int64
}

func (e Event) Save() error {
	query := `
	INSERT INTO events (name, description, location, dateTime, user_id)
	VALUES (?, ?, ?, ?, ?)`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		e.Name,
		e.Description,
		e.Location,
		e.DateTime,
		e.UserID,
	)

	if err != nil {
		return err
	}

	return err
}

func GetAllEvents() ([]Event, error) {
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
			&e.UserID,
		)

		if err != nil {
			return nil, err
		}

		events = append(events, e)
	}

	return events, nil
}

func GetEventByID(id int64) (*Event, error) {
	query := "SELECT * FROM events WHERE id = ?"
	r := db.DB.QueryRow(query, id)

	var e Event
	err := r.Scan(
		&e.ID,
		&e.Name,
		&e.Description,
		&e.Location,
		&e.DateTime,
		&e.UserID,
	)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (e Event) Update() error {
	query := `
	UPDATE events
	SET name = ?, description = ?, location=?, dateTime = ?
	WHERE id = ?
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.ID)
	return err
}

func (e Event) Delete() error {
	query := "DELETE FROM events WHERE id = ?"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.ID)
	return err
}
