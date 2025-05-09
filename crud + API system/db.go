package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/uuid"
)

type Document struct {
	Id      uuid.UUID
	Author  string
	Storage string
}

type DbDocuments struct {
	Filename  string
	Documents map[uuid.UUID]Document
}

func InitDB(DBfilename string) DbDocuments {
	Idocuments, err := parseDBfile(DBfilename)
	if err != nil {
		log.Fatalf("%e", err)
	}

	return DbDocuments{Filename: DBfilename, Documents: Idocuments}
}

func parseDBfile(filename string) (map[uuid.UUID]Document, error) {
	documents := map[uuid.UUID]Document{}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.FieldsPerRecord = 3

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		parsedUUID, err := uuid.Parse(record[0])
		if err != nil {
			log.Fatalf("Error with converting string to UUID: %v", err)
		}
		documents[parsedUUID] = Document{Id: parsedUUID, Author: record[1], Storage: record[2]}
	}

	return documents, nil
}

func (db *DbDocuments) SearchDocumentById(reqUUID string) string {
	parseUUID, err := uuid.Parse(reqUUID)
	if err != nil {
		log.Fatalf("Error with converting string to UUID: %v", err)
	}
	document, check := db.Documents[parseUUID]
	if check {
		return fmt.Sprintf("\nAuthor of document with id: %v is: %v", document.Id, document.Author)
	} else {
		return fmt.Sprintf("\nDocument with id: %v isnt exist", reqUUID)
	}
}

func (db *DbDocuments) CreateDocument(author string, storage string) string {
	file, err := os.OpenFile(db.Filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("\nОшибка открытия файла:", err)
		return "\nError"
	}
	defer file.Close()

	newId := uuid.New()
	_, err = file.Write([]byte(fmt.Sprintf("\n%v,%v,%v", newId, author, storage)))
	if err != nil {
		log.Fatal(err)
	}
	db.Documents[newId] = Document{Id: newId, Author: author, Storage: storage}
	return fmt.Sprintf("\nDocument successfully create with id: %v", newId)
}

func (db *DbDocuments) UpdateDocument(reqUUID string, newAuthor string, newStorage string) string {
	UUID, err := uuid.Parse(reqUUID)
	if err != nil {
		log.Fatalf("Error with converting string to UUID: %v", err)
	}

	document, check := db.Documents[UUID]
	if check {
		document.Author = newAuthor
		document.Storage = newStorage
		db.Documents[UUID] = document
	} else {
		return fmt.Sprintf("\nDocument with id: %v isnt exist", reqUUID)
	}
	db.SaveToFile()

	return fmt.Sprintf("\nDocument with id: %v successfully updated", reqUUID)
}

func (db *DbDocuments) DeleteDocumentById(reqUUID string) string {
	UUID, err := uuid.Parse(reqUUID)
	if err != nil {
		log.Fatalf("Error with converting string to UUID: %v", err)
	}

	_, check := db.Documents[UUID]
	if check {
		delete(db.Documents, UUID)
	} else {
		return fmt.Sprintf("\nDocument with id: %v isnt exist", reqUUID)
	}

	db.SaveToFile()

	return fmt.Sprintf("\nDocument with id: %v successfully deleted", reqUUID)
}

func (db *DbDocuments) SaveToFile() error {
	tmpFile, err := os.CreateTemp("", "tempdb-*.csv")
	fmt.Print(tmpFile.Name(), "\n")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	for k, _ := range db.Documents {
		tmpFile.WriteString(fmt.Sprintf("%v,%v,%v\n", db.Documents[k].Id.String(), db.Documents[k].Author, db.Documents[k].Storage))
	}

	err = tmpFile.Close()
	if err != nil {
		return err
	}

	err = os.Rename(tmpFile.Name(), db.Filename)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	db := InitDB("C:/learning_go/api/apiTask/documents.csv")
	fmt.Print(db.UpdateDocument("9a1e8c3a-5b9e-4e05-d2f7-1d47d2f7b8c3", "Реально Крутой Чел", "*реально крутые цитаты*"))
}
