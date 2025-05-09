package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/uuid"
)

type Document struct {
	Id      uuid.UUID "json:id"
	Author  string    "json:author"
	Storage string    "json:storage"
}

type DbDocuments struct {
	Filename  string
	Documents map[uuid.UUID]Document
}

func InitDB(DBfilename string) DbDocuments {
	Idocuments, err := parseDBfile(DBfilename)
	if err != nil {
		log.Printf("%e", err)
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
			log.Printf("Error with converting string to UUID: %v", err)
		}
		documents[parsedUUID] = Document{Id: parsedUUID, Author: record[1], Storage: record[2]}
	}

	return documents, nil
}

func (db *DbDocuments) SearchDocumentById(reqUUID string) (Document, error) {
	parseUUID, err := uuid.Parse(reqUUID)
	if err != nil {
		log.Printf("Error with converting string to UUID: %v", err)
	}
	document, check := db.Documents[parseUUID]
	if check {
		return document, nil
	} else {
		return Document{}, errors.New("Document dont exist")
	}
}

func (db *DbDocuments) CreateDocument(author string, storage string) (Document, error) {
	file, err := os.OpenFile(db.Filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return Document{}, fmt.Errorf("\nОшибка открытия файла:", err)
	}
	defer file.Close()

	newId := uuid.New()
	_, err = file.Write([]byte(fmt.Sprintf("\n%v,%v,%v", newId, author, storage)))
	if err != nil {
		return Document{}, err
	}
	db.Documents[newId] = Document{Id: newId, Author: author, Storage: storage}
	return db.Documents[newId], nil
}

func (db *DbDocuments) UpdateDocument(reqUUID string, newAuthor string, newStorage string) (Document, error) {
	UUID, err := uuid.Parse(reqUUID)
	if err != nil {
		return Document{}, err
	}

	document, check := db.Documents[UUID]
	if check {
		document.Author = newAuthor
		document.Storage = newStorage
		db.Documents[UUID] = document
	} else {
		return Document{}, fmt.Errorf("\nDocument with id: %v isnt exist", reqUUID)
	}
	db.SaveToFile()

	return db.Documents[UUID], err
}

func (db *DbDocuments) DeleteDocumentById(reqUUID string) error {
	UUID, err := uuid.Parse(reqUUID)
	if err != nil {
		return fmt.Errorf("Error with converting string to UUID: %v", err)
	}

	_, check := db.Documents[UUID]
	if check {
		delete(db.Documents, UUID)
	} else {
		return fmt.Errorf("\nDocument with id: %v isnt exist", reqUUID)
	}

	db.SaveToFile()

	return nil
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
