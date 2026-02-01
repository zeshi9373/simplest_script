package resident

func Process() {
	InitEntry()

	if len(Entry) > 0 {
		for _, entry := range Entry {
			go entry.Handler()
		}
	}
}
