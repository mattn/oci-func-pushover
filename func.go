package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	fdk "github.com/fnproject/fdk-go"
	"github.com/gregdel/pushover"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(func(ctx context.Context, in io.Reader, out io.Writer) {
		appToken := os.Getenv("PUSHOVER_APP_TOKEN")
		recepientToken := os.Getenv("PUSHOVER_RECIPIENT_TOKEN")
		if appToken == "" || recepientToken == "" {
			fdk.WriteStatus(out, 500)
			fmt.Fprintln(out, "PUSHOVER_APP_TOKEN and PUSHOVER_RECIPIENT_TOKEN environment variables must be set")
			return
		}

		b, err := io.ReadAll(in)
		if err != nil {
			fdk.WriteStatus(out, 500)
			fmt.Fprintf(out, "Error reading input: %v", err)
			return
		}

		message := pushover.NewMessage(string(b))

		message.CallbackURL = os.Getenv("PUSHOVER_CALLBACK_URL")
		message.DeviceName = os.Getenv("PUSHOVER_DEVICE_NAME")
		if expire, err := time.ParseDuration(os.Getenv("PUSHOVER_EXPIRE")); err == nil {
			message.Expire = expire
		}
		message.HTML = os.Getenv("PUSHOVER_HTML") == "true"
		switch os.Getenv("PUSHOVER_PRIORITY") {
		case "high":
			message.Priority = pushover.PriorityHigh
		case "low":
			message.Priority = pushover.PriorityLow
		case "emergency":
			message.Priority = pushover.PriorityEmergency
		case "lowest":
			message.Priority = pushover.PriorityLowest
		default:
			message.Priority = pushover.PriorityNormal
		}
		message.Monospace = os.Getenv("PUSHOVER_MONOSPACE") == "true"
		if retry, err := time.ParseDuration(os.Getenv("PUSHOVER_RETRY")); err == nil {
			message.Retry = retry
		}
		message.Sound = os.Getenv("PUSHOVER_SOUND")
		if ttl, err := time.ParseDuration(os.Getenv("PUSHOVER_TTL")); err == nil {
			message.TTL = ttl
		}
		message.Timestamp = time.Now().Unix()
		message.Title = os.Getenv("PUSHOVER_TITLE")
		message.URL = os.Getenv("PUSHOVER_URL")
		message.URLTitle = os.Getenv("PUSHOVER_URL_TITLE")

		app := pushover.New(appToken)
		response, err := app.SendMessage(
			message,
			pushover.NewRecipient(recepientToken),
		)
		if err != nil {
			fdk.WriteStatus(out, 500)
			fmt.Fprintf(out, "Error sending message: %v", err)
			return
		}
		fmt.Fprintln(out, response.String())
	}))
}
