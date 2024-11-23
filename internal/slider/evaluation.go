package slider

func (m Model) Val() int { return m.current }

func (m *Model) Set(v int) {
	m.current = v
	m.fixRange()
}

func (m *Model) Inc(v int) {
	m.current += v
	m.fixRange()
}

func (m *Model) Dec(v int) {
	m.current -= v
	m.fixRange()
}

func (m *Model) Pcnt() float64 {
	return float64(m.current) / float64(m.max)
}

func (m *Model) SetPcnt(p float64) {
	m.current = int(float64(m.max) * p)
	m.fixRange()
}

func (m *Model) IncPcnt(p float64) {
	m.current += int(float64(m.max) * p)
	m.fixRange()
}

func (m *Model) DecPcnt(p float64) {
	m.current -= int(float64(m.max) * p)
	m.fixRange()
}

func (m *Model) fixRange() {
	if m.current > m.max {
		m.current = m.max
	}
	if m.current < 0 {
		m.current = 0
	}
}
