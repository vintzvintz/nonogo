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

func StaticFileHandler( path, dir string ) http.Handler {
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

	if !exists {
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

type Action struct {
	bouton string
	ligne int
	colonne int
}



func parseParamInt( req *http.Request, nom string ) int {

	s := req.FormValue( nom )
	if s=="" {
		return -1
	}

	i, err := strconv.ParseInt( s ,10, 31 )
	if err!= nil {
		return -1
	}
	return int(i)
}


func decodeAction ( req *http.Request ) (a Action) {
	
	req.ParseForm()

	btn := req.FormValue( partie.PARAM_BOUTON )
	if btn!=partie.BOUTON_DROIT && btn!=partie.BOUTON_GAUCHE {
		return a
	}

	l := parseParamInt(req, partie.PARAM_LIGNE)
	c := parseParamInt(req, partie.PARAM_COL)
	if l>=0 && c>=0 {
		a.bouton = btn
		a.ligne = l
		a.colonne=c
	}
	return a
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

		// loggue la requete
		log.Printf( "%s %s", req.Method, req.RequestURI )

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
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie( w, &cookie )

	a := decodeAction(req)
	switch a.bouton {
	case partie.BOUTON_GAUCHE: 
		s.partie.Clique(partie.JOUE_PLEIN, a.ligne, a.colonne)
	case partie.BOUTON_DROIT: 
		s.partie.Clique(partie.JOUE_VIDE, a.ligne, a.colonne)
	}

	var err error
	var buf *bytes.Buffer
	if buf, err = s.partie.Html(); err !=nil {
		log.Printf( "Erreur de template : %v", err)
	}

	if _, err := w.Write( buf.Bytes() ); err!=nil {
		log.Printf("Erreur w.Write() %v\n", err)
	}

}
