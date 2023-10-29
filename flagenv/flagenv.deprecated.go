// Deprecated: *flag.FlagSet is difficult to use and hard to test.
package flagenv

import (
	"flag"
	"os"
	"time"

	"github.com/kunitsucom/util.go/env"
)

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
//
//nolint:revive
type FlagEnvSet struct {
	*flag.FlagSet
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func NewFlagEnvSet(name string, errorHandling flag.ErrorHandling) *FlagEnvSet {
	return &FlagEnvSet{flag.NewFlagSet(name, errorHandling)}
}

// CommandLine is the default set of command-line flags, parsed from os.Args.
// The top-level functions such as BoolVar, Arg, and so on are wrappers for the
// methods of CommandLine.
//
// Deprecated: *flag.FlagSet is difficult to use and hard to test.
//
//nolint:gochecknoglobals
var CommandLine = NewFlagEnvSet(os.Args[0], flag.ExitOnError)

//
// Bool.
//

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func (s *FlagEnvSet) BoolVar(p *bool, name string, environ string, value bool, usage string) {
	s.FlagSet.BoolVar(p, name, env.BoolOrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func BoolVar(p *bool, name string, environ string, value bool, usage string) {
	CommandLine.FlagSet.BoolVar(p, name, env.BoolOrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func (s *FlagEnvSet) Bool(name string, environ string, value bool, usage string) *bool {
	return s.FlagSet.Bool(name, env.BoolOrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func Bool(name string, environ string, value bool, usage string) *bool {
	return CommandLine.FlagSet.Bool(name, env.BoolOrDefault(environ, value), usage)
}

//
// Int
//

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func (s *FlagEnvSet) IntVar(p *int, name string, environ string, value int, usage string) {
	s.FlagSet.IntVar(p, name, env.IntOrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func IntVar(p *int, name string, environ string, value int, usage string) {
	CommandLine.FlagSet.IntVar(p, name, env.IntOrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func (s *FlagEnvSet) Int(name string, environ string, value int, usage string) *int {
	return s.FlagSet.Int(name, env.IntOrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func Int(name string, environ string, value int, usage string) *int {
	return CommandLine.FlagSet.Int(name, env.IntOrDefault(environ, value), usage)
}

//
// Int64
//

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func (s *FlagEnvSet) Int64Var(p *int64, name string, environ string, value int64, usage string) {
	s.FlagSet.Int64Var(p, name, env.Int64OrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func Int64Var(p *int64, name string, environ string, value int64, usage string) {
	CommandLine.FlagSet.Int64Var(p, name, env.Int64OrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func (s *FlagEnvSet) Int64(name string, environ string, value int64, usage string) *int64 {
	return s.FlagSet.Int64(name, env.Int64OrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func Int64(name string, environ string, value int64, usage string) *int64 {
	return CommandLine.FlagSet.Int64(name, env.Int64OrDefault(environ, value), usage)
}

//
// Uint
//

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func (s *FlagEnvSet) UintVar(p *uint, name string, environ string, value uint, usage string) {
	s.FlagSet.UintVar(p, name, env.UintOrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func UintVar(p *uint, name string, environ string, value uint, usage string) {
	CommandLine.FlagSet.UintVar(p, name, env.UintOrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func (s *FlagEnvSet) Uint(name string, environ string, value uint, usage string) *uint {
	return s.FlagSet.Uint(name, env.UintOrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func Uint(name string, environ string, value uint, usage string) *uint {
	return CommandLine.FlagSet.Uint(name, env.UintOrDefault(environ, value), usage)
}

//
// Uint64
//

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func (s *FlagEnvSet) Uint64Var(p *uint64, name string, environ string, value uint64, usage string) {
	s.FlagSet.Uint64Var(p, name, env.Uint64OrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func Uint64Var(p *uint64, name string, environ string, value uint64, usage string) {
	CommandLine.FlagSet.Uint64Var(p, name, env.Uint64OrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func (s *FlagEnvSet) Uint64(name string, environ string, value uint64, usage string) *uint64 {
	return s.FlagSet.Uint64(name, env.Uint64OrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func Uint64(name string, environ string, value uint64, usage string) *uint64 {
	return CommandLine.FlagSet.Uint64(name, env.Uint64OrDefault(environ, value), usage)
}

//
// String
//

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func (s *FlagEnvSet) StringVar(p *string, name string, environ string, value string, usage string) {
	s.FlagSet.StringVar(p, name, env.StringOrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func StringVar(p *string, name string, environ string, value string, usage string) {
	CommandLine.FlagSet.StringVar(p, name, env.StringOrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func (s *FlagEnvSet) String(name string, environ string, value string, usage string) *string {
	return s.FlagSet.String(name, env.StringOrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func String(name string, environ string, value string, usage string) *string {
	return CommandLine.FlagSet.String(name, env.StringOrDefault(environ, value), usage)
}

//
// Float64
//

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func (s *FlagEnvSet) Float64Var(p *float64, name string, environ string, value float64, usage string) {
	s.FlagSet.Float64Var(p, name, env.Float64OrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func Float64Var(p *float64, name string, environ string, value float64, usage string) {
	CommandLine.FlagSet.Float64Var(p, name, env.Float64OrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func (s *FlagEnvSet) Float64(name string, environ string, value float64, usage string) *float64 {
	return s.FlagSet.Float64(name, env.Float64OrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func Float64(name string, environ string, value float64, usage string) *float64 {
	return CommandLine.FlagSet.Float64(name, env.Float64OrDefault(environ, value), usage)
}

//
// Duration (Second)
//

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func (s *FlagEnvSet) SecondVar(p *time.Duration, name string, environ string, value time.Duration, usage string) {
	s.FlagSet.DurationVar(p, name, env.SecondOrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func SecondVar(p *time.Duration, name string, environ string, value time.Duration, usage string) {
	CommandLine.FlagSet.DurationVar(p, name, env.SecondOrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func (s *FlagEnvSet) Second(name string, environ string, value time.Duration, usage string) *time.Duration {
	return s.FlagSet.Duration(name, env.SecondOrDefault(environ, value), usage)
}

// Deprecated: *flag.FlagSet is difficult to use and hard to test.
func Second(name string, environ string, value time.Duration, usage string) *time.Duration {
	return CommandLine.FlagSet.Duration(name, env.SecondOrDefault(environ, value), usage)
}
