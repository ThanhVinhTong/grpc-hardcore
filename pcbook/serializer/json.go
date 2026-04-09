package serializer

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func ProtobufToJSON(message proto.Message) (string, error) {
	marshaler := protojson.MarshalOptions{
		UseProtoNames:     true,
		UseEnumNumbers:    true,
		EmitDefaultValues: true,
		Indent:            "  ",
	}

	data, err := marshaler.Marshal(message)
	if err != nil {
		return "", fmt.Errorf("cannot marshal proto message: %w", err)
	}

	return string(data), nil
}
