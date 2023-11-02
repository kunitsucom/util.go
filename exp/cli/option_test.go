package cliz

type testOption struct {
	Name        string
	Short       string
	Environment string
	Description string
	Default     *string
	value       *string
}

func (o *testOption) GetName() string         { return o.Name }
func (o *testOption) GetShort() string        { return o.Short }
func (o *testOption) GetEnvironment() string  { return o.Environment }
func (o *testOption) HasDefault() bool        { return o.Default != nil }
func (o *testOption) getDefault() interface{} { return o.Default }
func (o *testOption) HasValue() bool          { return o.value != nil }
func (o *testOption) GetDescription() string  { return o.Description }
func (o *testOption) private()                {}
