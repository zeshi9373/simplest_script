package resident

type ResidentHandler interface {
	Handler()
}

var Entry = make(map[string]ResidentHandler)

func InitEntry() {
}
