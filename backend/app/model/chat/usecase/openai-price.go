package usecase

import (
	"backend/app/pkg/hserr"

	"github.com/sashabaranov/go-openai"
	"github.com/shopspring/decimal"
	"golang.org/x/xerrors"
)

type openaiPrice struct {
	Input decimal.Decimal

	Output decimal.Decimal
}

//nolint:gomnd // This is not a magic number, it corresponds to a unique openai model price
var openaiPriceMap = map[string]openaiPrice{
	GPT3Dot5Turbo0125: {
		Input:  decimal.NewFromFloat(0.0005 / 1000),
		Output: decimal.NewFromFloat(0.0015 / 1000),
	},
}

func cacluteOpenaiChatUsagePrice(
	model string, usage openai.Usage,
) (inputPrice, outputPrice decimal.Decimal, err error) {
	price, ok := openaiPriceMap[model]
	if !ok {
		err = xerrors.New("price model not found")
		return decimal.Zero, decimal.Zero, hserr.NewInternalError(err, "price model not found")
	}

	inputPrice = price.Input.Mul(decimal.NewFromInt(int64(usage.PromptTokens)))
	outputPrice = price.Output.Mul(decimal.NewFromInt(int64(usage.CompletionTokens)))

	return inputPrice, outputPrice, nil
}
