package confirm

type confirmationState int

const (
	unconfirmed confirmationState = iota
	confirmedOnce
	confirmedAll
)

type Confirmation struct {
	state confirmationState
}

func (c *Confirmation) Next() bool {
	switch c.state {
	case confirmedOnce:
		c.state = unconfirmed
		return true
	case confirmedAll:
		return true
	default:
		return false
	}
}

func (c *Confirmation) ConfirmOnce() {
	c.state = confirmedOnce
}

func (c *Confirmation) ConfirmAll() {
	c.state = confirmedAll
}
