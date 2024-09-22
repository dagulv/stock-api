package mailer

type Email struct {
	FirstName    string
	LastName     string
	Email        string
	Subject      string
	HTMLTemplate *string
	Vars         map[string]string
}
