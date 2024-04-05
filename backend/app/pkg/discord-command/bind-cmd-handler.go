package discordcommand

import (
	"fmt"
	"reflect"
	"runtime/debug"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/xerrors"
)

type InteractionCreateExtraParser interface {
	InteractionCreateExtraParse(*discordgo.InteractionCreate) error
}

var _ InteractionCreateExtraParser = (*BaseAuthInteractionCreate)(nil)

type BaseAuthInteractionCreate struct {
	UserID string

	GuildID string

	ChannelID string

	Name string

	NickName string

	AvatarURL string
}

func (b *BaseAuthInteractionCreate) InteractionCreateExtraParse(i *discordgo.InteractionCreate) (err error) {
	b.UserID = i.Member.User.ID
	b.GuildID = i.GuildID
	b.ChannelID = i.ChannelID
	b.Name = i.Member.User.Username
	b.NickName = i.Member.Nick
	b.AvatarURL = i.Member.User.AvatarURL("")
	return nil
}

func genDiscordCommandHandler[reqType any, respType any](
	reqParseMap parseMap, h func(reqType) (respType, error),
) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	dcCmdHandler := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var err error

		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("\npanic: %+v\n", r)
				debug.PrintStack()

				err = xerrors.Errorf("panic: %+v", r)
			}
			if err != nil {
				dcInteractionErrResponse(s, i, err)
			}
		}()

		newReq, err := newReqAndBindDiscordCommandOptions[reqType](i, reqParseMap)
		if err != nil {
			err = xerrors.Errorf("new req and bind discord command options: %w", err)
			return
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
		if err != nil {
			err = xerrors.Errorf("interaction response deferred: %w", err)
			return
		}

		err = Validate.Struct(newReq)
		if err != nil {
			if ValidateErrorConvert != nil {
				err = ValidateErrorConvert(err)
			}

			return
		}

		resp, err := h(newReq)
		if err != nil {
			err = xerrors.Errorf("handler: %w", err)
			return
		}

		err = doFinishedResponse(s, i, resp)
		if err != nil {
			err = xerrors.Errorf("do finished response: %w", err)
			return
		}
	}

	return dcCmdHandler
}

func newReqAndBindDiscordCommandOptions[reqType any](
	i *discordgo.InteractionCreate, reqParseMap parseMap,
) (newReq reqType, err error) {
	options := i.ApplicationCommandData().Options

	var reqTypeIsPtr bool
	var reqReflectValue reflect.Value
	if reflect.TypeOf(newReq).Kind() == reflect.Ptr {
		reqTypeIsPtr = true
		reqReflectValue = reflect.New(reflect.TypeOf(newReq).Elem()).Elem()
	} else {
		reqReflectValue = reflect.ValueOf(&newReq).Elem()
	}

	for _, v := range options {
		switch v.Type {
		case discordgo.ApplicationCommandOptionString:
			reqReflectValue.Field(reqParseMap[v.Name]).SetString(v.StringValue())
		case discordgo.ApplicationCommandOptionInteger:
			reqReflectValue.Field(reqParseMap[v.Name]).SetInt(v.IntValue())
		case discordgo.ApplicationCommandOptionBoolean:
			reqReflectValue.Field(reqParseMap[v.Name]).SetBool(v.BoolValue())
		case discordgo.ApplicationCommandOptionNumber:
			reqReflectValue.Field(reqParseMap[v.Name]).SetFloat(v.FloatValue())
		}
	}

	if reqTypeIsPtr {
		reqReflectValue = reqReflectValue.Addr()
		newReq = reqReflectValue.Interface().(reqType)
	}

	if extraParser, ok := any(newReq).(InteractionCreateExtraParser); ok {
		err = extraParser.InteractionCreateExtraParse(i)
		if err != nil {
			return newReq, xerrors.Errorf("extra parser: %w", err)
		}
	}

	return newReq, nil
}

func doFinishedResponse(s *discordgo.Session, i *discordgo.InteractionCreate, resp any) (err error) {
	var dcRespData *discordgo.WebhookParams
	if marshaler, ok := resp.(DiscordWebhookParamsMarshaler); ok {
		dcRespData = marshaler.MarshalDiscordWebhookParams()
	} else {
		dcRespData = &discordgo.WebhookParams{
			Content: "success",
		}
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, dcRespData)
	if err != nil {
		err = xerrors.Errorf("response followup: %w", err)
		return err
	}

	return nil
}
