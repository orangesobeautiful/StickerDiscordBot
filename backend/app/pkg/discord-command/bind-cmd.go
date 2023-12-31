package discordcommand

import (
	"reflect"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func genDiscordApplicationCommane[reqType any](
	name, description string,
) (*discordgo.ApplicationCommand, parseMap) {
	options, reqParseMap := genDiscordApplicationCommandOptions[reqType]()

	return &discordgo.ApplicationCommand{
		Name:        name,
		Description: description,
		Options:     options,
	}, reqParseMap
}

type parseMap map[string]int

func genDiscordApplicationCommandOptions[reqType any]() ([]*discordgo.ApplicationCommandOption, parseMap) {
	var req reqType
	reqReflectType := reflect.TypeOf(req)
	if reqReflectType.Kind() != reflect.Struct {
		reqReflectType = reqReflectType.Elem()
		if reqReflectType.Kind() != reflect.Struct {
			panic("reqType must be struct")
		}
	}

	fieldNum := reqReflectType.NumField()
	options := make([]*discordgo.ApplicationCommandOption, 0, fieldNum)
	reqParseMap := make(map[string]int)
	for i := 0; i < fieldNum; i++ {
		field := reqReflectType.Field(i)

		var (
			name, description string
			required          bool
		)

		name = parseNameField(&field)
		required = parseRequiredField(&field)

		if tag := field.Tag.Get(DiscordCommandTagName); tag != "" {
			decodeResult := decodeDiscordCommandTag(tag)
			if decodeResult.name != nil {
				name = *decodeResult.name
			}
			if decodeResult.description != nil {
				description = *decodeResult.description
			}
			if decodeResult.required != nil {
				required = *decodeResult.required
			}
		}
		if description == "" {
			description = name
		}

		option := &discordgo.ApplicationCommandOption{
			Name:        name,
			Description: description,
			Required:    required,
		}

		switch field.Type.Kind() {
		case reflect.String:
			option.Type = discordgo.ApplicationCommandOptionString
		case reflect.Int:
			option.Type = discordgo.ApplicationCommandOptionInteger
		case reflect.Bool:
			option.Type = discordgo.ApplicationCommandOptionBoolean
		case reflect.Float32, reflect.Float64:
			option.Type = discordgo.ApplicationCommandOptionNumber
		default:
			panic("unsupported discord command type: " + field.Type.String())
		}

		options = append(options, option)
		reqParseMap[name] = i
	}

	return options, reqParseMap
}

func parseNameField(field *reflect.StructField) string {
	name := strings.ToLower(field.Name)
	for _, nameTagName := range ExternalNameTagNames {
		if tag := field.Tag.Get(nameTagName); tag != "" {
			name = strings.ToLower(tag)
		}
	}

	return name
}

func parseRequiredField(field *reflect.StructField) bool {
	required := false
	if tag := field.Tag.Get(ExternalValidateTagName); tag != "" {
		if strings.Contains(tag, "required") {
			required = true
		}
	}

	return required
}

type decodeDiscordCommandTagResult struct {
	name        *string
	description *string
	required    *bool
}

func decodeDiscordCommandTag(tag string) (result decodeDiscordCommandTagResult) {
	ss := strings.Split(tag, ",")
	for _, s := range ss {
		switch {
		case strings.HasPrefix(s, "name="):
			name := s[len("name="):]
			result.name = &name
		case strings.HasPrefix(s, "description="):
			description := s[len("description="):]
			result.description = &description
		case s == "required":
			required := true
			result.required = &required
		}
	}

	return result
}
