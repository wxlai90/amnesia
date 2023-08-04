package amnesia

import (
	"crypto/rand"
	"fmt"
)

type Amnesia struct {
	baseDir   string
	persistor Persistor
}

type Collection struct {
	name      string
	persistor Persistor
}

type Document map[string]interface{}

type Filter map[string]interface{}

func generateObjectID() string {
	buf := make([]byte, 16)
	rand.Read(buf)

	buf[6] = (buf[6] | 0x40) & 0x4F
	buf[8] = (buf[8] | 0x80) & 0xBF

	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:])
}

func New() *Amnesia {
	return &Amnesia{
		persistor: &memoryPersistor{
			collections: map[string]map[string]any{},
		},
	}
}

func NewWithFilePersistor(baseDir string) *Amnesia {
	return &Amnesia{
		persistor: &filePersistor{
			baseDir: baseDir,
		},
	}
}

func NewWithCustomPersistor(persistor Persistor) *Amnesia {
	return &Amnesia{
		persistor: persistor,
	}
}

func (am *Amnesia) Collection(name string) *Collection {
	return &Collection{
		persistor: am.persistor,
		name:      name,
	}
}

func (co *Collection) Find(filter Filter) []Document {
	all := co.persistor.Read(co.name)
	docs := []Document{}

	for _, doc := range all {
		match := true
		for k := range filter {
			if doc[k] != filter[k] {
				match = false
				break
			}
		}

		if match {
			docs = append(docs, doc)
		}
	}

	return docs
}

func (co *Collection) FindAll() []Document {
	return co.persistor.Read(co.name)
}

func (co *Collection) FindOne(filter Filter) Document {
	all := co.persistor.Read(co.name)
	for _, doc := range all {
		match := true
		for k := range filter {
			if doc[k] != filter[k] {
				match = false
				break
			}
		}

		if match {
			return doc
		}
	}

	return nil
}

func (co *Collection) Insert(x any) (string, error) {
	return co.persistor.Write(co.name, x)
}

func (co *Collection) Update(update Document) error {
	return co.persistor.Update(co.name, update)
}

func (co *Collection) Delete(filter Filter) error {
	return co.persistor.Delete(co.name, filter)
}
