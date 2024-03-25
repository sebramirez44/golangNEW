package models

import (
	"database/sql"
	"errors"
	"time"
)

//definimos Snuippet type, corresponden a los fields en la tabla de mysql
type Snippet struct {
	ID		int
	Title	string
	Content	string
	Created	time.Time
	Expires	time.Time
}

//snippetModel que wrappea sql.DB
type SnippetModel struct {
	DB *sql.DB
}

//insertamos snippets a la base de datos
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	//en Postgresql el stmt usa $N en vez de ? para placeholders
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?,?,UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	//nos regresa un sql.Result type, en postregsql Exec se hace de diferente forma.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	//LastInsertId() no funciona con el driver de Postgresql, buscar documentacion de ese para ver como do that shit
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	
	//id que tenemos es itn64 entonces lo convertimos a int normal chill
	return int(id), nil
}

//retrun snippet based on id
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := "SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?"

	row := m.DB.QueryRow(stmt, id)

	//inicializamos un pointer a un nuevo struct Snippet
	s := &Snippet{}

	// parametros son pointers a donde queremos insertar la informacion de la sql row, el numero de parametros
	// tiene que ser el mismo numero de valroes retornados.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// si no se encontro un registro, row.Scan regresa un sql.ErrorNoRows, con errors.ls checamos que sea ese error.
		// y retornamos un error custom ErrNoRecord
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	//Regresamos el snippet object si no hubo errores
	return s, nil
}

//return 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := "SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10"

	//Query regresa mas de una row
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	
	// antes de que la funcion retorne, cerramos sql.Rows resultset wtf is resultset
	// defer tiene que ir despues de checar error en la query porque si no tratariamos de cerrar un error y hay panic
	defer rows.Close()

	//empty slice to hold Snippet structs, slice es como un array
	snippets := []*Snippet{}

	//rows.Next() itera en las rows de un resultset que es lo que tenemos.
	//si iteramos por todos estos automaticamente se cierra el resultset, pero entonces no entiendo porq hicimos defer rows.Close()?
	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		//append el snippet de esta iteracion a la slice
		snippets = append(snippets, s)
	}

	//rows.Err() nos dice cualquier error que ocurrio durante la iteracion.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
