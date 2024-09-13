package main

import (
	"go-backend/models"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
)

// get existing tag
func TestGetTag(t *testing.T) {
	setup()
	defer teardown()

	userID := 1
	tagName := "test"

	tag, err := s.GetTag(userID, tagName)
	if err != nil {
		t.Errorf("handler returned error, %v", err.Error())
	}
	if tag.Name != tagName {
		t.Errorf("handler returned wrong tag, got %v want %v", tag.Name, tagName)
	}

}

func TestGetTagNotFound(t *testing.T) {
	setup()
	defer teardown()

	userID := 1
	tagName := "no tag by this name"

	tag, err := s.GetTag(userID, tagName)
	if err == nil {
		t.Errorf("handler did not return error when expecting error")
	}
	if tag.Name != "" {
		t.Errorf("handler returned a tag when not expecting one, got %v", tag.Name)
	}

}

func TestGetTagsRoute(t *testing.T) {
	setup()
	defer teardown()

	token, _ := generateTestJWT(1)

	req, err := http.NewRequest("GET", "/api/tags/cards/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(jwtMiddleware(s.GetTagsRoute))
	handler.ServeHTTP(rr, req)

	var tags []models.Tag
	parseJsonResponse(t, rr.Body.Bytes(), &tags)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	if len(tags) != 3 {
		t.Errorf("handler returned wrong number of tags, got %v want %v", len(tags), 3)
	}
}

// create new tag

func TestCreateTag(t *testing.T) {
	setup()
	defer teardown()

	userID := 1
	tagData := models.EditTagParams{
		Name:  "hello-world",
		Color: "black",
	}

	tag, err := s.CreateTag(userID, tagData)
	if err != nil {
		t.Errorf("handler returned error, %v", err.Error())
	}
	if tag.Name != tagData.Name {
		t.Errorf("handler returned wrong tag, got %v want %v", tag.Name, tagData.Name)
	}
	tag, err = s.GetTag(userID, tagData.Name)
	if err != nil {
		t.Errorf("handler returned error, %v", err.Error())
	}
	if tag.Name != tagData.Name {
		t.Errorf("handler returned a tag when not expecting one, got %v", tag.Name)
	}
}

// update tag (set new colour)

func TestEditTag(t *testing.T) {
	setup()
	defer teardown()

	userID := 1

	tagName := "test"
	tag, err := s.GetTag(userID, tagName)

	if err != nil {
		t.Errorf("handler returned error, %v", err.Error())
	}
	if tag.Name != tagName {
		t.Errorf("handler returned wrong tag, got %v want %v", tag.Name, tagName)
	}
	tagData := models.EditTagParams{
		Name:  "hello-world",
		Color: "red",
	}

	tag, err = s.EditTag(userID, tagName, tagData)
	if err != nil {
		t.Errorf("handler returned error, %v", err.Error())
	}
	if tag.Name != tagData.Name {
		t.Errorf("handler returned wrong tag, got %v want %v", tag.Name, tagData.Name)
	}
	tag, err = s.GetTag(userID, tagData.Name)
	if err != nil {
		t.Errorf("handler returned err %v", err)
	}
	if tag.Name != tagData.Name {
		t.Errorf("handler returned wrong tag, got %v want %v", tag.Name, tagData.Name)
	}
}

func TestCreateTagOverExisting(t *testing.T) {
	setup()
	defer teardown()

	userID := 1

	tagName := "test"
	tag, err := s.GetTag(userID, tagName)
	oldID := tag.ID

	if err != nil {
		t.Errorf("handler returned error, %v", err.Error())
	}
	if tag.Name != tagName {
		t.Errorf("handler returned wrong tag, got %v want %v", tag.Name, tagName)
	}
	tagData := models.EditTagParams{
		Name:  tagName,
		Color: "red",
	}
	tag, err = s.CreateTag(userID, tagData)
	if err != nil {
		t.Errorf("handler returned error, %v", err.Error())
	}
	if tag.ID != oldID {
		t.Errorf("handler returned wrong tag, got id %v want id %v", tag.ID, oldID)
	}
	tag, err = s.GetTag(userID, tagData.Name)
	if err != nil {
		t.Errorf("handler returned error, %v", err.Error())
	}
	if tag.ID != oldID {
		t.Errorf("handler returned wrong tag, got id %v want id %v", tag.ID, oldID)
	}
	if tag.Name != tagData.Name {
		t.Errorf("handler returned wrong tag, got %v want %v", tag.Name, tagData.Name)
	}

}

// add tags for a card

func TestAddTagToCard(t *testing.T) {

	setup()
	defer teardown()

	userID := 1
	cardPK := 1
	tagName := "test"

	var count int
	_ = s.db.QueryRow("SELECT count(*) FROM card_tags").Scan(&count)

	err := s.AddTagToCard(userID, tagName, cardPK)
	if err != nil {
		t.Errorf("handler returned error, %v", err.Error())
	}

	var newCount int
	_ = s.db.QueryRow("SELECT count(*) FROM card_tags").Scan(&newCount)
	if newCount != (count + 1) {
		t.Errorf("handler returned wrong number of card_tags, got %v want %v", newCount, count+1)
	}

}

func TestAddTagsFromCardQuery(t *testing.T) {
	setup()
	defer teardown()

	userID := 2

	cardPK := 23
	tagName := "to-read"

	var count int
	_ = s.db.QueryRow("SELECT count(*) FROM card_tags").Scan(&count)

	err := s.AddTagsFromCard(userID, cardPK)
	if err != nil {
		t.Errorf("handler returned error, %v", err.Error())
	}

	var newCount int
	_ = s.db.QueryRow("SELECT count(*) FROM card_tags").Scan(&newCount)
	if newCount != (count + 1) {
		t.Errorf("handler returned wrong number of card_tags, got %v want %v", newCount, count+1)
	}

	tags, err := s.QueryTagsForCard(userID, cardPK)
	if err != nil {
		t.Errorf("handler returned error, %v", err.Error())
	}
	if len(tags) != 1 {
		t.Errorf("handler returned wrong number of tags, got %v want %v", len(tags), 1)
	}

	if tags[0].Name != tagName {
		t.Errorf("wrong tag attached to card, got %v want %v", tags[0].Name, tagName)
	}

}

func TestAddTagsFromTaskQuery(t *testing.T) {
	setup()

	userID := 1
	taskPK := 3
	tagName := "to-read"

	var count int
	_ = s.db.QueryRow("SELECT count(*) FROM task_tags").Scan(&count)

	err := s.AddTagsFromTask(userID, taskPK)
	if err != nil {
		t.Errorf("handler returned error, %v", err.Error())
	}

	var newCount int
	_ = s.db.QueryRow("SELECT count(*) FROM task_tags").Scan(&newCount)
	if newCount != (count + 1) {
		t.Errorf("handler returned wrong number of task_tags, got %v want %v", newCount, count+1)
	}

	tags, err := s.QueryTagsForTask(userID, taskPK)
	if err != nil {
		t.Errorf("handler returned error, %v", err.Error())
	}
	if len(tags) != 1 {
		t.Errorf("handler returned wrong number of tags, got %v want %v", len(tags), 1)
	}

	if tags[0].Name != tagName {
		t.Errorf("wrong tag attached to card, got %v want %v", tags[0].Name, tagName)
	}
}

func TestParseTagsFromCardBody(t *testing.T) {
	setup()
	defer teardown()

	body := "hello world \n\n#to-read #hello#world"

	tags, err := s.ParseTagsFromCardBody(body)
	if err != nil {
		t.Errorf("handler returned error, %v", err.Error())
	}
	if len(tags) != 2 {
		t.Errorf("handler returned wrong number of tags, got %v want %v", len(tags), 2)
	}
	if tags[0] != "to-read" {
		t.Errorf("wrong tag returned, got %v want %v", tags[0], "to-read")
	}
	if tags[1] != "hello" {
		t.Errorf("wrong tag returned, got %v want %v", tags[1], "hello")
	}

	body = "hello world \n\nto-read#hello #world"
	tags, err = s.ParseTagsFromCardBody(body)
	if err != nil {
		t.Errorf("handler returned error, %v", err.Error())
	}
	if len(tags) != 1 {
		log.Printf("tags %v", tags)
		t.Errorf("handler returned wrong number of tags, got %v want %v", len(tags), 1)
	}
	if tags[0] != "world" {
		t.Errorf("wrong tag returned, got %v want %v", tags[0], "world")
	}

}

func TestDeleteTag(t *testing.T) {
	setup()
	defer teardown()

	id := 1
	token, _ := generateTestJWT(1)

	req, err := http.NewRequest("DELETE", "/api/tags/id/"+strconv.Itoa(id), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.SetPathValue("id", strconv.Itoa(id))

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/tags/id/{id}", jwtMiddleware(s.DeleteTagRoute))
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	_, err = s.GetTag(1, "test")
	if err == nil {
		t.Error("handler returned tag after it should have been deleted")
	}
}
