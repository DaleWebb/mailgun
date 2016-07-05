// Package mailgun implements the github.com/itsabot/abot/interface/email/driver
// interface.
package mailgun

import (
	"bytes"
	"net/http"
	"os"

	"github.com/itsabot/abot/core"
	"github.com/itsabot/abot/core/log"
	"github.com/itsabot/abot/shared/datatypes"
	"github.com/itsabot/abot/shared/interface/email"
	"github.com/itsabot/abot/shared/interface/email/driver"
	"github.com/julienschmidt/httprouter"
	"github.com/mailgun/mailgun-go"
)

func init() {
	email.Register("mailgun", &drv{})
}

type drv struct{}

func (d *drv) Open(r *httprouter.Router) (driver.Conn, error) {
	c := &conn{
		Mailgun: mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_API_KEY"), ""),
		Domain:  os.Getenv("MAILGUN_DOMAIN"),
		ApiKey:  os.Getenv("MAILGUN_API_KEY"),
	}

	hm := dt.NewHandlerMap([]dt.RouteHandler{
		{
			// Path is prefixed by "mailgun" automatically. Thus the
			// path below becomes "/mailgun"
			Path:   "/",
			Method: "POST",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				r.ParseForm()
				err := c.Receive(r.Form["from"][0], r.Form["subject"][0], r.Form["stripped-text"][0])
				if err != nil {
					log.Debug(err)
					w.WriteHeader(500)
					w.Write(nil)
					return
				}
				w.WriteHeader(200)
				w.Write(nil)
			},
		},
	})
	hm.AddRoutes("mailgun", r)
	return c, nil
}

type conn struct {
	Mailgun mailgun.Mailgun
	Domain  string
	ApiKey  string
}

func (c *conn) Receive(from, subj, html string) error {
	body := []byte(`{"CMD":"` + html + `", "FlexID": "` + "+13105555555" /*from*/ + `", "FlexIDType": 2}`)

	req, err := http.NewRequest("POST", "nil", bytes.NewBuffer(body))
	if err != nil {
		log.Debug(err)
		return err
	}

	ret, err := core.ProcessText(req)
	if err != nil {
		log.Debug(err)
		ret = "Something went wrong with my wiring... I'll get that fixed up soon."
		return err
	}

	return c.SendPlainText([]string{from}, "Abot", "Re: "+subj, ret)
}

func (c *conn) SendHTML(to []string, from, subj, html string) error {
	return c.SendPlainText(to, from, subj, html)
}

func (c *conn) SendPlainText(to []string, from, subj, plaintext string) error {
	m := c.Mailgun.NewMessage(from+" <abot@"+c.Domain+">", subj, plaintext, to[0])

	_, _, err := c.Mailgun.Send(m)

	if err != nil {
		log.Info(err)
		return err
	}
	return nil
}

// Close the connection, but since Mailgun connections are open as needed, there
// is nothing for us to close here. Return nil.
func (c *conn) Close() error {
	return nil
}
