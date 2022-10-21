package mirror

type Data struct {
	Hostname string   `json:"hostname"`
	Mirrors  []Mirror `json:"mirrors"`
}

type Mirror struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	Path      string   `json:"path"`
	Frequency int      `json:"frequency"`
	Protocols []string `json:"protocols"`
	Status    string   `json:"status"`
}
