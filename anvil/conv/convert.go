// Copyright (c) 2025 The anyvoxel Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package conv provide the helper function to convert data types & string
package conv

import (
	"context"
	"reflect"
	"strconv"
	"time"

	"github.com/anyvoxel/airmid/anvil/xerrors"
)

// ConvertToBool converts []string to bool. Only the first data is used.
func ConvertToBool(_ context.Context, data []string) (any, error) {
	origin := data[0]
	target, err := strconv.ParseBool(origin)
	if err != nil {
		return nil, xerrors.Wrapf(err, "Cann't convert %s to type bool", origin)
	}
	return target, nil
}

// ConvertToInt converts []string to int. Only the first data is used.
func ConvertToInt(_ context.Context, data []string) (any, error) {
	origin := data[0]
	target, err := strconv.ParseInt(origin, 10, 0)
	if err != nil {
		return nil, xerrors.Wrapf(err, "Cann't convert %s to type int", origin)
	}
	return int(target), nil
}

// ConvertToInt8 converts []string to int8. Only the first data is used.
func ConvertToInt8(_ context.Context, data []string) (any, error) {
	origin := data[0]
	target, err := strconv.ParseInt(origin, 10, 8)
	if err != nil {
		return nil, xerrors.Wrapf(err, "Cann't convert %s to type int8", origin)
	}
	return int8(target), nil
}

// ConvertToInt16 converts []string to int16. Only the first data is used.
func ConvertToInt16(_ context.Context, data []string) (any, error) {
	origin := data[0]
	target, err := strconv.ParseInt(origin, 10, 16)
	if err != nil {
		return nil, xerrors.Wrapf(err, "Cann't convert %s to type int16", origin)
	}
	return int16(target), nil
}

// ConvertToInt32 converts []string to int32. Only the first data is used.
func ConvertToInt32(_ context.Context, data []string) (any, error) {
	origin := data[0]
	target, err := strconv.ParseInt(origin, 10, 32)
	if err != nil {
		return nil, xerrors.Wrapf(err, "Cann't convert %s to type int32", origin)
	}
	return int32(target), nil
}

// ConvertToInt64 converts []string to int64. Only the first data is used.
func ConvertToInt64(_ context.Context, data []string) (any, error) {
	origin := data[0]
	target, err := strconv.ParseInt(origin, 10, 64)
	if err != nil {
		return nil, xerrors.Wrapf(err, "Cann't convert %s to type int64", origin)
	}
	return target, nil
}

// ConvertToUint converts []string to uint. Only the first data is used.
func ConvertToUint(_ context.Context, data []string) (any, error) {
	origin := data[0]
	target, err := strconv.ParseUint(origin, 10, 0)
	if err != nil {
		return nil, xerrors.Wrapf(err, "Cann't convert %s to type uint", origin)
	}
	return uint(target), nil
}

// ConvertToUint8 converts []string to uint8. Only the first data is used.
func ConvertToUint8(_ context.Context, data []string) (any, error) {
	origin := data[0]
	target, err := strconv.ParseUint(origin, 10, 8)
	if err != nil {
		return nil, xerrors.Wrapf(err, "Cann't convert %s to type uint8", origin)
	}
	return uint8(target), nil
}

// ConvertToUint16 converts []string to uint16. Only the first data is used.
func ConvertToUint16(_ context.Context, data []string) (any, error) {
	origin := data[0]
	target, err := strconv.ParseUint(origin, 10, 16)
	if err != nil {
		return nil, xerrors.Wrapf(err, "Cann't convert %s to type uint16", origin)
	}
	return uint16(target), nil
}

// ConvertToUint32 converts []string to uint32. Only the first data is used.
func ConvertToUint32(_ context.Context, data []string) (any, error) {
	origin := data[0]
	target, err := strconv.ParseUint(origin, 10, 32)
	if err != nil {
		return nil, xerrors.Wrapf(err, "Cann't convert %s to type uint32", origin)
	}
	return uint32(target), nil
}

