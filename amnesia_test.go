package amnesia

import "testing"

func TestAmnesia(t *testing.T) {
	tmpDir := t.TempDir()
	am := New(tmpDir)

	testInsertAndFindOne(t, am)
	testUpdate(t, am)
	testFind(t)
	testDelete(t, am)
}

func testInsertAndFindOne(t *testing.T, am *Amnesia) {
	col := am.Collection("test_collection")

	doc1 := Document{"name": "John Doe", "age": 30}
	doc2 := Document{"name": "Jane Smith", "age": 25}

	id1 := col.Insert(doc1)
	id2 := col.Insert(doc2)

	foundDoc1 := col.FindOne(Filter{"id": id1})
	foundDoc2 := col.FindOne(Filter{"id": id2})

	if foundDoc1 == nil || foundDoc2 == nil {
		t.Error("Failed to find inserted documents")
	}

	if foundDoc1["name"] != "John Doe" || foundDoc2["name"] != "Jane Smith" {
		t.Error("Found documents do not match the inserted ones")
	}
}

func testUpdate(t *testing.T, am *Amnesia) {
	col := am.Collection("test_collection")

	doc := Document{"name": "John Doe", "age": 30}
	id := col.Insert(doc)

	update := Document{"name": "Updated Name", "age": 35}
	col.Update(Filter{"id": id}, update)

	updatedDoc := col.FindOne(Filter{"id": id})
	if updatedDoc["name"].(string) != "Updated Name" || updatedDoc["age"].(float64) != 35 {
		t.Error("Failed to update the document")
	}
}

func testFind(t *testing.T) {
	tmpDir := t.TempDir()
	am := New(tmpDir)
	col := am.Collection("test_collection")

	col.Insert(Document{"name": "John Doe", "age": 30})
	col.Insert(Document{"name": "Jane Smith", "age": 25})

	filteredDocs := col.Find(Filter{"name": "Jane Smith"})
	if len(filteredDocs) != 1 {
		t.Error("Failed to find documents with the filter.")
	}
}

func testDelete(t *testing.T, am *Amnesia) {
	col := am.Collection("test_collection")

	doc := Document{"name": "John Doe", "age": 30}
	id := col.Insert(doc)

	col.Delete(Filter{"id": id})

	deletedDoc := col.FindOne(Filter{"id": id})
	if deletedDoc != nil {
		t.Error("Failed to delete documents")
	}
}
