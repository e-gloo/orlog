package commons

type Player struct {
	Name     string `json:"name"`
	Health   int    `json:"health"`
	Token    int    `json:"token"`
	Dices    [6]Die `json:"dices"`
	Position int    `json:"position"`
	// gods []God
}

func (p *Player) RollDices() {
	for idx := range p.Dices {
		if p.Dices[idx].Kept == false {
			p.Dices[idx].Roll()
		}
	}
}

func (p *Player) UnkeepDices() {
	for idx := range p.Dices {
		p.Dices[idx].Kept = false
	}
}

func (p *Player) AttackPlayer(player *Player) {
	nbArrows := 0
	nbAxes := 0
	for _, die := range p.Dices {
		if die.Face().Kind == Arrow {
			nbArrows++
		}
		if die.Face().Kind == Axe {
			nbAxes++
		}
	}

	nbHelmets := 0
	nbShields := 0
	for _, die := range player.Dices {
		if die.Face().Kind == Helmet {
			nbHelmets++
		}
		if die.Face().Kind == Shield {
			nbShields++
		}
	}
	arrowDamage := Max(nbArrows-nbShields, 0)
	axeDamage := Max(nbAxes-nbHelmets, 0)

	player.Health = player.Health - arrowDamage - axeDamage
}

func (p *Player) GainTokens() {
	nbMagics := 0

	for _, die := range p.Dices {
		if die.Face().Magic == true {
			nbMagics++
		}
	}

	p.Token += nbMagics
}

func (p *Player) StealTokens(player *Player) {
	nbThieves := 0

	for _, die := range p.Dices {
		if die.Face().Kind == Thief {
			nbThieves++
		}
	}

	nbToken := Min(nbThieves, player.Token)
	p.Token += nbToken
	player.Token -= nbToken
}
