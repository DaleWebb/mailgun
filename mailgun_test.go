package mailgun

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/itsabot/abot/core"
	"github.com/itsabot/abot/core/log"
	"github.com/julienschmidt/httprouter"
)

var router *httprouter.Router

func TestMain(m *testing.M) {
	if len(os.Getenv("MAILGUN_DOMAIN")) == 0 ||
		len(os.Getenv("MAILGUN_API_KEY")) == 0 {
		log.Info("must set MAILGUN_DOMAIN and MAILGUN_API_KEY env vars")
		os.Exit(1)
	}
	var err error
	router, err = core.NewServer()
	if err != nil {
		log.Info("couldn't boot server", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestHandler(t *testing.T) {
	u := fmt.Sprintf("http://localhost:%s/mailgun", os.Getenv("PORT"))
	data := []byte(`{ "stripped-text": "Hello!", "sender": "dale_webb@hotmail.com", "recipient": "dale_webb@hotmail.com", "subject": "test" }`)
	c, _ := request("POST", u, data)
	if c != 200 {
		t.Fatal("expected 200. got", c)
	}
}

func request(method, path string, data []byte) (int, string) {
	r, err := http.NewRequest(method, path, bytes.NewBuffer(data))
	if err != nil {
		return 0, "err completing request: " + err.Error()
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	return w.Code, w.Body.String()
}
