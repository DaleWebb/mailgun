package mailgun

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
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
	form := url.Values{}
	form.Add("body", "Hello!")
	form.Add("sender", "dale_webb@hotmail.com")
	form.Add("recipient", "dale_webb@hotmail.com")
	form.Add("subject", "test")
	c, _ := request("POST", u, form.Encode())
	if c != 200 {
		t.Fatal("expected 200. got", c)
	}
}

func request(method, path string, data string) (int, string) {
	r, err := http.NewRequest(method, path, strings.NewReader(data))
	if err != nil {
		return 0, "err completing request: " + err.Error()
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	return w.Code, w.Body.String()
}
