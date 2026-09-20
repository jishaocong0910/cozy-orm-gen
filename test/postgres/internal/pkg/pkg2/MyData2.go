package pkg2

type MyData2 map[string]string

func (m MyData2) ToValue() string {
	return ""
}

func (m *MyData2) ToField(val string) {
}
