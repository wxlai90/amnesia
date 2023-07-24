package amnesia

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Amnesia struct {
	baseDir string
}

type Collection struct {
	collection string
	baseDir    string
}

type Document map[string]interface{}

type Filter map[string]interface{}

func New(baseDir string) *Amnesia {
	return &Amnesia{
		baseDir: baseDir,
	}
}

func (am *Amnesia) Collection(collection string) *Collection {
	return &Collection{
		collection: collection,
		baseDir:    am.baseDir,
	}
}

func (co *Collection) Find(filter Filter) []Document {
	docs, err := readAllDocumentsInCollection(co.collection, co.baseDir)
	if err != nil {
		panic(err)
	}

	res := []Document{}
	for _, doc := range docs {
		if matchesFilter(doc, filter) {
			res = append(res, doc)
		}
	}

	return res
}

func (co *Collection) FindAll() []Document {
	docs, err := readAllDocumentsInCollection(co.collection, co.baseDir)
	if err != nil {
		panic(err)
	}

	return docs
}

func (co *Collection) FindOne(filter Filter) Document {
	docs, err := readAllDocumentsInCollection(co.collection, co.baseDir)
	if err != nil {
		panic(err)
	}

	for _, doc := range docs {
		if matchesFilter(doc, filter) {
			return doc
		}
	}

	return nil
}

func (co *Collection) Insert(doc Document) string {
	idValue, idExists := doc["id"]
	if !idExists {
		idValue = generateObjectID()
		doc["id"] = idValue
	}

	err := ensureCollectionDirExists(co.baseDir, co.collection)
	if err != nil {
		panic(err)
	}

	writeDocumentToDisk(co.collection, idValue.(string), doc, co.baseDir)

	return idValue.(string)
}

func (co *Collection) Update(filter Filter, update Document) {
	docs, err := readAllDocumentsInCollection(co.collection, co.baseDir)
	if err != nil {
		panic(err)
	}

	for _, doc := range docs {
		if matchesFilter(doc, filter) {
			updatedDoc := applyUpdate(doc, update)
			idValue, idExists := updatedDoc["id"]
			if idExists {
				err := ensureCollectionDirExists(co.baseDir, co.collection)
				if err != nil {
					panic(err)
				}

				writeDocumentToDisk(co.collection, idValue.(string), updatedDoc, co.baseDir)
			}
		}
	}
}

func (co *Collection) Delete(filter Filter) {
	docs, err := readAllDocumentsInCollection(co.collection, co.baseDir)
	if err != nil {
		panic(err)
	}

	for _, doc := range docs {
		if matchesFilter(doc, filter) {
			idValue, idExists := doc["id"]
			if idExists {
				deleteDocumentFromDisk(co.collection, idValue.(string), co.baseDir)
			}
		}
	}
}

func matchesFilter(doc Document, filter Filter) bool {
	for k, v := range filter {
		if doc[k] != v {
			return false
		}
	}

	return true
}

func applyUpdate(doc Document, update Document) Document {
	for k, v := range update {
		doc[k] = v
	}

	return doc
}

func writeDocumentToDisk(collection, id string, doc Document, baseDir string) {
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return
	}

	filename := path.Join(baseDir, collection, id+".json")
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return
	}
}

func readAllDocumentsInCollection(collection, baseDir string) ([]Document, error) {
	files, err := ioutil.ReadDir(path.Join(baseDir, collection))
	if err != nil {
		return nil, err
	}

	docs := []Document{}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			b, err := ioutil.ReadFile(path.Join(baseDir, collection, file.Name()))
			if err != nil {
				return nil, err
			}
			var doc Document
			err = json.Unmarshal(b, &doc)
			if err != nil {
				return nil, err
			}
			docs = append(docs, doc)
		}
	}

	return docs, nil
}

func deleteDocumentFromDisk(collection, id, baseDir string) {
	filename := path.Join(baseDir, collection, id+".json")
	_ = os.Remove(filename)
}

func generateObjectID() string {
	buf := make([]byte, 16)
	rand.Read(buf)

	buf[6] = (buf[6] | 0x40) & 0x4F
	buf[8] = (buf[8] | 0x80) & 0xBF

	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:])
}

func ensureCollectionDirExists(baseDir, collection string) error {
	_, err := os.Stat(path.Join(baseDir, collection))
	if os.IsNotExist(err) {
		err = os.MkdirAll(path.Join(baseDir, collection), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
