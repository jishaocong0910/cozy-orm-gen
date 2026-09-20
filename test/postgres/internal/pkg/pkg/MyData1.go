package pkg

type MyData1 []string

func (m MyData1) ToValue() string {
	return ""
}

func (m *MyData1) ToField(val string) {
}
