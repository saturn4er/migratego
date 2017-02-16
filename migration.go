package migrates

type Migration struct {
	Number     int
	Name       string
	UpScript   string
	DownScript string
}

type byNumber []Migration

func (s byNumber) Len() int {
	return len(s)
}

func (s byNumber) Less(i, j int) bool {
	return s[i].Number < s[j].Number
}

func (s byNumber) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
