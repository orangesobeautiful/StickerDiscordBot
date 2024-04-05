package qdarnt

import (
	"encoding/base64"
	"unicode/utf8"

	pb "github.com/qdrant/go-client/qdrant"
	"golang.org/x/exp/constraints"
	"google.golang.org/protobuf/runtime/protoimpl"
)

func NewPayload(v map[string]any) (map[string]*pb.Value, error) {
	x := make(map[string]*pb.Value, len(v))
	for k, v := range v {
		if !utf8.ValidString(k) {
			return nil, protoimpl.X.NewError("invalid UTF-8 in string: %q", k)
		}
		var err error
		x[k], err = NewValue(v)
		if err != nil {
			return nil, err
		}
	}

	return x, nil
}

// NewValue
//
//nolint:gocyclo
func NewValue(v any) (*pb.Value, error) {
	switch v := v.(type) {
	case nil:
		return NewNullValue(), nil
	case bool:
		return NewBoolValue(v), nil
	case int:
		return NewIntegerValue(v), nil
	case int32:
		return NewIntegerValue(v), nil
	case int64:
		return NewIntegerValue(v), nil
	case uint:
		return NewIntegerValue(v), nil
	case uint32:
		return NewIntegerValue(v), nil
	case uint64:
		return NewIntegerValue(v), nil
	case float32:
		return NewFloatValue(v), nil
	case float64:
		return NewFloatValue(v), nil
	case string:
		if !utf8.ValidString(v) {
			return nil, protoimpl.X.NewError("invalid UTF-8 in string: %q", v)
		}
		return NewStringValue(v), nil
	case []byte:
		s := base64.StdEncoding.EncodeToString(v)
		return NewStringValue(s), nil
	case map[string]any:
		v2, err := NewStruct(v)
		if err != nil {
			return nil, err
		}
		return NewStructValue(v2), nil
	case []any:
		v2, err := NewList(v)
		if err != nil {
			return nil, err
		}
		return NewListValue(v2), nil
	default:
		return nil, protoimpl.X.NewError("invalid type: %T", v)
	}
}

func NewNullValue() *pb.Value {
	return &pb.Value{
		Kind: &pb.Value_NullValue{
			NullValue: pb.NullValue_NULL_VALUE,
		},
	}
}

func NewBoolValue(v bool) *pb.Value {
	return &pb.Value{
		Kind: &pb.Value_BoolValue{
			BoolValue: v,
		},
	}
}

func NewIntegerValue[T constraints.Integer](v T) *pb.Value {
	return &pb.Value{
		Kind: &pb.Value_IntegerValue{
			IntegerValue: int64(v),
		},
	}
}

func NewFloatValue[T constraints.Float](v T) *pb.Value {
	return &pb.Value{
		Kind: &pb.Value_DoubleValue{
			DoubleValue: float64(v),
		},
	}
}

func NewStringValue(v string) *pb.Value {
	return &pb.Value{
		Kind: &pb.Value_StringValue{
			StringValue: v,
		},
	}
}

func NewStruct(v map[string]any) (*pb.Struct, error) {
	x := &pb.Struct{Fields: make(map[string]*pb.Value, len(v))}
	for k, v := range v {
		if !utf8.ValidString(k) {
			return nil, protoimpl.X.NewError("invalid UTF-8 in string: %q", k)
		}
		var err error
		x.Fields[k], err = NewValue(v)
		if err != nil {
			return nil, err
		}
	}
	return x, nil
}

func NewStructValue(v *pb.Struct) *pb.Value {
	return &pb.Value{
		Kind: &pb.Value_StructValue{
			StructValue: v,
		},
	}
}

func NewList(v []any) (*pb.ListValue, error) {
	x := &pb.ListValue{Values: make([]*pb.Value, len(v))}
	for i, v := range v {
		var err error
		x.Values[i], err = NewValue(v)
		if err != nil {
			return nil, err
		}
	}
	return x, nil
}

func NewListValue(v *pb.ListValue) *pb.Value {
	return &pb.Value{
		Kind: &pb.Value_ListValue{
			ListValue: v,
		},
	}
}
