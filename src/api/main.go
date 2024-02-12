package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-memdb"
)

var db *memdb.MemDB

func main() {
	fmt.Println("Iniciando banco de dados")
	if err := initializeBD(); err != nil {
		log.Fatal("Falha ao inicializar banco", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/pets", Create).Methods("POST")
	router.HandleFunc("/pets/{id}", Get).Methods("GET")
	router.HandleFunc("/pets/{id}", Put).Methods("PUT")
	router.HandleFunc("/pets/{id}", Delete).Methods("DELETE")

	fmt.Println("Iniciando API Pet Store")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// CREATE - POST
func Create(w http.ResponseWriter, r *http.Request) {
	// Declarando variavel pet para realizar a desearilização do body enviado
	// na requisição e transformar na struct Pet
	var pet *Pet

	//Transformando (Desearilizando) o body para struct pet
	if err := json.NewDecoder(r.Body).Decode(&pet); err != nil {
		http.Error(w, "Falha ao desearilizar body", http.StatusBadRequest)
		return
	}

	// Cria uma transação de escrita
	txn := db.Txn(true)

	//Inserindo no banco de dados
	if err := txn.Insert("pet", pet); err != nil {
		http.Error(w, "Falha ao inserir no banco de dados", http.StatusInternalServerError)
		return
	}

	// Confirma transação no banco
	txn.Commit()

	//Retorna para requisição o código 201 para informar que foi criado com sucesso
	w.WriteHeader(http.StatusCreated)
}

// READ - GET
func Get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	// Cria uma transação de leitura
	txn := db.Txn(false)

	res, err := txn.First("pet", "id", id)
	if err != nil {
		http.Error(w, "Falha ao buscar no banco", http.StatusNotFound)
		return
	}

	if res == nil {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(res.(*Pet))
}

// UPDATE - PUT
func Put(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(r.Body)
}

// DELETE - DELETE
func Delete(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(r.Body)
}

// Structs
type Pet struct {
	ID     int    `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
}

func initializeBD() error {
	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"pet": {
				Name: "pet",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "ID"},
					},
				},
			},
		},
	}

	// Create a new data base
	dba, err := memdb.NewMemDB(schema)
	if err != nil {
		return err
	}

	db = dba
	return nil
}

// string = texto  ex.: "abc123"
// int = inteiro ex.: 123, -132
// float = decimais ex.: 1.23
// bool = verdadeiro ou false ex.: true
// interface{} = generico
// any = alias(apelido) interface{}
// byte[] = lista de bytes
// uint = inteiro positivo
