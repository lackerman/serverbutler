package viewmodels

type Home struct {
	Title   string
	Heading string
}

func GetHome() Home {
	return Home{
		Title:   "Server Butler",
		Heading: "Server Butler",
	}
}
