package speedsensor

import (
	"errors"
	"fmt"
	Log "github.com/kevinchapron/BasicLogger/Logging"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"math"
	"sync"
	"time"
)

var mutex sync.Mutex

type ADSxValue struct {
	Value float64
	Error error
}

func (v *ADSxValue) String() string {
	if v == nil || v.Error != nil {
		return "---"
	}
	return fmt.Sprintf("%.3f", v.Value)
}
func (v *ADSxValue) GetValue() float64 {
	if v.Value < IR_SENSOR_MIN_DISTANCE || v.Value > IR_SENSOR_MAX_DISTANCE {
		return 0
	}
	return v.Value
}

type ADSxDATA struct {
	Values    map[uint]*ADSxValue
	Timestamp time.Time
}

type Bot struct {
	Adaptor *raspi.Adaptor
	Driver  *i2c.ADS1x15Driver

	Data     chan *ADSxDATA
	callback *func(*ADSxDATA)

	waiter chan bool

	running   bool
	frequency int
}

func (b *Bot) Start() error {
	err := b.Driver.Start()
	if err == nil {
		b.running = true
		b.Run()
	}
	return err
}
func (b *Bot) Stop() error {
	err := b.Driver.Halt()
	if err == nil {
		b.running = false
	}
	return err
}

func (b *Bot) Run() {
	go b.RunCallback()
	timer := time.NewTicker(time.Second / FREQUENCY)
	var i uint
	for {
		<-timer.C
		go func() {
			a := ADSxDATA{
				Values:    make(map[uint]*ADSxValue),
				Timestamp: time.Now(),
			}
			for i = 0; i < 3; i++ {
				v, err := b.Driver.Read(int(i), b.Driver.DefaultGain, 475)
				if err != nil {
					Log.Error(err)
				}
				a.Values[i] = b.VoltageToDistance(v)
			}
			b.Data <- &a
		}()
	}
}
func (b *Bot) RunCallback() {
	for {
		(*b.callback)(<-b.Data)
	}
}
func (b *Bot) SetFrequency(f int) {
	b.frequency = f
}
func (b *Bot) VoltageToDistance(f float64) *ADSxValue {
	v := 35.274*math.Pow(f, 4) - 244.1*math.Pow(f, 3) + 618.21*math.Pow(f, 2) - 704.42*f + 353.5
	return NewADSxValue(v)
}

func (b *Bot) SetCallback(f *func(*ADSxDATA)) {
	b.callback = f
}

func NewADSxValue(v float64) *ADSxValue {
	vr := ADSxValue{
		Value: v,
		Error: nil,
	}
	if v < IR_SENSOR_MIN_DISTANCE || v > IR_SENSOR_MAX_DISTANCE {
		vr.Error = errors.New("Value not in range " + string(IR_SENSOR_MIN_DISTANCE) + "-" + string(IR_SENSOR_MAX_DISTANCE))
	}
	return &vr
}
func CreateBot() *Bot {
	ada := raspi.NewAdaptor()
	dr := i2c.NewADS1115Driver(ada, i2c.WithBus(1), i2c.WithAddress(0x48))

	return &Bot{
		Adaptor: ada,
		Driver:  dr,

		Data: make(chan *ADSxDATA),
	}
}
