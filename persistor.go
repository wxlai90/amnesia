package amnesia

type Persistor interface {
	Write(collection string, x any) (string, error)
	Read(collection string) []Document
	Update(collection string, doc Document) error
	Delete(collection string, filter Filter) error
}
