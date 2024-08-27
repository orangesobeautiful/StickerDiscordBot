package server

import (
	"backend/app/domain"
	discordcommand "backend/app/pkg/discord-command"

	"github.com/go-playground/validator/v10"
)

func (s *Server) initDCCommandManager(
	validate *validator.Validate,
	eh *errHandler,
	commandsRepo domain.DiscordCommandRepository,
) {
	dcCommandManager := s.newDCCommandManager(validate, eh, commandsRepo)
	s.dcCommandManager = dcCommandManager
}

func (s *Server) newDCCommandManager(
	validate *validator.Validate,
	eh *errHandler,
	commandsRepo domain.DiscordCommandRepository,
) discordcommand.Manager {
	discordcommand.Validate = validate
	setDcCommandValidateErrConverter(eh)
	dcCommandManager := discordcommand.New(commandsRepo)

	return dcCommandManager
}

func setDcCommandValidateErrConverter(eh *errHandler) {
	bindErrConverter := eh.getBindErrConvert("en")
	discordcommand.ValidateErrorConvert = func(err error) error {
		return bindErrConverter(err)
	}
}
