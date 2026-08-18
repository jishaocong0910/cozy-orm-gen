package pkg_2

type MyData3 string

func (m MyData3) ToValue() string {
	return ""
}

func (m *MyData3) ToField(val string) {
}
