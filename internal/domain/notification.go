package domain

type Notification struct {
	Title    *string
	Priority *int
	Sound    *string
	URL      *string
	URLTitle *string
	Device   *string
	Message  string
}
