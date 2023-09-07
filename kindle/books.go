package kindle

type Book struct {
	UUID   string `db:"p_uuid"`
	CDEKey string `db:"p_cdeKey"`
}

type Books []*Book

func (b Books) IdxByUUID(uuid string) int {
	for idx, e := range b {
		if e.UUID == uuid {
			return idx
		}
	}
	return -1
}

func (b Books) BookByUUID(uuid string) *Book {
	idx := b.IdxByUUID(uuid)
	if idx == -1 {
		return nil
	}

	return b[idx]
}

func (b Books) BookByCDEKey(cde_key string) *Book {
	for _, book := range b {
		if book.CDEKey == cde_key {
			return book
		}
	}
	return nil
}

func (b Books) BooksByCDEKeys(cde_keys []string) Books {
	var result Books
	for _, key := range cde_keys {
		book := b.BookByCDEKey(key)
		if book != nil {
			result = append(result, book)
		}
	}
	return result
}

func (b Books) CDEKeys() []string {
	result := make([]string, 0, len(b))
	for _, b := range b {
		result = append(result, b.CDEKey)
	}
	return result
}
