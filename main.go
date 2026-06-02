package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/mail"
	"strings"
	"time"

	"github.com/alash3al/go-smtpsrv"
)

func main() {
	// INIT SYSTEMS
	initDB()
	startWorker()

	log.Printf("Starting SMTP Webhook Server on %s", *flagListenAddr)
	log.Printf("Forwarding emails to webhook: %s", *flagWebhook)

	cfg := smtpsrv.ServerConfig{
		ReadTimeout:     time.Duration(*flagReadTimeout) * time.Second,
		WriteTimeout:    time.Duration(*flagWriteTimeout) * time.Second,
		ListenAddr:      *flagListenAddr,
		MaxMessageBytes: int(*flagMaxMessageSize),
		BannerDomain:    *flagServerName,

		Handler: smtpsrv.HandlerFunc(func(c *smtpsrv.Context) error {

			log.Println("EMAIL RECEIVED - HANDLER TRIGGERED")

			msg, err := c.Parse()
			if err != nil {
				log.Println("PARSE ERROR:", err)
				return err
			}

			log.Println("MESSAGE PARSED SUCCESSFULLY")

			spfResult, _, _ := c.SPF()

			// CREATE EMAIL MESSAGE (SAFE INIT)
			jsonData := EmailMessage{
				ID:            msg.MessageID,
				Date:          msg.Date.String(),
				References:    msg.References,
				SPFResult:     spfResult.String(),
				ResentDate:    msg.ResentDate.String(),
				ResentID:      msg.ResentMessageID,
				Subject:       msg.Subject,
				Attachments:   []*EmailAttachment{},
				EmbeddedFiles: []*EmailEmbeddedFile{},
			}

			// SAFE INITIALIZATION (IMPORTANT FIX)
			jsonData.Body.Text = string(msg.TextBody)
			jsonData.Body.HTML = string(msg.HTMLBody)

			jsonData.Addresses.From = transformStdAddressToEmailAddress([]*mail.Address{c.From()})[0]
			jsonData.Addresses.To = transformStdAddressToEmailAddress([]*mail.Address{c.To()})[0]

			// DOMAIN CHECK
			toSplited := strings.Split(jsonData.Addresses.To.Address, "@")
			if len(*flagDomain) > 0 && (len(toSplited) < 2 || toSplited[1] != *flagDomain) {
				log.Println("DOMAIN NOT ALLOWED:", *flagDomain)
				return errors.New("unauthorized TO domain")
			}

			jsonData.Addresses.Cc = transformStdAddressToEmailAddress(msg.Cc)
			jsonData.Addresses.Bcc = transformStdAddressToEmailAddress(msg.Bcc)
			jsonData.Addresses.ReplyTo = transformStdAddressToEmailAddress(msg.ReplyTo)
			jsonData.Addresses.InReplyTo = msg.InReplyTo

			if resentFrom := transformStdAddressToEmailAddress(msg.ResentFrom); len(resentFrom) > 0 {
				jsonData.Addresses.ResentFrom = resentFrom[0]
			}

			jsonData.Addresses.ResentTo = transformStdAddressToEmailAddress(msg.ResentTo)
			jsonData.Addresses.ResentCc = transformStdAddressToEmailAddress(msg.ResentCc)
			jsonData.Addresses.ResentBcc = transformStdAddressToEmailAddress(msg.ResentBcc)

			// ATTACHMENTS
			for _, a := range msg.Attachments {
				data, _ := io.ReadAll(a.Data)

				jsonData.Attachments = append(jsonData.Attachments, &EmailAttachment{
					Filename:    a.Filename,
					ContentType: a.ContentType,
					Data:        base64.StdEncoding.EncodeToString(data),
				})
			}

			// EMBEDDED FILES
			for _, a := range msg.EmbeddedFiles {
				data, _ := io.ReadAll(a.Data)

				jsonData.EmbeddedFiles = append(jsonData.EmbeddedFiles, &EmailEmbeddedFile{
					CID:         a.CID,
					ContentType: a.ContentType,
					Data:        base64.StdEncoding.EncodeToString(data),
				})
			}

			// =========================
			// NEW ARCHITECTURE STARTS
			// =========================

			// SAVE TO DB
			id := saveEmail(
				jsonData.Addresses.From.Address,
				jsonData.Addresses.To.Address,
				jsonData.Subject,
				jsonData.Body.Text,
			)

			if id == -1 {
				log.Println("FAILED TO SAVE EMAIL TO DB")
				return errors.New("db error")
			}

			// PUSH TO WORKER QUEUE (ASYNC)
			jobQueue <- Job{
				ID:      id,
				Payload: jsonData,
			}

			log.Println("EMAIL QUEUED FOR PROCESSING, ID:", id)

			return nil
		}),
	}

	fmt.Println(smtpsrv.ListenAndServe(&cfg))
}
