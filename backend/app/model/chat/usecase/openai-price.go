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
	openai.GPT4oMini20240718: {
		Input:  decimal.New(150, -9),
		Output: decimal.New(600, -9),
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
