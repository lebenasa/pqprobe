// Package musics is model code to table musics in database.
package musics




import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type (
	// Table is interface to table musics.
	Table interface {
		// Select retrieve musics' record with specific primary keys value.
		Select(ID int64) (Record, error)

		// Insert new records to table musics.
		Insert(tx *sqlx.Tx, data Data) (Record, error)
	}

	// table implements interface Table.
	table struct {
		dbSlave *sqlx.DB
	}

	// Record is interface to table musics' record.
	Record interface {
		// Fields returns records' data.
		Data() Data
		// Updates records' data.
		Update(tx *sqlx.Tx, data Data) (Record, error)
	}

	// record implements interface Record.
	record struct {
		data Data
	}

	// Data contains table musics' fields.
	Data struct {
	
		ID int64 `db:"id"`
	
		Artist string `db:"artist"`
	
		Title string `db:"title"`
	
		Album string `db:"album"`
	
		ReleaseDate time.Time `db:"release_date"`
	
		LastPlayed time.Time `db:"last_played"`
	
		Rating float64 `db:"rating"`
	
		Description string `db:"description"`
	
		valid       bool
	}
)

var (
	// ErrInvalidData returned when trying to modify a record with handcrafted data.
	ErrInvalidData = errors.New("invalid data")
)

// Bind slave database connection to a new Table.
func Bind(slave *sqlx.DB) (t Table, err error) {
	if slave == nil {
		return nil, errors.New("musics: nil slave connection")
	}
	return table{dbSlave: slave}, nil
}

// Select returns Record with given primary key.
func (t table) Select(ID int64) (r Record, err error) {
	rec := record{}
	err = t.dbSlave.Select(&rec, `
	SELECT
	
		id,
	
		artist,
	
		title,
	
		album,
	
		release_date,
	
		last_played,
	
		rating,
	
		description
	
	FROM
		musics
	WHERE
		
			id = $1
		
	`, 
	
		ID)
	if err != nil {
		return nil, errors.Wrap(err, "musics: select error")
	}

	rec.data.valid = true
	return rec, nil
}

// Insert inserts new Record into the Table and returns newly created Record if insertion succeed.
func (t table) Insert(tx *sqlx.Tx, data Data) (r Record, err error) {
	if tx == nil {
		return nil, errors.New("musics: nil tx")
	}
	err = tx.QueryRowx(`
		INSERT INTO
			musics (
				
					artist,
				
					title,
				
					album,
				
					release_date,
				
					last_played,
				
					rating,
				
					description
				
			)
		VALUES
			(
				
					$1,
				
					$2,
				
					$3,
				
					$4,
				
					$5,
				
					$6,
				
					$7
				
			)
		RETURNING
			
				id,
			
				artist,
			
				title,
			
				album,
			
				release_date,
			
				last_played,
			
				rating,
			
				description
			
		`, 
		
			data.Artist,
			data.Title,
			data.Album,
			data.ReleaseDate,
			data.LastPlayed,
			data.Rating,
			data.Description).StructScan(&data)
	if err != nil {
		return nil, errors.Wrap(err, "musics: insert failed")
	}

	return record{data: data}, nil
}

func (r record) Data() Data {
	return r.data
}

func (r record) Update(tx *sqlx.Tx, data Data) (newRecord Record, err error) {
	if tx == nil {
		return nil, errors.New("musics: nil tx")
	}
	if !data.valid {
		return nil, ErrInvalidData
	}

	_, err = tx.Exec(`
		UPDATE
			musics
		SET
			
				artist = $1,
			
				title = $2,
			
				album = $3,
			
				release_date = $4,
			
				last_played = $5,
			
				rating = $6,
			
				description = $7
			
		WHERE
		
			id = $9
		
		`, 
		
			data.Artist,
		
			data.Title,
		
			data.Album,
		
			data.ReleaseDate,
		
			data.LastPlayed,
		
			data.Rating,
		
			data.Description,
		
		
			data.ID)
	if err != nil {
		return nil, errors.Wrap(err, "musics: update failed")
	}

	return record{data: data}, nil
}
