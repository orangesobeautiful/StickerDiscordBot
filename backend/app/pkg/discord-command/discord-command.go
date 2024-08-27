package discordcommand

import (
	"context"

	"backend/app/domain"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/xerrors"
)

type Register interface {
	Add(command *discordgo.ApplicationCommand, handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) error
	MustAdd(command *discordgo.ApplicationCommand, handler func(s *discordgo.Session, i *discordgo.InteractionCreate))
}

type Manager interface {
	Register
	MigrateAllCommand(ctx context.Context, s *discordgo.Session, guildID string) (err error)
	GetHandler() func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var _ Manager = (*manager)(nil)

type manager struct {
	commands        []*discordgo.ApplicationCommand
	commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)

	commandsRepo domain.DiscordCommandRepository
}

func New(commandsRepo domain.DiscordCommandRepository) Manager {
	return &manager{
		commandHandlers: make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)),
		commandsRepo:    commandsRepo,
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

func (d *manager) MigrateAllCommand(ctx context.Context, s *discordgo.Session, guildID string) (err error) {
	migrator := newCommnadMigrator(d.commandsRepo, s, guildID, d.commands)

	err = migrator.Migrate(ctx)
	if err != nil {
		return xerrors.Errorf("migrate all command: %w", err)
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
