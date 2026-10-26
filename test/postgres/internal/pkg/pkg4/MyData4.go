package pkg_2

type MyData4 []string

func (m MyData4) ToArg() string {
	return ""
}

func (m *MyData4) ToField(val string) {
}
