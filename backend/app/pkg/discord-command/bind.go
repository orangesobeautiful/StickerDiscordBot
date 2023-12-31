package discordcommand

import (
	"github.com/bwmarrin/discordgo"
	"github.com/go-playground/locales/en_US"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	DiscordCommandTagName   = "dccmd"
	ExternalNameTagNames    = []string{"json"}
	ExternalValidateTagName = "binding"

	enLocalesTranslator = en_US.New()
	UniTranslator       = ut.New(enLocalesTranslator)

	Validate = validator.New(validator.WithRequiredStructEnabled())
)

type DiscordWebhookParamsMarshaler interface {
	MarshalDiscordWebhookParams() *discordgo.WebhookParams
}

func DiscordCommandRegister[reqType any, respType any](name, description string,
	h func(reqType) (respType, error),
) (
	discordApplicationCommand *discordgo.ApplicationCommand,
	discordCommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate),
) {
	dcCmd, reqParseMap := genDiscordApplicationCommane[reqType](name, description)
	dcCmdHandler := genDiscordCommandHandler(reqParseMap, h)
	return dcCmd, dcCmdHandler
}
