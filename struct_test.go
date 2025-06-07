package env_test

import (
	"encoding/json"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xdrm-io/env"
)

func TestReadStruct(t *testing.T) {
	type nominal struct {
		Field1  string        `env:"FIELD1"`
		Field2  []uint8       `env:"FIELD2"`
		Field3  []string      `env:"FIELD3"`
		Field4  int           `env:"FIELD4"`
		Field5  int8          `env:"FIELD5"`
		Field6  int16         `env:"FIELD6"`
		Field7  int32         `env:"FIELD7"`
		Field8  int64         `env:"FIELD8"`
		Field9  uint          `env:"FIELD9"`
		Field10 uint8         `env:"FIELD10"`
		Field11 uint16        `env:"FIELD11"`
		Field12 uint32        `env:"FIELD12"`
		Field13 uint64        `env:"FIELD13"`
		Field14 float32       `env:"FIELD14"`
		Field15 float64       `env:"FIELD15"`
		Field16 bool          `env:"FIELD16"`
		Field17 time.Time     `env:"FIELD17"`
		Field18 time.Duration `env:"FIELD18"`
		Field19 slog.Level    `env:"FIELD19"`
	}

	type customUnexported struct {
		field string `env:"VARNAME"`
	}
	type required struct {
		Field string `env:"VARNAME,required"`
	}

	type customString struct {
		Field string `env:"VARNAME"`
	}
	type customBytes struct {
		Field []byte `env:"VARNAME"`
	}
	type customStrings struct {
		Field []string `env:"VARNAME"`
	}
	type customInt struct {
		Field int `env:"VARNAME"`
	}
	type customInt8 struct {
		Field int8 `env:"VARNAME"`
	}
	type customInt16 struct {
		Field int16 `env:"VARNAME"`
	}
	type customInt32 struct {
		Field int32 `env:"VARNAME"`
	}
	type customInt64 struct {
		Field int64 `env:"VARNAME"`
	}
	type customUint8 struct {
		Field uint8 `env:"VARNAME"`
	}
	type customUint16 struct {
		Field uint16 `env:"VARNAME"`
	}
	type customUint32 struct {
		Field uint32 `env:"VARNAME"`
	}
	type customUint64 struct {
		Field uint64 `env:"VARNAME"`
	}
	type customFloat32 struct {
		Field float32 `env:"VARNAME"`
	}
	type customFloat64 struct {
		Field float64 `env:"VARNAME"`
	}
	type customBool struct {
		Field bool `env:"VARNAME"`
	}
	type customTime struct {
		Field time.Time `env:"VARNAME"`
	}
	type customDuration struct {
		Field time.Duration `env:"VARNAME"`
	}
	type customSlogLevel struct {
		Field slog.Level `env:"VARNAME"`
	}
	type unsupported struct {
		Field []any `env:"VARNAME"`
	}
	tt := []struct {
		name     string
		receiver any
		env      map[string]string
		expect   any
		err      error
	}{
		{
			name:     "nil",
			receiver: nil,
			err:      env.ErrNotPtr,
		},
		{
			name:     "not a pointer",
			receiver: struct{}{},
			err:      env.ErrNotPtr,
		},
		{
			name:     "not a pointer to struct",
			receiver: &[]int{},
			err:      env.ErrNotStructPtr,
		},
		{
			name:     "unexported field fails",
			receiver: &customUnexported{},
			err:      env.ErrFieldUnexported,
		},
		{
			name:     "required field ok",
			receiver: &required{},
			env:      map[string]string{"VARNAME": "value"},
			expect:   required{Field: "value"},
		},
		{
			name:     "required field fails",
			receiver: &required{},
			env:      map[string]string{},
			err:      env.ErrFieldRequired,
		},
		{
			name:     "string field ok",
			receiver: &customString{},
			env:      map[string]string{"VARNAME": "value"},
			expect:   customString{Field: "value"},
		},
		{
			name:     "[]byte field ok",
			receiver: &customBytes{},
			env:      map[string]string{"VARNAME": "value"},
			expect:   customBytes{Field: []byte("value")},
		},
		{
			name:     "[]string field ok",
			receiver: &customStrings{},
			env:      map[string]string{"VARNAME": "value1,value2"},
			expect:   customStrings{Field: []string{"value1", "value2"}},
		},
		{
			name:     "int field ok",
			receiver: &customInt{},
			env:      map[string]string{"VARNAME": "-13"},
			expect:   customInt{Field: -13},
		},
		{
			name:     "int field fail",
			receiver: &customInt{},
			env:      map[string]string{"VARNAME": "a"},
			err:      env.ErrFieldDecode,
		},
		{
			name:     "int8 field ok",
			receiver: &customInt8{},
			env:      map[string]string{"VARNAME": "1"},
			expect:   customInt8{Field: 1},
		},
		{
			name:     "int8 field fail",
			receiver: &customInt8{},
			env:      map[string]string{"VARNAME": "a"},
			err:      env.ErrFieldDecode,
		},
		{
			name:     "int16 field ok",
			receiver: &customInt16{},
			env:      map[string]string{"VARNAME": "1"},
			expect:   customInt16{Field: 1},
		},
		{
			name:     "int16 field fail",
			receiver: &customInt16{},
			env:      map[string]string{"VARNAME": "a"},
			err:      env.ErrFieldDecode,
		},
		{
			name:     "int32 field ok",
			receiver: &customInt32{},
			env:      map[string]string{"VARNAME": "1"},
			expect:   customInt32{Field: 1},
		},
		{
			name:     "int32 field fail",
			receiver: &customInt32{},
			env:      map[string]string{"VARNAME": "a"},
			err:      env.ErrFieldDecode,
		},
		{
			name:     "int64 field ok",
			receiver: &customInt64{},
			env:      map[string]string{"VARNAME": "1"},
			expect:   customInt64{Field: 1},
		},
		{
			name:     "int64 field fail",
			receiver: &customInt64{},
			env:      map[string]string{"VARNAME": "a"},
			err:      env.ErrFieldDecode,
		},
		{
			name:     "uint8 field ok",
			receiver: &customUint8{},
			env:      map[string]string{"VARNAME": "1"},
			expect:   customUint8{Field: 1},
		},
		{
			name:     "uint8 field fail",
			receiver: &customUint8{},
			env:      map[string]string{"VARNAME": "a"},
			err:      env.ErrFieldDecode,
		},
		{
			name:     "uint16 field ok",
			receiver: &customUint16{},
			env:      map[string]string{"VARNAME": "1"},
			expect:   customUint16{Field: 1},
		},
		{
			name:     "uint16 field fail",
			receiver: &customUint16{},
			env:      map[string]string{"VARNAME": "a"},
			err:      env.ErrFieldDecode,
		},
		{
			name:     "uint32 field ok",
			receiver: &customUint32{},
			env:      map[string]string{"VARNAME": "1"},
			expect:   customUint32{Field: 1},
		},
		{
			name:     "uint32 field fail",
			receiver: &customUint32{},
			env:      map[string]string{"VARNAME": "a"},
			err:      env.ErrFieldDecode,
		},
		{
			name:     "uint64 field ok",
			receiver: &customUint64{},
			env:      map[string]string{"VARNAME": "1"},
			expect:   customUint64{Field: 1},
		},
		{
			name:     "uint64 field fail",
			receiver: &customUint64{},
			env:      map[string]string{"VARNAME": "a"},
			err:      env.ErrFieldDecode,
		},
		{
			name:     "float32 field ok",
			receiver: &customFloat32{},
			env:      map[string]string{"VARNAME": "1.23"},
			expect:   customFloat32{Field: 1.23},
		},
		{
			name:     "float32 field fail",
			receiver: &customFloat32{},
			env:      map[string]string{"VARNAME": "a"},
			err:      env.ErrFieldDecode,
		},
		{
			name:     "float64 field ok",
			receiver: &customFloat64{},
			env:      map[string]string{"VARNAME": "1.23"},
			expect:   customFloat64{Field: 1.23},
		},
		{
			name:     "float64 field fail",
			receiver: &customFloat64{},
			env:      map[string]string{"VARNAME": "a"},
			err:      env.ErrFieldDecode,
		},
		{
			name:     "bool field ok",
			receiver: &customBool{},
			env:      map[string]string{"VARNAME": "true"},
			expect:   customBool{Field: true},
		},
		{
			name:     "bool field fail",
			receiver: &customBool{},
			env:      map[string]string{"VARNAME": "a"},
			err:      env.ErrFieldDecode,
		},
		{
			name:     "time.Time field ok",
			receiver: &customTime{},
			env:      map[string]string{"VARNAME": "2025-01-01T00:00:00Z"},
			expect:   customTime{Field: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
		{
			name:     "time.Time field fail",
			receiver: &customTime{},
			env:      map[string]string{"VARNAME": "a"},
			err:      env.ErrFieldDecode,
		},
		{
			name:     "time.Duration field ok",
			receiver: &customDuration{},
			env:      map[string]string{"VARNAME": "2h"},
			expect:   customDuration{Field: 2 * time.Hour},
		},
		{
			name:     "time.Duration field fail",
			receiver: &customDuration{},
			env:      map[string]string{"VARNAME": "a"},
			err:      env.ErrFieldDecode,
		},
		{
			name:     "slog.Level field ok",
			receiver: &customSlogLevel{},
			env:      map[string]string{"VARNAME": "info"},
			expect:   customSlogLevel{Field: slog.LevelInfo},
		},
		{
			name:     "slog.Level field fail",
			receiver: &customSlogLevel{},
			env:      map[string]string{"VARNAME": "a"},
			err:      env.ErrFieldDecode,
		},
		{
			name:     "unsupported field type fails",
			receiver: &unsupported{},
			env:      map[string]string{"VARNAME": "value"},
			err:      env.ErrFieldUnsupported,
		},
		{
			name:     "nominal ok",
			receiver: &nominal{},
			env: map[string]string{
				"FIELD1":  "value",
				"FIELD2":  "value",
				"FIELD3":  "value1,value2",
				"FIELD4":  "-13",
				"FIELD5":  "-1",
				"FIELD6":  "-1",
				"FIELD7":  "-1",
				"FIELD8":  "-1",
				"FIELD9":  "1",
				"FIELD10": "1",
				"FIELD11": "1",
				"FIELD12": "1",
				"FIELD13": "1",
				"FIELD14": "1.23",
				"FIELD15": "-1.23",
				"FIELD16": "true",
				"FIELD17": "2025-01-01T00:00:00Z",
				"FIELD18": "2h",
				"FIELD19": "info",
			},
			expect: nominal{
				Field1:  "value",
				Field2:  []uint8("value"),
				Field3:  []string{"value1", "value2"},
				Field4:  -13,
				Field5:  -1,
				Field6:  -1,
				Field7:  -1,
				Field8:  -1,
				Field9:  1,
				Field10: 1,
				Field11: 1,
				Field12: 1,
				Field13: 1,
				Field14: 1.23,
				Field15: -1.23,
				Field16: true,
				Field17: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Field18: 2 * time.Hour,
				Field19: slog.LevelInfo,
			},
			err: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tc.env {
				os.Setenv(k, v)
			}

			err := env.ReadStruct(tc.receiver)
			if tc.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)

			actual, err := json.Marshal(tc.receiver)
			require.NoError(t, err)
			expect, err := json.Marshal(tc.expect)
			require.NoError(t, err)
			require.JSONEq(t, string(expect), string(actual))
		})
	}
}

