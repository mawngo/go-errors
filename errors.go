// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

//nolint
// The idea of writing errors package in thanos is highly motivated from the Tast project of Chromium OS Authors. However, instead of
// copying the package, we end up writing our own simplified logic borrowing some ideas from the errors and github.com/pkg/errors.
// A big thanks to all of them.

// Package errors provides basic utilities to manipulate errors with a useful stacktrace. It combines the
// benefits of errors.New and fmt.Errorf world into a single package.
package errors

import (
	//lint:ignore faillint Custom errors package needs to import standard library errors.
	"errors"
	"fmt"
	"strings"
)

// base is the fundamental struct that implements the error interface and the acts as the backbone of this errors package.
type base struct {
	// info contains the error message passed through calls like errors.Wrap, errors.New.
	info string
	// stacktrace stores information about the program counters - i.e. where this error was generated.
	stack stacktrace
	// err is the actual error which is being wrapped with a stacktrace and message information.
	err error
}

// Error implements the error interface.
func (b *base) Error() string {
	if b.err != nil {
		e := b.err.Error()
		if e == b.info {
			return e
		}
		return b.info + ": " + e
	}
	return b.info
}

// String implements the [fmt.Stringer] interface.
// String returns error message of this error only.
// For full error chain message use Error instead.
func (b *base) String() string {
	return b.info
}

// Unwrap implements the error Unwrap interface.
func (b *base) Unwrap() error {
	return b.err
}

// Format implements the [fmt.Formatter] interface to support the formatting of an error chain with "%+v" verb.
// Whenever error is printed with %+v format verb, stacktrace info gets dumped to the output.
func (b *base) Format(s fmt.State, verb rune) {
	if verb == 'v' && s.Flag('+') {
		_, _ = s.Write([]byte(formatErrorChain(b)))
		return
	}
	_, _ = s.Write([]byte(b.Error()))
}

// Newf formats according to a format specifier and returns a new error with a stacktrace
// with recent call frames.
// Each call to Newf returns a distinct error value even if the text is
// identical. An alternative of the errors.New function.
func Newf(format string, args ...any) error {
	info := format
	if len(args) > 0 {
		info = fmt.Sprintf(format, args...)
	}
	return &base{
		info:  info,
		stack: newStackTrace(),
		err:   nil,
	}
}

// New create a new error with a stacktrace with recent call frames.
// Each call to New returns a distinct error value even if the text is identical.
// Deprecated: use [Newf] for error with stacktrace, use [Raw] for error without stacktrace.
func New(message string) error {
	return &base{
		info:  message,
		stack: newStackTrace(),
		err:   nil,
	}
}

// Wrapf returns a new error by formatting the error message with the supplied format specifier
// and wrapping another error with a stacktrace containing recent call frames.
//
// If cause is nil, this method returns nil.
func Wrapf(cause error, format string, args ...any) error {
	if cause == nil {
		return nil
	}
	info := format
	if len(args) > 0 {
		info = fmt.Sprintf(format, args...)
	}
	return &base{
		info:  info,
		stack: newStackTrace(),
		err:   cause,
	}
}

// Wrap returns a new error by wrapping another error with a stacktrace containing recent call frames.
// If cause is nil, this method returns nil.
//
// If you want to add context msg to the error, use [Wrapf].
func Wrap(cause error) error {
	if cause == nil {
		return nil
	}
	return &base{
		info:  cause.Error(),
		stack: newStackTrace(),
		err:   cause,
	}
}

// Cause returns the result of repeatedly calling the Unwrap method on err, if err's
// type implements an Unwrap method. Otherwise, Cause returns the last encountered error.
// The difference between Unwrap and Cause is the first one performs unwrapping of one level
// but Cause returns the last err (whether it's nil or not) where it failed to assert
// the interface containing the Unwrap method.
// This is a replacement of errors.Cause without the causer interface from pkg/errors which
// actually can be sufficed through the errors.Is function. But considering some use cases
// where we need to peel off all the external layers applied through errors.Wrap family,
// it is useful (where external SDK doesn't use errors.Is internally).
func Cause(err error) error {
	for err != nil {
		e, ok := err.(interface {
			Unwrap() error
		})
		if !ok {
			return err
		}
		err = e.Unwrap()
	}
	return nil
}

// formatErrorChain formats an error chain.
func formatErrorChain(err error) string {
	var buf strings.Builder
	for err != nil {
		var e *base
		if errors.As(err, &e) {
			buf.WriteString(e.info)
			buf.WriteString("\n")
			buf.WriteString(fmt.Sprintf("%v", e.stack))
			err = e.err
		} else {
			buf.WriteString(err.Error())
			buf.WriteString("\n")
			err = nil
		}
	}
	return buf.String()
}

// The functions `Is`, `As` & `Unwrap` provides a thin wrapper around the builtin errors
// package in go. Just for the sake of completeness and correct autocompletion behaviors from
// IDEs they have been wrapped using functions instead of using variable to reference them
// as first class functions (eg: var Is = errros.Is ).

// Is is a wrapper of built-in errors.Is. It reports whether any error in err's
// chain matches target. The chain consists of err itself followed by the sequence
// of errors obtained by repeatedly calling Unwrap.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As is a wrapper of built-in [errors.As]. It finds the first error in err's
// chain that matches target, and if one is found, sets target to that error
// value and returns true. Otherwise, it returns false.
func As(err error, target any) bool {
	return errors.As(err, target)
}

// Unwrap is a wrapper of built-in errors.Unwrap.
// Unwrap returns the result of calling the Unwrap method on err, if err's type contains an Unwrap method
// returning error. Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Join is a wrapper of built-in [errors.Join]
// Join returns an error that wraps the given errors.
// Any nil error values are discarded.
// Join returns nil if every value in errs is nil.
// The error formats as the concatenation of the strings obtained
// by calling the Error method of each element of errs, with a newline
// between each string.
//
// A non-nil error returned by Join implements the Unwrap() []error method.
func Join(errs ...error) error {
	return errors.Join(errs...)
}

// Raw is a wrapper of built-in [errors.New].
// Raw create an error without stacktrace,
// for defining error constant without having to import the go standard errors package.
// Use [Newf] if you want to return an error with a stacktrace.
func Raw(msg string) error {
	return errors.New(msg)
}

// ErrUnsupported is a wrapper of built-in [errors.ErrUnsupported]
// [errors.ErrUnsupported] indicates that a requested operation cannot be performed,
// because it is unsupported. For example, a call to [os.Link] when using a
// file system that does not support hard links.
//
// Functions and methods should not return this error but should instead
// return an error including appropriate context that satisfies
//
//	errors.Is(err, errors.ErrUnsupported)
//
// either by directly wrapping ErrUnsupported or by implementing an [Is] method.
//
// Functions and methods should document the cases in which an error
// wrapping this will be returned.
var ErrUnsupported = errors.ErrUnsupported
