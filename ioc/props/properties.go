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

package props

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/anyvoxel/airmid/anvil/conv"
	"github.com/anyvoxel/airmid/anvil/xerrors"
	"github.com/anyvoxel/airmid/anvil/xreflect"
	slogctx "github.com/veqryn/slog-context"
)

// NewProperties return the Properties impl.
func NewProperties() Properties {
	return propertiesImpl(make(map[string]string))
}

type propertiesImpl map[string]string

//nolint:revive,cyclop
func (p propertiesImpl) Get(ctx context.Context, key string, opts ...GetOption) (ret any, err error) {
	opt := defaultGetOption()
	for _, o := range opts {
		o.Apply(opt)
	}

	if err = opt.Complete(); err != nil {
		return nil, err
	}

	var targetValue reflect.Value
	targetValue, err = xreflect.IndirectToSetableValue(opt.Target)
	if err != nil {
		return nil, err
	}

	var propValues []string
	switch {
	case targetValue.Kind() == reflect.Slice:
		vstrs, err := p.doGetSlice(key)
		if err != nil {
			if !xerrors.Is(err, xerrors.ErrNotFound) || opt.Default == nil {
				return nil, err
			}

			if *opt.Default == "" {
				vstrs = []string{}
			} else {
				vstrs = strings.Split(*opt.Default, ",")
			}
		}
		propValues = vstrs
	default:
		vstr, err := p.doGet(key)
		if err != nil {
			if !xerrors.Is(err, xerrors.ErrNotFound) || opt.Default == nil {
				return nil, err
			}

			vstr = *opt.Default
		}
		propValues = []string{vstr}
	}

	cv, err := conv.ConvertTo(context.Background(), targetValue.Type(), propValues)
	if err != nil {
		return nil, err
	}
	targetValue.Set(reflect.ValueOf(cv).Convert(targetValue.Type()))
	return opt.Target.Interface(), nil
}

func (p propertiesImpl) doGet(key string) (string, error) {
	val, ok := p[key]
	if ok {
		return val, nil
	}

	return "", xerrors.WrapNotFound("property with key='%v' not found", key)
}

type indexString struct {
	i int64
	v string
}

func (p propertiesImpl) doGetSlice(key string) ([]string, error) {
	val, ok := p[key]
	if ok {
		return []string{val}, nil
	}

	ret := []indexString{}
	for k, v := range p {
		if !strings.HasPrefix(k, key) {
			continue
		}

		k = k[len(key):]
		if !strings.HasPrefix(k, "[") || !strings.HasSuffix(k, "]") {
			continue
		}

		k = k[1 : len(k)-1]
		i, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			return nil, err
		}

		ret = append(ret, indexString{
			i: i,
			v: v,
		})
	}

	sort.Slice(ret, func(i int, j int) bool {
		return ret[i].i < ret[j].i
	})

	if len(ret) > 0 {
		rr := make([]string, 0, len(ret))
		for _, v := range ret {
			rr = append(rr, v.v)
		}

		return rr, nil
	}

	return nil, xerrors.WrapNotFound("property slice with key='%v' not found", key)
}

//nolint:revive,exhaustive
func (p propertiesImpl) Set(ctx context.Context, key string, val any) error {
	switch v := reflect.ValueOf(val); v.Kind() {
	case reflect.Map:
		// If the val is a map, we expand the val with keys and set it recursive
		for _, k := range v.MapKeys() {
			kstr, err := conv.ToString(k.Interface())
			if err != nil {
				return xerrors.Wrapf(err, "Cannot convert map's key '%v' to string", k)
			}

			kstr = fmt.Sprintf("%s.%s", key, kstr)
			kvalue := v.MapIndex(k).Interface()
			err = p.Set(ctx, kstr, kvalue)
			if err != nil {
				return xerrors.Wrapf(err, "Cannot set val for map's key '%v'", kstr)
			}
		}
	case reflect.Array, reflect.Slice:
		// If the val is a array/slice, we expand the val with index and set it recursive
		for i := 0; i < v.Len(); i++ {
			kstr := fmt.Sprintf("%s[%d]", key, i)
			kvalue := v.Index(i).Interface()
			err := p.Set(ctx, kstr, kvalue)
			if err != nil {
				return xerrors.Wrapf(err, "Cannot set val for array/slice index's key '%v'", kstr)
			}
		}
	default:
		value, err := conv.ToString(val)
		if err != nil {
			return xerrors.Wrapf(err, "Cannot convert value for key '%s' to string", key)
		}
		p[key] = value
		slogctx.FromCtx(ctx).DebugContext(
			ctx,
			"set property success",
			slog.String("Key", key),
			slog.String("Value", value),
		)
	}

	return nil
}
