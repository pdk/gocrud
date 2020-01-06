package example

import "time"

type Account struct {
	ID              int64     `db:"account_id"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
	CreatedByUserID int64     `db:"created_by_user_id"`
	UpdatedByUserID int64     `db:"updated_by_user_id"`
	Name            string    `db:"name"`
	URLStub         string    `db:"url_stub"`
}
