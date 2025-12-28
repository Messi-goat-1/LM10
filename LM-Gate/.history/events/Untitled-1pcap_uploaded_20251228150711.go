package events

type PCAPUploaded struct {
	FileID   string
	FilePath string
}

type Event interface {
	Name() string
}

func (e PCAPUploaded) Name() string {
	return "pcap.uploaded"
}

type Handler func(Event) error

type Dispatcher struct {
	handlers map[string][]Handler
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string][]Handler),
	}
}

func (d *Dispatcher) Register(eventName string, h Handler) {
	d.handlers[eventName] = append(d.handlers[eventName], h)
}

func (d *Dispatcher) Dispatch(e Event) error {
	if hs, ok := d.handlers[e.Name()]; ok {
		for _, h := range hs {
			if err := h(e); err != nil {
				return err
			}
		}
	}
	return nil
}
