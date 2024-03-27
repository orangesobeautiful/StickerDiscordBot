package qdarnt

import (
	"fmt"

	vectordatabase "backend/app/pkg/vector-database"
	searchfilter "backend/app/pkg/vector-database/search-filter"

	pb "github.com/qdrant/go-client/qdrant"
	"golang.org/x/xerrors"
)

func convertSearchFilter(filter *searchfilter.FilterInstance) (pbFilter *pb.Filter, err error) {
	pbFilter = new(pb.Filter)

	filterType := filter.GetType()
	switch filterType {
	case searchfilter.FilterTypeAnd:
		andFilterConditions := filter.GetAnd().GetConditions()
		pbConditions := make([]*pb.Condition, 0, len(andFilterConditions))
		for _, condition := range andFilterConditions {
			pbCondition, err := convertCondition(condition)
			if err != nil {
				return nil, xerrors.Errorf("convert condition: %w", err)
			}
			pbConditions = append(pbConditions, pbCondition)
		}
		pbFilter.Must = pbConditions
	default:
		return nil, vectordatabase.NewSearchFilterConvertError(fmt.Sprintf("unknow filter type: %v", filterType))
	}

	return pbFilter, nil
}

func convertCondition(condition *searchfilter.ConditionInstance) (pbCondition *pb.Condition, err error) {
	pbCondition = new(pb.Condition)

	conditionType := condition.GetType()
	switch conditionType {
	case searchfilter.ConditionTypeField:
		conditionField := condition.GetField()
		fieldKey := conditionField.GetKey()

		fieldMatch := conditionField.GetMatch()

		pbMatch, err := convertMatch(fieldMatch)
		if err != nil {
			return nil, xerrors.Errorf("convert match: %w", err)
		}

		pbCondition.ConditionOneOf = &pb.Condition_Field{
			Field: &pb.FieldCondition{
				Key:   fieldKey,
				Match: pbMatch,
			},
		}
	default:
		return nil, vectordatabase.NewSearchFilterConvertError(fmt.Sprintf("unknow condition type: %v", conditionType))
	}

	return pbCondition, nil
}

func convertMatch(match *searchfilter.MatchInstance) (pbMatch *pb.Match, err error) {
	pbMatch = new(pb.Match)

	matchType := match.GetType()
	switch matchType {
	case searchfilter.MatchTypeIntegers:
		inIntegers := match.GetInIntegers()
		pbMatch.MatchValue = &pb.Match_Integers{
			Integers: &pb.RepeatedIntegers{
				Integers: inIntegers.GetValue(),
			},
		}
	default:
		return nil, vectordatabase.NewSearchFilterConvertError(fmt.Sprintf("unknow match type: %v", matchType))
	}

	return pbMatch, nil
}
