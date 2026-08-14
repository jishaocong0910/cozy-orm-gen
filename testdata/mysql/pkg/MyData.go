package pkg

type MyData struct {
}

func (m MyData) ToValue() string {
	return ""
}

func (m *MyData) ToField(val string) {
}
