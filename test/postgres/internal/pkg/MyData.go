package pkg

type MyData struct {
}

func (m MyData) ToArg() string {
	return ""
}

func (m *MyData) ToField(val string) {
}
