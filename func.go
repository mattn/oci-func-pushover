package main

import (
	"context"
	"fmt"
	"io"
	"os"

	fdk "github.com/fnproject/fdk-go"
	"github.com/gregdel/pushover"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(func(ctx context.Context, in io.Reader, out io.Writer) {
		appToken := os.Getenv("APP_TOKEN")
		recepientToken := os.Getenv("RECIPIENT_TOKEN")
		if appToken == "" || recepientToken == "" {
			fdk.WriteStatus(out, 500)
			fmt.Fprintln(out, "APP_TOKEN and RECIPIENT_TOKEN environment variables must be set")
			return
		}

		b, err := io.ReadAll(in)
		if err != nil {
			fdk.WriteStatus(out, 500)
			fmt.Fprintf(out, "Error reading input: %v", err)
			return
		}

		app := pushover.New(appToken)
		response, err := app.SendMessage(
			pushover.NewMessage(string(b)),
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
