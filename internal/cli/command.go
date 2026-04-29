package cli

type Command struct {
	Name string
	Args []string
}

type HandlerFunc func(*State, Command) error
