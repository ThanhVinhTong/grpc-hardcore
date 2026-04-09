package serializer

import (
	"fmt"
	"os"

	"google.golang.org/protobuf/proto"
)

func WriteProtobufToJSONFile(message proto.Message, filename string) error {
	data, err := ProtobufToJSON(message)
	if err != nil {
		return fmt.Errorf("cannot marshal proto message: %w", err)
	}

	err = os.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("cannot write data to file: %w", err)
	}

	return nil
}

func WriteProtobufToBinaryFile(message proto.Message, filename string) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("cannot marshal proto message: %w", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("cannot write data to file: %w", err)
	}

	return nil
}

func ReadProtobufFromBinaryFile(filename string, message proto.Message) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("cannot read data from file: %w", err)
	}

	err = proto.Unmarshal(data, message)
	if err != nil {
		return fmt.Errorf("cannot unmarshal proto message: %w", err)
	}

	return nil
}
