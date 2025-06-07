package env

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// DecoderFn decodes a string value into a specific type
type DecoderFn func(raw string) (any, error)

var decoders = map[string]DecoderFn{
	"string":        func(raw string) (any, error) { return raw, nil },
	"[]uint8":       func(raw string) (any, error) { return []byte(raw), nil }, // []byte
	"[]string":      func(raw string) (any, error) { return strings.Split(raw, ","), nil },
	"int":           func(raw string) (any, error) { v, err := strconv.ParseInt(raw, 10, 64); return int(v), err },
	"int8":          func(raw string) (any, error) { v, err := strconv.ParseInt(raw, 10, 8); return int8(v), err },
	"int16":         func(raw string) (any, error) { v, err := strconv.ParseInt(raw, 10, 16); return int16(v), err },
	"int32":         func(raw string) (any, error) { v, err := strconv.ParseInt(raw, 10, 32); return int32(v), err },
	"int64":         func(raw string) (any, error) { v, err := strconv.ParseInt(raw, 10, 64); return int64(v), err },
	"uint":          func(raw string) (any, error) { v, err := strconv.ParseUint(raw, 10, 64); return uint(v), err },
	"uint8":         func(raw string) (any, error) { v, err := strconv.ParseUint(raw, 10, 8); return uint8(v), err },
	"uint16":        func(raw string) (any, error) { v, err := strconv.ParseUint(raw, 10, 16); return uint16(v), err },
	"uint32":        func(raw string) (any, error) { v, err := strconv.ParseUint(raw, 10, 32); return uint32(v), err },
	"uint64":        func(raw string) (any, error) { v, err := strconv.ParseUint(raw, 10, 64); return uint64(v), err },
	"float32":       func(raw string) (any, error) { v, err := strconv.ParseFloat(raw, 32); return float32(v), err },
	"float64":       func(raw string) (any, error) { v, err := strconv.ParseFloat(raw, 64); return float64(v), err },
	"bool":          func(raw string) (any, error) { v, err := strconv.ParseBool(raw); return bool(v), err },
	"time.Time":     func(raw string) (any, error) { return time.Parse(time.RFC3339, raw) },
	"time.Duration": func(raw string) (any, error) { return time.ParseDuration(raw) },
	"slog.Level": func(raw string) (any, error) {
		switch strings.TrimSpace(strings.ToLower(raw)) {
		case "debug":
			return slog.LevelDebug, nil
		case "warn":
			return slog.LevelWarn, nil
		case "error":
			return slog.LevelError, nil
		case "info":
			return slog.LevelInfo, nil
		default:
			return slog.LevelInfo, fmt.Errorf("invalid slog.Level: %q", raw)
		}
	},
}

// ReadStruct fills the fields of a struct with the values from the environment
// Struct tags are defined as :
// - `env:"key"`
// - `env:"key,required"` : if the environment variable is not set, an error is
// returned
func ReadStruct(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrNotPtr
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return ErrNotStructPtr
	}

	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fieldValue := rv.Field(i)

		if !fieldValue.CanSet() {
			return fmt.Errorf("field %q: %w", field.Name, ErrFieldUnexported)
		}

		decoded, err := decodeField(field)
		if errors.Is(err, ErrFieldNoEnvTag) {
			continue
		}
		if err != nil {
			return fmt.Errorf("field %q: %w", field.Name, err)
		}

		// skip decoded nil (not set and not required)
		if decoded == nil {
			continue
		}

		decodedValue := reflect.ValueOf(decoded)

		switch field.Type.Kind() {
		case reflect.Slice:
			if !decodedValue.IsValid() || decodedValue.Kind() != reflect.Slice {
				return fmt.Errorf("field %q: %w: cannot convert to slice", field.Name, ErrFieldDecode)
			}
			// Always create a new slice with the correct size to avoid index out of bounds
			fieldValue.Set(reflect.MakeSlice(field.Type, decodedValue.Len(), decodedValue.Len()))

			for i := 0; i < decodedValue.Len(); i++ {
				fieldValue.Index(i).Set(decodedValue.Index(i))
			}
		case reflect.Ptr:
			if fieldValue.IsNil() {
				fieldValue.Set(reflect.New(field.Type.Elem()))
			}
			if decodedValue.IsValid() {
				fieldValue.Elem().Set(decodedValue)
			}
		default:
			if decodedValue.IsValid() {
				fieldValue.Set(decodedValue)
			}
		}
	}

	return nil
}

func decodeField(field reflect.StructField) (any, error) {
	tag := field.Tag.Get("env")
	if tag == "" {
		return nil, ErrFieldNoEnvTag
	}

	// parse tag
	parts := strings.Split(tag, ",")
	envName := parts[0]
	required := false
	if len(parts) > 1 && parts[1] == "required" {
		required = true
	}

	// read the value
	raw, set := Read(envName)
	if !set {
		if required {
			return nil, fmt.Errorf("%w (%s)", ErrFieldRequired, envName)
		}
		return nil, nil
	}

	typeName := field.Type.String()
	if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() != reflect.Invalid {
		// For pointers, use the underlying type's decoder
		typeName = field.Type.Elem().String()
	}
	if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() != reflect.Invalid {
		typeName = `[]` + field.Type.Elem().String()
	}

	// decode
	for name, decoder := range decoders {
		if name != typeName {
			continue
		}

		decoded, err := decoder(raw)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrFieldDecode, err)
		}
		return decoded, nil
	}
	return nil, fmt.Errorf("%w: %q", ErrFieldUnsupported, typeName)
}
