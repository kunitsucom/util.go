package protoext

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type EnumStringer interface {
	Number() protoreflect.EnumNumber
	Descriptor() protoreflect.EnumDescriptor
}

// EnumValueOptionsString is a helper function to get the string of google.protobuf.EnumValueOptions set in the enum of proto.
// Get the string `male` from the enum `SEX_MALE` for a .proto file like the following:
//
//	extend google.protobuf.EnumValueOptions {
//	  optional string enum_stringer = 50000;
//	}
//
//	enum Sex {
//	  SEX_UNSPECIFIED = 0 [(enum_stringer) = "unspecified"];
//	  SEX_MALE = 1 [(enum_stringer) = "male"];
//	  SEX_FEMALE = 2 [(enum_stringer) = "female"];
//	  SEX_OTHER = 3 [(enum_stringer) = "other"];
//	}
func EnumValueOptionsString(xt protoreflect.ExtensionType, e EnumStringer) string {
	ext, ok := proto.GetExtension(e.Descriptor().Values().Get(int(e.Number())).Options(), xt).(string)

	if !ok {
		return ""
	}

	return ext
}
