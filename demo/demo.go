package demo

import "errors"

type AnyValue struct {
	Val any
	Err error
}

func (a AnyValue) Bytes() ([]byte, error) {
	if a.Err != nil {
		return nil, a.Err
	}
	str, ok := a.Val.([]byte)
	if !ok {
		return nil, errors.New("无法转换的类型")
	}
	return str, nil
}

func (a AnyValue) Int64() (int64, error) {
	if a.Err != nil {
		return 0, a.Err
	}
	str, ok := a.Val.(int64)
	if !ok {
		return 0, errors.New("无法转换的类型")
	}
	return str, nil
}

