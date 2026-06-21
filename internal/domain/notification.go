package domain

type Notification struct {
	Priority *int
	Retry    *int // Optional; defaults to 60 for emergency priority (2)
	Expire   *int // Optional; defaults to 3600 for emergency priority (2)
	Message  string
	Title    string
	Sound    string
	URL      string
	URLTitle string
	Device   string
}
