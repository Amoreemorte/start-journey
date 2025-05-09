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
	Id      uuid.UUID `json:"id"`
	Author  string    `json:"author"`
	Content string    `json:"content"`
}

type Dbdocuments struct {
	filename  string
	documents map[uuid.UUID]Document
}

func InitDB(DBfilename string) Dbdocuments {
	Idocuments, err := parseDBfile(DBfilename)
	if err != nil {
		log.Printf("%e", err)
	}

	return Dbdocuments{filename: DBfilename, documents: Idocuments}
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
		documents[parsedUUID] = Document{Id: parsedUUID, Author: record[1], Content: record[2]}
	}

	return documents, nil
}

func (db *Dbdocuments) SearchDocumentById(reqUUID string) (Document, error) {
	parseUUID, err := uuid.Parse(reqUUID)
	if err != nil {
		log.Printf("Error with converting string to UUID: %v", err)
	}
	document, check := db.documents[parseUUID]
	if check {
		return document, nil
	} else {
		return Document{}, errors.New("Document dont exist")
	}
}

func (db *Dbdocuments) CreateDocument(author string, content string) (Document, error) {
	file, err := os.OpenFile(db.filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return Document{}, fmt.Errorf("\nОшибка открытия файла:", err)
	}
	defer file.Close()

	newId := uuid.New()
	_, err = file.Write([]byte(fmt.Sprintf("\n%v,%v,%v", newId, author, content)))
	if err != nil {
		return Document{}, err
	}
	db.documents[newId] = Document{Id: newId, Author: author, Content: content}
	return db.documents[newId], nil
}

func (db *Dbdocuments) UpdateDocument(reqUUID string, newAuthor string, newContent string) (Document, error) {
	UUID, err := uuid.Parse(reqUUID)
	if err != nil {
		return Document{}, err
	}

	document, check := db.documents[UUID]
	if check {
		document.Author = newAuthor
		document.Content = newContent
		db.documents[UUID] = document
	} else {
		return Document{}, fmt.Errorf("\nDocument with id: %v isnt exist", reqUUID)
	}
	db.SaveToFile()

	return db.documents[UUID], err
}

func (db *Dbdocuments) DeleteDocumentById(reqUUID string) error {
	UUID, err := uuid.Parse(reqUUID)
	if err != nil {
		return fmt.Errorf("Error with converting string to UUID: %v", err)
	}

	_, check := db.documents[UUID]
	if check {
		delete(db.documents, UUID)
	} else {
		return fmt.Errorf("\nDocument with id: %v isnt exist", reqUUID)
	}

	db.SaveToFile()

	return nil
}

func (db *Dbdocuments) SaveToFile() error {
	tmpFile, err := os.CreateTemp("", "tempdb-*.csv")
	fmt.Print(tmpFile.Name(), "\n")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	for k, _ := range db.documents {
		tmpFile.WriteString(fmt.Sprintf("%v,%v,%v\n", db.documents[k].Id.String(), db.documents[k].Author, db.documents[k].Content))
	}

	err = tmpFile.Close()
	if err != nil {
		return err
	}

	err = os.Rename(tmpFile.Name(), db.filename)
	if err != nil {
		return err
	}

	return nil
}