// Test case for slice size mismatch bug fix
func TestReadStruct_SliceSizeMismatch(t *testing.T) {
	type testStruct struct {
		Tags []string `env:"TAGS"`
	}

	// Pre-initialize with different size to trigger the bug
	config := &testStruct{
		Tags: []string{"existing", "values", "that", "should", "be", "replaced"},
	}

	os.Clearenv()
	os.Setenv("TAGS", "new,values")

	err := env.ReadStruct(config)
	require.NoError(t, err)
	require.Equal(t, []string{"new", "values"}, config.Tags)
}

// Test case for nil decoded value bug fix
func TestReadStruct_NilDecodedValue(t *testing.T) {
	type testStruct struct {
		OptionalField string `env:"OPTIONAL_FIELD"`
		RequiredField string `env:"REQUIRED_FIELD,required"`
	}

	config := &testStruct{}

	os.Clearenv()
	os.Setenv("REQUIRED_FIELD", "required_value")
	// OPTIONAL_FIELD is not set

	err := env.ReadStruct(config)
	require.NoError(t, err)
	require.Equal(t, "", config.OptionalField) // Should remain empty/zero value
	require.Equal(t, "required_value", config.RequiredField)
}

// Test case for pointer fields with nil values
func TestReadStruct_PointerNilValue(t *testing.T) {
	type testStruct struct {
		OptionalPtr *string `env:"OPTIONAL_PTR"`
		RequiredPtr *string `env:"REQUIRED_PTR,required"`
	}

	config := &testStruct{}

	os.Clearenv()
	os.Setenv("REQUIRED_PTR", "required_value")
	// OPTIONAL_PTR is not set

	err := env.ReadStruct(config)
	require.NoError(t, err)
	require.Nil(t, config.OptionalPtr) // Should remain nil
	require.NotNil(t, config.RequiredPtr)
	require.Equal(t, "required_value", *config.RequiredPtr)
}

// Test case for slice fields with nil values
func TestReadStruct_SliceNilValue(t *testing.T) {
	type testStruct struct {
		OptionalSlice []string `env:"OPTIONAL_SLICE"`
		RequiredSlice []string `env:"REQUIRED_SLICE,required"`
	}

	config := &testStruct{}

	os.Clearenv()
	os.Setenv("REQUIRED_SLICE", "val1,val2")
	// OPTIONAL_SLICE is not set

	err := env.ReadStruct(config)
	require.NoError(t, err)
	require.Nil(t, config.OptionalSlice) // Should remain nil
	require.Equal(t, []string{"val1", "val2"}, config.RequiredSlice)
}
