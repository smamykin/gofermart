package money

type IntMoney int

func FromFloat(rubles float64) IntMoney {
	return IntMoney(rubles * 100)
}

func (m *IntMoney) AsFloat() float64 {
	cents := float64(*m)

	return cents / 100.
}
