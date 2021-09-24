package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"os/user"

	flag "github.com/spf13/pflag"
)

var (
	host = func() string {
		if v, err := os.Hostname(); err != nil {
			return v
		}

		return "localhost"
	}()

	username = func() string {
		if user, _ := user.Current(); user != nil && user.Username != "" {
			return user.Username
		}

		return "nobody"
	}()

	fromAddr = envOrStr("MH_SENDMAIL_FROM", username+"@"+host)
	smtpAddr = envOrStr("MH_SENDMAIL_SMTP_ADDR", "localhost:1025")
)

// set via -ldflags.
var version = "development"

func main() {
	var verbose bool

	// override defaults from cli flags
	flag.StringVar(&smtpAddr, "smtp-addr", smtpAddr, "SMTP server address")
	flag.StringVarP(&fromAddr, "from", "f", fromAddr, "SMTP sender")
	flag.BoolP("long-i", "i", true, "Ignored. This flag exists for sendmail compatibility.")
	flag.BoolP("long-o", "o", true, "Ignored. This flag exists for sendmail compatibility.")
	flag.BoolP("long-t", "t", true, "Ignored. This flag exists for sendmail compatibility.")
	flag.StringP("long-N", "N", "", "Ignored. This flag exists for sendmail compatibility.")
	flag.BoolVarP(&verbose, "verbose", "v", false, "Verbose mode (sends debug output to stderr)")
	flag.Parse()

	// allow recipient to be passed as an argument
	recip := flag.Args()

	if verbose {
		fmt.Fprintln(os.Stderr, "mhsendmail", version, smtpAddr, fromAddr)
	}

	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading stdin")
		os.Exit(11)
	}

	msg, err := mail.ReadMessage(bytes.NewReader(body))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing message body")
		os.Exit(11)
	}

	if len(recip) == 0 {
		// We only need to parse the message to get a recipient if none where
		// provided on the command line.
		recip = append(recip, msg.Header.Get("To"))
	}

	err = smtp.SendMail(smtpAddr, nil, fromAddr, recip, body)
	if err != nil {
		log.Fatalf("error sending mail: %v", err)
	}
}

func envOrStr(name, fallback string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}

	return fallback
}