// ConvertToUint64 converts []string to uint64. Only the first data is used.
func ConvertToUint64(_ context.Context, data []string) (any, error) {
	origin := data[0]
	target, err := strconv.ParseUint(origin, 10, 64)
	if err != nil {
		return nil, xerrors.Wrapf(err, "Cann't convert %s to type uint64", origin)
	}
	return target, nil
}

// ConvertToFloat32 converts []string to float32. Only the first data is used.
func ConvertToFloat32(_ context.Context, data []string) (any, error) {
	origin := data[0]
	target, err := strconv.ParseFloat(origin, 32)
	if err != nil {
		return nil, xerrors.Wrapf(err, "Cann't convert %s to type float32", origin)
	}
	return float32(target), nil
}

// ConvertToFloat64 converts []string to float64. Only the first data is used.
func ConvertToFloat64(_ context.Context, data []string) (any, error) {
	origin := data[0]
	target, err := strconv.ParseFloat(origin, 64)
	if err != nil {
		return nil, xerrors.Wrapf(err, "Cann't convert %s to type float64", origin)
	}
	return target, nil
}

// ConvertToString return the first element in []string.
func ConvertToString(_ context.Context, data []string) (any, error) {
	return data[0], nil
}

// ConvertToTimeDuration return the first element in []string as time.Duration.
func ConvertToTimeDuration(_ context.Context, data []string) (any, error) {
	return time.ParseDuration(data[0])
}

// ConvertTo return data to typ.
//
//nolint:revive,exhaustive,cyclop
func ConvertTo(ctx context.Context, typ reflect.Type, data []string) (any, error) {
	// NOTE: we must check the typ first, otherwise if the type is time.Duration
	// the kind will be reflect.Int64
	if typ == reflect.TypeOf(time.Duration(0)) {
		return ConvertToTimeDuration(ctx, data)
	}

	switch typ.Kind() {
	case reflect.Slice:
		ret := reflect.MakeSlice(typ, 0, len(data))
		for _, v := range data {
			i, err := ConvertTo(ctx, typ.Elem(), []string{v})
			if err != nil {
				return nil, err
			}

			ret = reflect.Append(ret, reflect.ValueOf(i))
		}
		return ret.Interface(), nil
	case reflect.Bool:
		return ConvertToBool(ctx, data)
	case reflect.Int:
		return ConvertToInt(ctx, data)
	case reflect.Int8:
		return ConvertToInt8(ctx, data)
	case reflect.Int16:
		return ConvertToInt16(ctx, data)
	case reflect.Int32:
		return ConvertToInt32(ctx, data)
	case reflect.Int64:
		return ConvertToInt64(ctx, data)
	case reflect.Uint:
		return ConvertToUint(ctx, data)
	case reflect.Uint8:
		return ConvertToUint8(ctx, data)
	case reflect.Uint16:
		return ConvertToUint16(ctx, data)
	case reflect.Uint32:
		return ConvertToUint32(ctx, data)
	case reflect.Uint64:
		return ConvertToUint64(ctx, data)
	case reflect.Float32:
		return ConvertToFloat32(ctx, data)
	case reflect.Float64:
		return ConvertToFloat64(ctx, data)
	case reflect.String:
		return ConvertToString(ctx, data)
	}

	return nil, xerrors.Errorf("Unsupport target type %s", typ.String())
}

// ToString convert i to string.
//
//nolint:revive,cyclop
func ToString(i any) (string, error) {
	switch s := i.(type) {
	case string:
		return s, nil
	case bool:
		return strconv.FormatBool(s), nil
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64), nil
	case int:
		return strconv.Itoa(s), nil
	case int64:
		return strconv.FormatInt(s, 10), nil
	case int32:
		return strconv.Itoa(int(s)), nil
	case int16:
		return strconv.FormatInt(int64(s), 10), nil
	case int8:
		return strconv.FormatInt(int64(s), 10), nil
	case uint:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint64:
		return strconv.FormatUint(s, 10), nil
	case uint32:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(s), 10), nil
	case time.Duration:
		return s.String(), nil
	}

	return "", xerrors.Errorf("Unsupport target type '%T'", i)
}
