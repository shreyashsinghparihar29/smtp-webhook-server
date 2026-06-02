package main

import "flag"

var (
	flagServerName = flag.String(
		"name",
		"smtp-webhook-server",
		"SMTP server name",
	)

	flagListenAddr = flag.String(
		"listen",
		":2525",
		"SMTP address to listen on",
	)

	flagWebhook = flag.String(
		"webhook",
		"http://localhost:8080/my/webhook",
		"Webhook endpoint to receive email payloads",
	)

	flagMaxMessageSize = flag.Int64(
		"msglimit",
		1024*1024*2,
		"Maximum incoming message size in bytes",
	)

	flagReadTimeout = flag.Int(
		"timeout.read",
		60,
		"Read timeout in seconds",
	)

	flagWriteTimeout = flag.Int(
		"timeout.write",
		60,
		"Write timeout in seconds",
	)

	flagAuthUSER = flag.String(
		"user",
		"",
		"SMTP authentication username",
	)

	flagAuthPASS = flag.String(
		"pass",
		"",
		"SMTP authentication password",
	)

	flagDomain = flag.String(
		"domain",
		"",
		"Allowed recipient domain",
	)
)

func init() {
	flag.Parse()
}
