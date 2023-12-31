package discordcommand

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/xerrors"
)

type Register interface {
	Add(command *discordgo.ApplicationCommand, handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) error
	MustAdd(command *discordgo.ApplicationCommand, handler func(s *discordgo.Session, i *discordgo.InteractionCreate))
}

type Manager interface {
	Register
	RegisterAllCommand(s *discordgo.Session, guildID string) (err error)
	GetHandler() func(s *discordgo.Session, i *discordgo.InteractionCreate)
	DeleteAllCommand(s *discordgo.Session, guildID string) (err error)
}

var _ Manager = (*manager)(nil)

type manager struct {
	commands           []*discordgo.ApplicationCommand
	commandHandlers    map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	registeredCommands []*discordgo.ApplicationCommand
}

func New() Manager {
	return &manager{
		commandHandlers: make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)),
	}
}

func (d *manager) Add(
	command *discordgo.ApplicationCommand, handler func(s *discordgo.Session, i *discordgo.InteractionCreate),
) (err error) {
	if d.commandHandlers[command.Name] != nil {
		return xerrors.New("command already exists")
	}

	d.commands = append(d.commands, command)
	d.commandHandlers[command.Name] = handler

	return nil
}

func (d *manager) MustAdd(
	command *discordgo.ApplicationCommand, handler func(s *discordgo.Session, i *discordgo.InteractionCreate),
) {
	err := d.Add(command, handler)
	if err != nil {
		panic(err)
	}
}

func (d *manager) RegisterAllCommand(s *discordgo.Session, guildID string) (err error) {
	for _, v := range d.commands {
		var registedCmd *discordgo.ApplicationCommand
		registedCmd, err = s.ApplicationCommandCreate(s.State.User.ID, guildID, v)
		if err != nil {
			return xerrors.Errorf("create command: '%v': %w", v.Name, err)
		}

		d.registeredCommands = append(d.registeredCommands, registedCmd)
	}

	return nil
}

func (d *manager) GetHandler() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		handler := d.commandHandlers[i.ApplicationCommandData().Name]
		if handler == nil {
			return
		}

		handler(s, i)
	}
}

func (d *manager) DeleteAllCommand(s *discordgo.Session, guildID string) (err error) {
	for _, v := range d.registeredCommands {
		slog.Info("delete command: ", slog.String("id", v.ID), slog.String("name", v.Name))
		err = s.ApplicationCommandDelete(s.State.User.ID, guildID, v.ID)
		if err != nil {
			return xerrors.Errorf("delete command: '%v': %w", v.Name, err)
		}
	}

	return nil
}
