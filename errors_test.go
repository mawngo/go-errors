// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

package errors

import (
	//lint:ignore faillint Custom errors package tests need to import standard library errors.
	stderrors "errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"
)

const msg = "test_error_message"
const wrapper = "test_wrapper"

var ErrTest = Raw("global_defined_error")

func TestNewf(t *testing.T) {
	err := Newf(msg)
	if err.Error() != msg {
		t.Fatalf("the root error message must match")
	}

	reg := regexp.MustCompile(msg + `[ \n]+> github\.com\/mawngo\/go-errors\.TestNewf	.*\/go-errors\/errors_test\.go:\d+`)
	if !reg.MatchString(fmt.Sprintf("%+v", err)) {
		t.Fatalf("matching stacktrace in errors.New")
	}
}

func TestNewfFormatted(t *testing.T) {
	fmtMsg := msg + " key=%v"
	expectedMsg := msg + " key=value"

	err := Newf(fmtMsg, "value")
	if err.Error() != expectedMsg {
		t.Fatalf("the root error message must match")
	}
	reg := regexp.MustCompile(expectedMsg + `[ \n]+> github\.com\/mawngo\/go-errors\.TestNewfFormatted	.*\/go-errors\/errors_test\.go:\d+`)
	if !reg.MatchString(fmt.Sprintf("%+v", err)) {
		t.Fatalf("matching stacktrace in errors.New with format string")
	}
}

func TestWrapf(t *testing.T) {
	err := Newf(msg)
	err = Wrapf(err, wrapper)

	expectedMsg := wrapper + ": " + msg
	if err.Error() != expectedMsg {
		t.Fatalf("the root error message must match")
	}

	reg := regexp.MustCompile(`test_wrapper[ \n]+> github\.com\/mawngo\/go-errors\.TestWrapf	.*\/go-errors\/errors_test\.go:\d+
[[:ascii:]]+test_error_message[ \n]+> github\.com\/mawngo\/go-errors\.TestWrapf	.*\/go-errors\/errors_test\.go:\d+`)

	errMsg := fmt.Sprintf("%+v", err)
	if !reg.MatchString(errMsg) {
		t.Fatalf("matching stacktrace in errors.Wrapf")
	}
}

func TestWrap(t *testing.T) {
	err := Newf(msg)
	err = Wrap(err)

	expectedMsg := msg
	if err.Error() != expectedMsg {
		t.Fatalf("the root error message must match")
	}

	reg := regexp.MustCompile(`test_error_message[ \n]+> github\.com\/mawngo\/go-errors\.TestWrap	.*\/go-errors\/errors_test\.go:\d+
[[:ascii:]]+test_error_message[ \n]+> github\.com\/mawngo\/go-errors\.TestWrap	.*\/go-errors\/errors_test\.go:\d+`)

	errMsg := fmt.Sprintf("%+v", err)
	if !reg.MatchString(errMsg) {
		t.Fatalf("matching stacktrace in errors.Wrap")
	}

}

func TestUnwrap(t *testing.T) {
	// test with base error
	err := Newf(msg)

	for i, tc := range []struct {
		err      error
		expected string
		isNil    bool
	}{
		{
			// no wrapping
			err:   err,
			isNil: true,
		},
		{
			err:      Wrapf(err, wrapper),
			expected: "test_error_message",
		},
		{
			err:      Wrapf(Wrapf(err, wrapper), wrapper),
			expected: "test_wrapper: test_error_message",
		},
		// check primitives errors
		{
			err:   stderrors.New("std-error"),
			isNil: true,
		},
		{
			err:      Wrapf(stderrors.New("std-error"), wrapper),
			expected: "std-error",
		},
		{
			err:   nil,
			isNil: true,
		},
	} {
		t.Run("TestCase"+strconv.Itoa(i), func(t *testing.T) {
			unwrapped := Unwrap(tc.err)
			if tc.isNil {
				if unwrapped != nil {
					t.Fatalf("expected nil, got %v", unwrapped)
				}
				return
			}
			if tc.expected != unwrapped.Error() {
				t.Fatalf("Unwrap must match expected output")
			}
		})
	}
}

func TestCause(t *testing.T) {
	// test with base error that implements interface containing Unwrap method
	err := Newf(msg)

	for i, tc := range []struct {
		err      error
		expected string
		isNil    bool
	}{
		{
			// no wrapping
			err:   err,
			isNil: true,
		},
		{
			err:   Wrapf(err, wrapper),
			isNil: true,
		},
		{
			err:   Wrapf(Wrapf(err, wrapper), wrapper),
			isNil: true,
		},
		// check primitives errors
		{
			err:      stderrors.New("std-error"),
			expected: "std-error",
		},
		{
			err:      Wrapf(stderrors.New("std-error"), wrapper),
			expected: "std-error",
		},
		{
			err:   nil,
			isNil: true,
		},
	} {
		t.Run("TestCase"+strconv.Itoa(i), func(t *testing.T) {
			cause := Cause(tc.err)
			if tc.isNil {
				if cause != nil {
					t.Fatalf("expected nil, got %v", cause)
				}
				return
			}
			if tc.expected != cause.Error() {
				t.Fatalf("Cause must match expected output")
			}
		})
	}
}

func TestErrorIs(t *testing.T) {
	// test with base error that implements interface containing Unwrap method
	err := Wrap(ErrTest)
	if !stderrors.Is(err, ErrTest) {
		t.Fatalf("expected error to be equal to ErrTest")
	}
}

func TestErrorAs(t *testing.T) {
	// test with base error that implements interface containing Unwrap method
	err := Wrap(ErrTest)
	var e *base
	if !stderrors.As(err, &e) {
		t.Fatalf("expected error to be assignable to base error")
	}
}
