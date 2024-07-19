package commons

type InitGodPower struct {
	Description string
	TokenCost   int
}

type InitGod struct {
	Emoji       string
	Name        string
	Description string
	Priority    int
	Levels      [3]InitGodPower
}
