package cli

import (
	"flag"
	"fmt"
	"strings"
)

// StringSlice is an opaque type for []string to satisfy flag.Value and flag.Getter
type StringSlice []string

// Set appends the string value to the list of values
func (f *StringSlice) Set(value string) error {
	*f = append(*f, value)
	return nil
}

// String returns a readable representation of this value (for usage defaults)
func (f *StringSlice) String() string {
	return fmt.Sprintf("%s", *f)
}

// Value returns the slice of strings set by this flag
func (f *StringSlice) Value() []string {
	return *f
}

// Get returns the slice of strings set by this flag
func (f *StringSlice) Get() interface{} {
	return *f
}

// StringSliceFlag is a flag with type *StringSlice
type StringSliceFlag struct {
	Name      string
	Usage     string
	EnvVar    string
	FilePath  string
	Required  bool
	Hidden    bool
	TakesFile bool
	Value     *StringSlice
}

// String returns a readable representation of this value
// (for usage defaults)
func (f StringSliceFlag) String() string {
	return FlagStringer(f)
}

// GetName returns the name of the flag
func (f StringSliceFlag) GetName() string {
	return f.Name
}

// IsRequired returns whether or not the flag is required
func (f StringSliceFlag) IsRequired() bool {
	return f.Required
}

// TakesValue returns true of the flag takes a value, otherwise false
func (f StringSliceFlag) TakesValue() bool {
	return true
}

// GetUsage returns the usage string for the flag
func (f StringSliceFlag) GetUsage() string {
	return f.Usage
}

// GetValue returns the flags value as string representation and an empty
// string if the flag takes no value at all.
func (f StringSliceFlag) GetValue() string {
	if f.Value != nil {
		return f.Value.String()
	}
	return ""
}

// Apply populates the flag given the flag set and environment
// Ignores errors
func (f StringSliceFlag) Apply(set *flag.FlagSet) {
	_ = f.ApplyWithError(set)
}

// ApplyWithError populates the flag given the flag set and environment
func (f StringSliceFlag) ApplyWithError(set *flag.FlagSet) error {
	if envVal, ok := flagFromFileEnv(f.FilePath, f.EnvVar); ok {
		newVal := &StringSlice{}
		for _, s := range strings.Split(envVal, ",") {
			s = strings.TrimSpace(s)
			if err := newVal.Set(s); err != nil {
				return fmt.Errorf("could not parse %s as string value for flag %s: %s", envVal, f.Name, err)
			}
		}
		if f.Value == nil {
			f.Value = newVal
		} else {
			*f.Value = *newVal
		}
	}

	eachName(f.Name, func(name string) {
		if f.Value == nil {
			f.Value = &StringSlice{}
		}
		set.Var(f.Value, name, f.Usage)
	})

	return nil
}

// StringSlice looks up the value of a local StringSliceFlag, returns
// nil if not found
func (c *Context) StringSlice(name string) []string {
	return lookupStringSlice(name, c.flagSet)
}

// GlobalStringSlice looks up the value of a global StringSliceFlag, returns
// nil if not found
func (c *Context) GlobalStringSlice(name string) []string {
	if fs := lookupGlobalFlagSet(name, c); fs != nil {
		return lookupStringSlice(name, fs)
	}
	return nil
}

func lookupStringSlice(name string, set *flag.FlagSet) []string {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := (f.Value.(*StringSlice)).Value(), error(nil)
		if err != nil {
			return nil
		}
		return parsed
	}
	return nil
}
