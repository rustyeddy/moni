package moni

type LogCard struct {
	*Card
}

func NewLogCard(msg string) (c *Card) {
	c = NewCard("Log")
	c.Title = "Logs ... "
	c.Cols = "col-6"
	c.Text = msg
	return c
}
