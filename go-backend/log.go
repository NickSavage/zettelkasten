package main

import (
	"go-backend/models"
	"log"
)

func (s *Server) logCardView(cardPK int, userID int) {
	_, err := s.db.Exec(`
   INSERT INTO 
   card_views 
   (card_pk, user_id, created_at) 
   VALUES ($1, $2, CURRENT_TIMESTAMP);`, cardPK, userID)

	if err != nil {
		// Log the error
		log.Printf("Error logging card view for cardPK %d and userID %d: %v", cardPK, userID, err)
	}
}

func (s *Server) LogLastLogin(user models.User) {
	_, err := s.db.Exec(`UPDATE users SET last_login = NOW() WHERE id = $1`, user.ID)
	if err != nil {
		// Log the error
		log.Printf("Error logging card view for userID %v: %v", user.ID, err)
	}

}
