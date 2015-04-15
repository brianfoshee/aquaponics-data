package notify

import (
	"log"
	"strconv"
	"time"

	"github.com/brianfoshee/aquaponics-data/models"
)

type Manager struct {
	notifiers []Notifier
	Ch        chan *models.Reading
}

func NewManager() *Manager {
	m := &Manager{
		Ch: make(chan *models.Reading, 10),
	}
	return m
}

// Run is started in a goroutine
func (m *Manager) Run() {
	b := map[string]time.Time{}
	for r := range m.Ch {
		r := r
		t := time.Now()
		i, ok := b[r.Device.Identifier]
		a := r.ShouldAlert()
		if (a && !ok) || (a && ok && t.Sub(i) > 1*time.Minute) {
			b[r.Device.Identifier] = t
			tds := strconv.FormatFloat(r.SensorData.TDS, 'f', 2, 64)
			m.send("Error: Readings are off. TDS: " + tds)
		}
	}
}

func (m *Manager) Register(n Notifier) {
	m.notifiers = append(m.notifiers, n)
}

func (m Manager) send(msg string) {
	b := []byte(msg)
	for _, n := range m.notifiers {
		_, err := n.Write(b)
		if err != nil {
			log.Printf("Unable to notify", err)
		}
	}
}
