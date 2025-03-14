package backends

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/freakmaxi/go-guerrilla/mail"
)

// ----------------------------------------------------------------------------------
// Processor Name: hasher
// ----------------------------------------------------------------------------------
// Description   : Generates a unique md5 checksum id for an email
// ----------------------------------------------------------------------------------
// Config Options: None
// --------------:-------------------------------------------------------------------
// Input         : e.MailFrom, e.Subject, e.RcptTo
//
//	: assuming e.Subject was generated by "headersparser" processor
//
// ----------------------------------------------------------------------------------
// Output        : Checksum stored in e.Hash
// ----------------------------------------------------------------------------------
func init() {
	processors["hasher"] = func() Decorator {
		return Hasher()
	}
}

// The hasher decorator computes a hash of the email for each recipient
// It appends the hashes to envelope's Hashes slice.
func Hasher() Decorator {
	return func(p Processor) Processor {
		return ProcessWith(func(e *mail.Envelope, task SelectTask) (Result, error) {

			if task == TaskSaveMail {
				// base hash, use subject from and timestamp-nano
				h := md5.New()
				ts := fmt.Sprintf("%d", time.Now().UnixNano())
				_, _ = io.Copy(h, strings.NewReader(e.MailFrom.String()))
				_, _ = io.Copy(h, strings.NewReader(e.Subject))
				_, _ = io.Copy(h, strings.NewReader(ts))
				// using the base hash, calculate a unique hash for each recipient
				for i := range e.RcptTo {
					h2 := h
					_, _ = io.Copy(h2, strings.NewReader(e.RcptTo[i].String()))
					sum := h2.Sum([]byte{})
					e.Hashes = append(e.Hashes, fmt.Sprintf("%x", sum))
				}
				return p.Process(e, task)
			} else {
				return p.Process(e, task)
			}

		})
	}
}
