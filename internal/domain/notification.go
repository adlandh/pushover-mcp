package domain

type Notification struct {
	Priority *int
	Message  string
	Title    string
	Sound    string
	URL      string
	URLTitle string
	Device   string
}
