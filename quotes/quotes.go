package quotes

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/bruno-chavez/restedancestor/database"
	"github.com/satori/go.uuid"
)

// init is used to seed the rand.Intn function.
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// QuoteType describes a quote.
type QuoteType struct {
	Quote string    `json:"quote"`
	Uuid  uuid.UUID `json:"uuid"`
	Score int       `json:"score"`
}

// Quotes is used to parse the whole json database.
type Quotes struct {
	Data    []QuoteType `json:"data"`
	Indexes indexes     `json:"index"`
}

// Random returns a random quote from a Quotes type.
func (q Quotes) Random() QuoteType {
	qd := q.Data
	return qd[rand.Intn(len(qd))]
}

func (q Quotes) Len() int {
	return len(q.Data)
}

func (q Quotes) Swap(i int, j int) {
	qd := q.Data
	qd[i], qd[j] = qd[j], qd[i]
}

func (q Quotes) Less(i int, j int) bool {
	qd := q.Data
	return qd[i].Score > qd[j].Score
}

// Parser fetches from database.json and puts it on a struct.
func Parser(data database.Database) Quotes {
	q := Quotes{}
	err := json.Unmarshal(data.Read(), &q)
	if err != nil {
		log.Fatal(err)
	}

	return q
}

// LikeQuote increments the score of the quote
func (q Quotes) LikeQuote(db database.Database, uuid string) {
	offset, _ := q.OffsetQuoteFromUUID(uuid)
	q.Data[*offset].Score++

	if err := unparser(db, q); err != nil {
		log.Fatal(err)
	}
}

// DislikeQuote decrements the score of the quote
func (q Quotes) DislikeQuote(db database.Database, uuid string) {
	offset, _ := q.OffsetQuoteFromUUID(uuid)
	q.Data[*offset].Score--

	if err := unparser(db, q); err != nil {
		log.Fatal(err)
	}
}

// OffsetQuoteFromUUID find the uuid in the slice and returns its offset
func (q Quotes) OffsetQuoteFromUUID(uuid string) (*int, error) {
	for k, quote := range q.Data {
		if quote.Uuid.String() == uuid {
			return &k, nil
		}
	}

	return nil, errors.New("unknown")
}

// index is an inverted index : for each word, list all documents that contain it.
type index struct {
	Word  string      `json:"word"`
	Uuids []uuid.UUID `json:"uuids"`
}

type indexes []index

// Index indexes all the data
func (q Quotes) Index(db database.Database) {
	q.Indexes = make(indexes, 0)
	const limitSize = 3
	for _, quote := range q.Data {
		words := strings.FieldsFunc(quote.Quote, func(r rune) bool {
			switch r {
			case '\'', '!', ',', '.', ' ':
				return true
			}
			return false
		})
		for _, word := range words {
			if len(word) > limitSize {
				// First Find offset
				o := -1
				for k, index := range q.Indexes {
					if index.Word == word {
						o = k
						break
					}
				}

				// Doesn't exist, create index
				if o == -1 {
					idx := index{
						Word:  word,
						Uuids: []uuid.UUID{quote.Uuid},
					}
					q.Indexes = append(q.Indexes, idx)
					continue
				}

				// Exists, append u
				index := &(q.Indexes[o])
				index.Uuids = append(index.Uuids, quote.Uuid)
			}
		}
	}
	unparser(db, q)
}

// unparser writes a slice into database.
func unparser(db database.Database, quotes Quotes) error {
	writeJSON, err := json.MarshalIndent(quotes, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return db.Write(writeJSON)
}
