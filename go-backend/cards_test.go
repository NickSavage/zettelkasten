package main

import (
	"go-backend/models"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func makeCardRequestSuccess(t *testing.T) *httptest.ResponseRecorder {

	token, _ := generateTestJWT(1)

	req, err := http.NewRequest("GET", "/api/cards/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(jwtMiddleware(s.getCard))
	handler.ServeHTTP(rr, req)

	return rr
}

func TestGetCardSuccess(t *testing.T) {
	setup()
	defer teardown()

	var logCount int
	_ = s.db.QueryRow("SELECT count(*) FROM card_views").Scan(&logCount)
	if logCount != 0 {
		t.Errorf("wrong log count, got %v want %v", logCount, 0)
	}
	rr := makeCardRequestSuccess(t)

	if status := rr.Code; status != http.StatusOK {
		log.Printf("err %v", rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	_ = s.db.QueryRow("SELECT count(*) FROM card_views").Scan(&logCount)
	if logCount != 1 {
		t.Errorf("wrong log count, got %v want %v", logCount, 1)
	}
	var card models.Card
	parseJsonResponse(t, rr.Body.Bytes(), &card)
	if card.ID != 1 {
		t.Errorf("handler returned wrong card, got %v want %v", card.ID, 1)
	}
	if card.UserID != 1 {
		t.Errorf("handler returned card for wrong user, got %v want %v", card.UserID, 1)
	}

}

func TestGetCardWrongUser(t *testing.T) {
	setup()
	defer teardown()

	token, _ := generateTestJWT(2)

	req, err := http.NewRequest("GET", "/api/cards/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(jwtMiddleware(s.getCard))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
	if rr.Body.String() != "unable to access card\n" {
		t.Errorf("handler returned wrong body, got %v want %v", rr.Body.String(), "unable to access card\n")
	}
	var logCount int
	_ = s.db.QueryRow("SELECT count(*) FROM card_views").Scan(&logCount)
	if logCount != 0 {
		t.Errorf("wrong log count, got %v want %v", logCount, 0)
	}
}

func TestGetParentCardId(t *testing.T) {
	cardID := "SP170/A.1/A.1/A.1/A.1"
	expected := "SP170/A.1/A.1/A.1/A"
	result := getParentIdAlternating(cardID)
	if result != expected {
		t.Errorf("function returned wrong result, got %v want %v", result, expected)
	}

	cardID = "1"
	expected = "1"
	result = getParentIdAlternating(cardID)
	if result != expected {
		t.Errorf("function returned wrong result, got %v want %v", result, expected)
	}

}

func TestGetCardSuccessParent(t *testing.T) {
	setup()
	defer teardown()

	rr := makeCardRequestSuccess(t)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var card models.Card
	parseJsonResponse(t, rr.Body.Bytes(), &card)
	if card.Parent.CardID != card.CardID {
		t.Errorf("wrong card parent returned. got %v want %v", card.Parent.CardID, card.CardID)
	}

}

func TestExtractBacklinks(t *testing.T) {
	text := "This is a sample text with [link1] and [another link]."
	expected := []string{"link1", "another link"}
	result := extractBacklinks(text)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestGetCardSuccessDirectLinks(t *testing.T) {
	setup()
	defer teardown()

	rr := makeCardRequestSuccess(t)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var card models.Card
	parseJsonResponse(t, rr.Body.Bytes(), &card)
	if len(card.DirectLinks) == 0 {
		t.Errorf("direct links was empty. got %v want %v", len(card.DirectLinks), 1)
	}

	expected := extractBacklinks(card.Body)

	if len(card.DirectLinks) > 0 && card.DirectLinks[0].CardID != expected[0] {
		t.Errorf("linked to wrong card, got %v want %v", card.DirectLinks[0].CardID, expected[0])

	}

}
