package server

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"vintz.fr/nonogram/partie"
)


const SESSION_COOKIE = "NonoSession"

func ImagesHandler( path, dir string ) http.Handler {
	return http.StripPrefix( path, http.FileServer(http.Dir(dir)) )
}


type Session struct {
	sid int64
	etat partie.Etat
	nbHits int
	partie partie.Partie
}


type SessionStore struct {
	sidGenerator rand.Source
	sessions map[int64]*Session
}


func NewSessionStore() *SessionStore {
	return &SessionStore{
		sidGenerator: rand.NewSource( time.Now().UnixNano() ),
		sessions:make( map[int64]*Session ) }
}


func (store *SessionStore) getSession( sid int64 ) *Session {

	s, exists := store.sessions[sid]

	if exists {
		log.Printf("Session %d existante", sid)
	}
	if !exists {
		log.Printf("Session %d inexistante", sid)
		sid = (store.sidGenerator.Int63()>>1) + 1   // non-null random integer
		s = &Session{sid: sid}
		store.sessions[sid] = s
		log.Printf("Session %d créée", sid)
	}
	if s.sid != sid {
		panic( "SessionStore : incohérence dans les id de session" )
	}
	return s
}


func MakeNonoHandler( ) http.HandlerFunc {

	sessionStore := NewSessionStore()

	// le sessionStore est capturé dans le handler
	handler := func (w http.ResponseWriter, req *http.Request) {

		// recupere le sid dans SESSION_COOKIE dans la requete
		var sid int64
		cookie, err := req.Cookie(SESSION_COOKIE)
		if err == nil {
			sid, err = strconv.ParseInt(cookie.Value, 10, 64)
			if err!=nil {
				log.Printf("Cookie de session '%v' invalide : %v", cookie.Value, err)
			}
		}

		// Recupere la session existante ou cree une nouvelle
		s := sessionStore.getSession(sid)

		// appelle le handler 
		nonoHandleFunc(s, w, req )
		}

	return handler
}

func nonoHandleFunc(s *Session, w http.ResponseWriter, req *http.Request)  {


	if s.etat == partie.PAS_COMMENCE {
		s.partie = partie.NewPartieDefault()
		s.etat = partie.EN_COURS
	}
	
	cookie := http.Cookie{
		Name: SESSION_COOKIE,
		Value: fmt.Sprintf("%d", s.sid),
		Path: "/",
	}

	http.SetCookie( w, &cookie )

	var err error
	var buf *bytes.Buffer
	if buf, err = s.partie.Html(); err !=nil {
		log.Printf( "Erreur de template : %v", err)
	}

	if _, err := w.Write( buf.Bytes() ); err!=nil {
		log.Printf("Erreur w.Write() %v\n", err)
	}

}
