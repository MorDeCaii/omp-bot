package certificate

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ozonmp/omp-bot/internal/service/loyalty/certificate"
	"log"
	"strconv"
)

func (c *LoyaltyCertificateCommander) New(inputMessage *tgbotapi.Message) {
	args := inputMessage.CommandArguments()

	outputMsg := ""
	parsedData := CertificateData{}
	err := json.Unmarshal([]byte(args), &parsedData)
	if err != nil {
		outputMsg = "Pass valid JSON serialized data as parameter"
		msg := tgbotapi.NewMessage(
			inputMessage.Chat.ID,
			outputMsg,
		)

		_, err = c.bot.Send(msg)
		if err != nil {
			log.Printf("LoyaltyCertificateCommander.New: error sending reply message to chat - %v", err)
		}
		return
	}
	lastIndex := c.certificateService.Certificates[len(c.certificateService.Certificates) - 1].Id

	newCertificate := certificate.Certificate{
		Id:          lastIndex + 1,
		SellerTitle: parsedData.SellerTitle,
		Amount:      parsedData.Amount,
		ExpireDate:  parsedData.ExpireDate,
	}

	newId, err := c.certificateService.Create(newCertificate)
	outputMsg = "Certificate with ID " + strconv.Itoa(int(newId))
	if err != nil {
		log.Printf("failed to create certificate with id %d: %v", newCertificate.Id, err)
		outputMsg += " already exists"
	} else {
		outputMsg += " was created"
	}

	msg := tgbotapi.NewMessage(
		inputMessage.Chat.ID,
		outputMsg,
	)

	_, err = c.bot.Send(msg)
	if err != nil {
		log.Printf("LoyaltyCertificateCommander.New: error sending reply message to chat - %v", err)
	}
}