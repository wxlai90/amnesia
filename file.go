package amnesia

import (
	"encoding/json"
	"os"
	"path"
	"strings"
)

func (fp *filePersistor) Read(collection string) []Document {
	_, err := os.Stat(path.Join(fp.baseDir, collection))
	if err != nil {
		return []Document{}
	}

	files, err := os.ReadDir(path.Join(fp.baseDir, collection))
	if err != nil {
		return nil
	}

	docs := []Document{}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			b, err := os.ReadFile(path.Join(fp.baseDir, collection, file.Name()))
			if err != nil {
				return nil
			}
			var doc Document
			err = json.Unmarshal(b, &doc)
			if err != nil {
				return nil
			}
			docs = append(docs, doc)
		}
	}

	return docs
}

func (fp *filePersistor) Write(collection string, x any) (string, error) {
	err := ensureCollectionDirExists(fp.baseDir, collection)
	if err != nil {
		return "", err
	}

	switch x.(type) {
	case string:
		panic("cannot be primitive type")
	case bool:
		panic("cannot be primitive type")
	case uint:
		panic("cannot be primitive type")
	case int:
		panic("cannot be primitive type")
	}

	var doc Document
	data, err := json.Marshal(x)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(data, &doc)
	if err != nil {
		return "", err
	}

	idValue, idExists := doc["id"]
	if !idExists {
		idValue = generateObjectID()
		doc["id"] = idValue
	}

	idValueStr, ok := idValue.(string)
	if !ok {
		return "", ErrIdNotOfStringType
	}

	data, err = json.Marshal(doc)
	if err != nil {
		return "", err
	}

	filename := path.Join(fp.baseDir, collection, idValueStr+".json")
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return "", err
	}

	return idValueStr, nil
}

func (fp *filePersistor) Update(collection string, doc Document) error {
	id, ok := doc["id"].(string)
	if !ok {
		return ErrIdNotFound
	}

	oldDoc := findDocById(path.Join(fp.baseDir, collection, id+".json"))
	if oldDoc == nil {
		return ErrDocNotExist
	}

	_, err := fp.Write(collection, doc)
	return err
}

func (fp *filePersistor) Delete(collection string, filter Filter) error {
	docs := fp.Read(collection)

	for _, doc := range docs {
		if matchesFilter(doc, filter) {
			id := doc["id"].(string)
			filename := path.Join(fp.baseDir, collection, id+".json")
			err := os.Remove(filename)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type filePersistor struct {
	baseDir string
}

func matchesFilter(doc Document, filter Filter) bool {
	for k := range filter {
		if doc[k] != filter[k] {
			return false
		}
	}

	return true
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

func findDocById(fullpath string) Document {
	b, err := os.ReadFile(fullpath)
	if err != nil {
		return nil
	}
	var doc Document
	err = json.Unmarshal(b, &doc)
	if err != nil {
		return nil
	}

	return doc
}
