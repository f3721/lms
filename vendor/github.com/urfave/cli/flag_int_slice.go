package cli

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// IntSlice is an opaque type for []int to satisfy flag.Value and flag.Getter
type IntSlice []int

// Set parses the value into an integer and appends it to the list of values
func (f *IntSlice) Set(value string) error {
	tmp, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*f = append(*f, tmp)
	return nil
}

// String returns a readable representation of this value (for usage defaults)
func (f *IntSlice) String() string {
	return fmt.Sprintf("%#v", *f)
}

// Value returns the slice of ints set by this flag
func (f *IntSlice) Value() []int {
	return *f
}

// Get returns the slice of ints set by this flag
func (f *IntSlice) Get() interface{} {
	return *f
}

// IntSliceFlag is a flag with type *IntSlice
type IntSliceFlag struct {
	Name     string
	Usage    string
	EnvVar   string
	FilePath string
	Required bool
	Hidden   bool
	Value    *IntSlice
}

// String returns a readable representation of this value
// (for usage defaults)
func (f IntSliceFlag) String() string {
	return FlagStringer(f)
}

// GetName returns the name of the flag
func (f IntSliceFlag) GetName() string {
	return f.Name
}

// IsRequired returns whether or not the flag is required
func (f IntSliceFlag) IsRequired() bool {
	return f.Required
}

// TakesValue returns true of the flag takes a value, otherwise false
func (f IntSliceFlag) TakesValue() bool {
	return true
}

// GetUsage returns the usage string for the flag
func (f IntSliceFlag) GetUsage() string {
	return f.Usage
}

// GetValue returns the flags value as string representation and an empty
// string if the flag takes no value at all.
func (f IntSliceFlag) GetValue() string {
	if f.Value != nil {
		return f.Value.String()
	}
	return ""
}

// Apply populates the flag given the flag set and environment
// Ignores errors
func (f IntSliceFlag) Apply(set *flag.FlagSet) {
	_ = f.ApplyWithError(set)
}

// ApplyWithError populates the flag given the flag set and environment
func (f IntSliceFlag) ApplyWithError(set *flag.FlagSet) error {
	if envVal, ok := flagFromFileEnv(f.FilePath, f.EnvVar); ok {
		newVal := &IntSlice{}
		for _, s := range strings.Split(envVal, ",") {
			s = strings.TrimSpace(s)
			if err := newVal.Set(s); err != nil {
				return fmt.Errorf("could not parse %s as int slice value for flag %s: %s", envVal, f.Name, err)
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
			f.Value = &IntSlice{}
		}
		set.Var(f.Value, name, f.Usage)
	})

	return nil
}

// IntSlice looks up the value of a local IntSliceFlag, returns
// nil if not found
func (c *Context) IntSlice(name string) []int {
	return lookupIntSlice(name, c.flagSet)
}

// GlobalIntSlice looks up the value of a global IntSliceFlag, returns
// nil if not found
func (c *Context) GlobalIntSlice(name string) []int {
	if fs := lookupGlobalFlagSet(name, c); fs != nil {
		return lookupIntSlice(name, fs)
	}
	return nil
}

func lookupIntSlice(name string, set *flag.FlagSet) []int {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := (f.Value.(*IntSlice)).Value(), error(nil)
		if err != nil {
			return nil
		}
		return parsed
	}
	return nil
}
