package main

import (
	"sync"
)

type state struct {
	sync.Mutex
	failureThreshold     uint
	failedCommunications uint
	total                uint
	counter              uint
	text                 string
	loading              bool
}

func (m *state) IncFailures() {
	m.Lock()
	defer m.Unlock()
	m.failedCommunications++
}

func (m *state) GetFailures() uint {
	m.Lock()
	defer m.Unlock()
	return m.failedCommunications
}

func (m *state) Set(value *ProgressReport) {
	m.Lock()
	defer m.Unlock()
	if value == nil {
		return
	}
	m.total = value.Total
	m.loading = value.Loading
	m.text = value.Text
	m.counter = value.Current
	m.failedCommunications = 0
}

func (m *state) IsLoading() bool {
	m.Lock()
	defer m.Unlock()
	return m.loading
}

func (m *state) StopDueToFailures() bool {
	m.Lock()
	defer m.Unlock()
	return m.failedCommunications > m.failureThreshold
}

func (m *state) GetText() string {
	m.Lock()
	defer m.Unlock()
	return m.text
}

func (m *state) GetPercentage() float64 {
	m.Lock()
	defer m.Unlock()

	if m.total == 0 {
		return 0.0
	}

	return float64(m.counter) / float64(m.total)
}
