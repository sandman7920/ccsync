package kindle

type CollEntry struct {
	UUID  string `db:"p_uuid"`
	Title string `db:"p_titles_0_nominal"`
	Books Books
}

func (c *CollEntry) Members() []string {
	members := make([]string, 0, len(c.Books))
	for _, m := range c.Books {
		members = append(members, m.UUID)
	}
	return members
}

type Collection []*CollEntry

func (c Collection) IdxByUUID(uuid string) int {
	for idx, e := range c {
		if e.UUID == uuid {
			return idx
		}
	}
	return -1
}

func (c Collection) IdxByTitle(title string) int {
	for idx, e := range c {
		if e.Title == title {
			return idx
		}
	}
	return -1
}
