package amnesia

import (
	"encoding/json"
)

type memoryPersistor struct {
	collections map[string]map[string]any
}

func (mp *memoryPersistor) Read(collection string) []Document {
	docs := make([]Document, len(mp.collections))
	i := 0

	for _, doc := range mp.collections {
		docs[i] = doc
		i++
	}

	return docs
}

func (mp *memoryPersistor) Write(collection string, x any) (string, error) {
	var doc Document
	data, err := json.Marshal(x)
	if err != nil {
		return "", err
	}

	json.Unmarshal(data, &doc)

	idValue, idExists := doc["id"]
	if !idExists {
		idValue = generateObjectID()
		doc["id"] = idValue
	}

	idValueStr, ok := idValue.(string)
	if !ok {
		return "", ErrIdNotOfStringType
	}

	mp.collections[idValueStr] = doc

	return idValueStr, nil
}

func (mp *memoryPersistor) Update(collection string, doc Document) error {
	id, ok := doc["id"]
	if !ok {
		return ErrIdNotFound
	}

	idStr, ok := id.(string)
	if !ok {
		return ErrIdNotOfStringType
	}

	mp.collections[collection][idStr] = doc
	return nil
}

func (mp *memoryPersistor) Delete(collection string, filter Filter) error {
	id, ok := filter["id"].(string)
	if !ok {
		return ErrIdNotOfStringType
	}

	delete(mp.collections, id)
	return nil
}
