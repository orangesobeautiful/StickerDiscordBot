package discordcommand

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"log/slog"

	"backend/app/domain"
	"backend/app/ent"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/xerrors"
)

type commnadMigrator struct {
	commandsRepo domain.DiscordCommandRepository

	session *discordgo.Session

	guildID string

	targetCmds []*discordgo.ApplicationCommand

	dcRegisteredCmds []*discordgo.ApplicationCommand

	savingCmds []*ent.DiscordCommand

	savingCmdNameMap map[string]*ent.DiscordCommand
}

func newCommnadMigrator(
	commandsRepo domain.DiscordCommandRepository,
	session *discordgo.Session,
	guildID string,
	targetCmds []*discordgo.ApplicationCommand,
) commnadMigrator {
	return commnadMigrator{
		commandsRepo: commandsRepo,
		session:      session,
		guildID:      guildID,
		targetCmds:   targetCmds,
	}
}

func (m *commnadMigrator) Migrate(ctx context.Context) error {
	var err error

	err = m.loadAllDCRegisteredCommand()
	if err != nil {
		return xerrors.Errorf("load all discord registered commands: %w", err)
	}

	err = m.loadAllSavingCommand(ctx)
	if err != nil {
		return xerrors.Errorf("load all saving commands: %w", err)
	}

	err = m.deleteUnnecessaryCommand()
	if err != nil {
		return xerrors.Errorf("delete unnecessary command: %w", err)
	}

	err = m.handleTargetCommands(ctx)
	if err != nil {
		return xerrors.Errorf("handle target commands: %w", err)
	}

	return nil
}

func (m *commnadMigrator) loadAllDCRegisteredCommand() error {
	var err error
	m.dcRegisteredCmds, err = m.dcCommands()
	if err != nil {
		return xerrors.Errorf("get all discord registered commands: %w", err)
	}

	return nil
}

func (m *commnadMigrator) dcCommands() ([]*discordgo.ApplicationCommand, error) {
	return m.session.ApplicationCommands(m.session.State.User.ID, m.guildID)
}

func (m *commnadMigrator) loadAllSavingCommand(ctx context.Context) error {
	var err error
	savingCmds, err := m.commandsRepo.GetAll(ctx)
	if err != nil {
		return xerrors.Errorf("get all saving commands: %w", err)
	}

	savingCmdNameMap := make(map[string]*ent.DiscordCommand, len(savingCmds))
	for _, v := range savingCmds {
		savingCmdNameMap[v.Name] = v
	}

	m.savingCmds = savingCmds
	m.savingCmdNameMap = savingCmdNameMap

	return nil
}

func (m *commnadMigrator) deleteUnnecessaryCommand() error {
	var err error

	for _, dcRegisteredCmd := range m.dcRegisteredCmds {
		if _, exist := m.savingCmdNameMap[dcRegisteredCmd.Name]; !exist {
			slog.Info("delete command: ",
				slog.String("discord_id", dcRegisteredCmd.ID),
				slog.String("name", dcRegisteredCmd.Name),
			)

			err = m.dcCommandDelete(dcRegisteredCmd.ID)
			if err != nil {
				return xerrors.Errorf("delete command: '%s': %w", dcRegisteredCmd.Name, err)
			}
		}
	}

	return nil
}

func (m *commnadMigrator) dcCommandDelete(cmdID string) error {
	return m.session.ApplicationCommandDelete(m.session.State.User.ID, m.guildID, cmdID)
}

func (m *commnadMigrator) handleTargetCommands(ctx context.Context) error {
	var err error

	for _, targetCmd := range m.targetCmds {
		err = m.handleSingleTargetCommand(ctx, targetCmd)
		if err != nil {
			return xerrors.Errorf("handle single target command: '%s': %w", targetCmd.Name, err)
		}
	}

	return nil
}

func (m *commnadMigrator) handleSingleTargetCommand(
	ctx context.Context,
	targetCmd *discordgo.ApplicationCommand,
) error {
	var err error

	if savingCmd, exist := m.savingCmdNameMap[targetCmd.Name]; exist {
		err = m.checkingCommandAndUpdate(ctx, targetCmd, savingCmd)
		if err != nil {
			return xerrors.Errorf("try update command: '%s': %w", targetCmd.Name, err)
		}

		return nil
	}

	err = m.registerCommand(ctx, targetCmd)
	if err != nil {
		return xerrors.Errorf("register command: '%s': %w", targetCmd.Name, err)
	}

	return nil
}

func (m *commnadMigrator) checkingCommandAndUpdate(
	ctx context.Context,
	targetCmd *discordgo.ApplicationCommand,
	savingCmd *ent.DiscordCommand,
) error {
	if isCommandHashEqual(targetCmd, [sha256.Size]byte(savingCmd.Sha256Checksum)) {
		return nil
	}

	targetCmd.ID = savingCmd.DiscordID

	err := m.updateCommand(ctx, targetCmd)
	if err != nil {
		return xerrors.Errorf("update command: '%s': %w", targetCmd.Name, err)
	}

	return nil
}

func (m *commnadMigrator) updateCommand(
	ctx context.Context,
	cmd *discordgo.ApplicationCommand,
) error {
	updatedCmd, err := m.dcCommandEdit(cmd)
	if err != nil {
		return xerrors.Errorf("edit discord command: %w", err)
	}

	updatedCmdChecksum := commandSha256Sum(cmd)

	err = m.commandsRepo.UpdateByName(ctx, cmd.Name, updatedCmd.ID, updatedCmdChecksum[:])
	if err != nil {
		return xerrors.Errorf("update repository command: %w", err)
	}

	return nil
}

func (m *commnadMigrator) dcCommandEdit(
	cmd *discordgo.ApplicationCommand,
) (*discordgo.ApplicationCommand, error) {
	return m.session.ApplicationCommandEdit(m.session.State.User.ID, m.guildID, cmd.ID, cmd)
}

func (m *commnadMigrator) registerCommand(ctx context.Context, cmd *discordgo.ApplicationCommand) error {
	slog.Info("register command: ", slog.String("name", cmd.Name))

	registedCmd, err := m.dcCommandCreate(cmd)
	if err != nil {
		return xerrors.Errorf("create discord command: %w", err)
	}

	cmdChecksum := commandSha256Sum(cmd)
	err = m.commandsRepo.Add(ctx, cmd.Name, registedCmd.ID, cmdChecksum[:])
	if err != nil {
		return xerrors.Errorf("add command to repository: %w", err)
	}

	return nil
}

func (m *commnadMigrator) dcCommandCreate(
	cmd *discordgo.ApplicationCommand,
) (*discordgo.ApplicationCommand, error) {
	return m.session.ApplicationCommandCreate(m.session.State.User.ID, m.guildID, cmd)
}

func isCommandHashEqual(cmd *discordgo.ApplicationCommand, checksum [sha256.Size]byte) bool {
	return commandSha256Sum(cmd) == checksum
}

func commandSha256Sum(cmd *discordgo.ApplicationCommand) [sha256.Size]byte {
	copyCmd := *cmd
	copyCmd.ID = ""

	jsonBs, _ := json.Marshal(copyCmd)

	return sha256.Sum256(jsonBs)
}
