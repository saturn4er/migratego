package migrates

type Migration struct {
	Version    int
	Name       string
	UpScript   string
	DownScript string
}

type byVersion []Migration

func (s byVersion) Len() int {
	return len(s)
}

func (s byVersion) Less(i, j int) bool {
	return s[i].Version < s[j].Version
}

func (s byVersion) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
