package sdk

import (
	"github.com/yunis-du/actflow-agent-sdk/go-sdk/pb"
	"google.golang.org/protobuf/types/known/structpb"
)

// protoContextToContext converts proto Context to SDK Context.
func protoContextToContext(ctx *pb.Context) *Context {
	if ctx == nil {
		return &Context{
			Env:  make(map[string]string),
			Vars: make(map[string]any),
		}
	}

	vars := make(map[string]any)
	for k, v := range ctx.Vars {
		vars[k] = protoValueToAny(v)
	}

	return &Context{
		PID:  ctx.Pid,
		Env:  ctx.Env,
		Vars: vars,
	}
}

// protoValueToAny converts structpb.Value to any.
func protoValueToAny(v *structpb.Value) any {
	if v == nil {
		return nil
	}

	switch x := v.Kind.(type) {
	case *structpb.Value_NullValue:
		return nil
	case *structpb.Value_NumberValue:
		return x.NumberValue
	case *structpb.Value_StringValue:
		return x.StringValue
	case *structpb.Value_BoolValue:
		return x.BoolValue
	case *structpb.Value_StructValue:
		return protoStructToMap(x.StructValue)
	case *structpb.Value_ListValue:
		return protoListToSlice(x.ListValue)
	default:
		return nil
	}
}

// protoStructToMap converts structpb.Struct to map[string]any.
func protoStructToMap(s *structpb.Struct) map[string]any {
	if s == nil {
		return nil
	}

	m := make(map[string]any)
	for k, v := range s.Fields {
		m[k] = protoValueToAny(v)
	}
	return m
}

// protoListToSlice converts structpb.ListValue to []any.
func protoListToSlice(l *structpb.ListValue) []any {
	if l == nil {
		return nil
	}

	result := make([]any, len(l.Values))
	for i, v := range l.Values {
		result[i] = protoValueToAny(v)
	}
	return result
}

// anyToProtoValue converts any to structpb.Value.
func anyToProtoValue(v any) *structpb.Value {
	if v == nil {
		return structpb.NewNullValue()
	}

	switch x := v.(type) {
	case bool:
		return structpb.NewBoolValue(x)
	case int:
		return structpb.NewNumberValue(float64(x))
	case int32:
		return structpb.NewNumberValue(float64(x))
	case int64:
		return structpb.NewNumberValue(float64(x))
	case float32:
		return structpb.NewNumberValue(float64(x))
	case float64:
		return structpb.NewNumberValue(x)
	case string:
		return structpb.NewStringValue(x)
	case []any:
		values := make([]*structpb.Value, len(x))
		for i, item := range x {
			values[i] = anyToProtoValue(item)
		}
		return structpb.NewListValue(&structpb.ListValue{Values: values})
	case map[string]any:
		fields := make(map[string]*structpb.Value)
		for k, val := range x {
			fields[k] = anyToProtoValue(val)
		}
		return structpb.NewStructValue(&structpb.Struct{Fields: fields})
	default:
		// Try to use structpb.NewValue for other types
		val, err := structpb.NewValue(v)
		if err != nil {
			return structpb.NewNullValue()
		}
		return val
	}
}

// outputToProto converts SDK Output to proto AgentOutput.
func outputToProto(output *Output) *pb.AgentOutput {
	if output == nil {
		return &pb.AgentOutput{
			Status: pb.NodeExecutionStatus_FAILED,
			Error:  "nil output",
		}
	}

	return &pb.AgentOutput{
		Status:    pb.NodeExecutionStatus(output.Status),
		Outputs:   anyToProtoValue(output.Outputs),
		Error:     output.Error,
		Exception: output.Exception,
	}
}

