package fault

// Message is a type of error that is just string message.
// this type can be used to create error constants instead of variables.
// See: https://dave.cheney.net/2016/04/07/constant-errors
type Message string

func (m Message) Error() string {
	return string(m)
}
