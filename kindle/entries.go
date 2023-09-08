package kindle

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const SQL_COLLECTION_ENTRIES string = "SELECT p_uuid, p_titles_0_nominal FROM Entries WHERE p_type = 'Collection'"
const SQL_BOOK_ENTRIES string = "SELECT p_uuid, p_cdeKey from Entries WHERE p_type = 'Entry:Item' AND p_cdeType = 'EBOK' AND p_location LIKE '/mnt/us/documents/%'"

type Entries struct {
	Collection Collection
	Books      Books
	IsCcAware  bool
}

func get_Collection(db *sqlx.DB) (Collection, error) {
	var result Collection
	rows, err := db.Queryx(SQL_COLLECTION_ENTRIES)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry CollEntry
		rows.StructScan(&entry)
		if err != nil {
			return nil, err
		}
		result = append(result, &entry)
	}
	return result, nil
}

func get_Books(db *sqlx.DB) (Books, error) {
	var result Books
	rows, err := db.Queryx(SQL_BOOK_ENTRIES)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry Book
		rows.StructScan(&entry)
		if err != nil {
			return nil, err
		}
		result = append(result, &entry)
	}

	return result, nil
}

func NewEntries(db_file string) (*Entries, error) {
	db, err := sqlx.Connect("sqlite3", db_file)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	collection, err := get_Collection(db)
	if err != nil {
		return nil, err
	}

	books, err := get_Books(db)
	if err != nil {
		return nil, err
	}

	cursor, err := db.Query("SELECT i_collection_uuid, coalesce(i_member_uuid,'') as i_member_uuid FROM Collections")
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	for cursor.Next() {
		var uuid, member string
		err := cursor.Scan(&uuid, &member)
		if err != nil {
			return nil, err
		}
		if idx := collection.IdxByUUID(uuid); idx != -1 {
			collection := collection[idx]
			if idx_book := books.IdxByUUID(member); idx_book != -1 {
				collection.Books = append(collection.Books, books[idx_book])
			}
		}
	}

	var cc_aware bool
	db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('Entries') WHERE name='p_collectionCount'").Scan(&cc_aware)

	return &Entries{
		Collection: collection,
		Books:      books,
		IsCcAware:  cc_aware,
	}, nil
}
