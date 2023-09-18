//nolint:dupword
package stringz_test

import (
	"fmt"

	stringz "github.com/kunitsucom/util.go/strings"
)

func ExampleIsSnakeCase() {
	fmt.Println(stringz.IsSnakeCase("snake_case"))
	fmt.Println(stringz.IsSnakeCase("kebab-case"))
	fmt.Println(stringz.IsSnakeCase("camelCase"))
	fmt.Println(stringz.IsSnakeCase("PascalCase"))
	fmt.Println(stringz.IsSnakeCase("go"))
	fmt.Println(stringz.IsSnakeCase("type_script"))
	fmt.Println(stringz.IsSnakeCase("postgre_sql"))
	fmt.Println(stringz.IsSnakeCase("goV1_21"))
	fmt.Println(stringz.IsSnakeCase("MySQLV8"))
	// Output:
	// true
	// false
	// false
	// false
	// false
	// true
	// true
	// false
	// false
}

func ExampleIsKebabCase() {
	fmt.Println(stringz.IsKebabCase("snake_case"))
	fmt.Println(stringz.IsKebabCase("kebab-case"))
	fmt.Println(stringz.IsKebabCase("camelCase"))
	fmt.Println(stringz.IsKebabCase("PascalCase"))
	fmt.Println(stringz.IsKebabCase("go"))
	fmt.Println(stringz.IsKebabCase("type-script"))
	fmt.Println(stringz.IsKebabCase("postgre-sql"))
	fmt.Println(stringz.IsKebabCase("goV1_21"))
	fmt.Println(stringz.IsKebabCase("MySQLV8"))
	// Output:
	// false
	// true
	// false
	// false
	// false
	// true
	// true
	// false
	// false
}

func ExampleIsCamelCase() {
	fmt.Println(stringz.IsCamelCase("snake_case"))
	fmt.Println(stringz.IsCamelCase("kebab-case"))
	fmt.Println(stringz.IsCamelCase("camelCase"))
	fmt.Println(stringz.IsCamelCase("PascalCase"))
	fmt.Println(stringz.IsCamelCase("go"))
	fmt.Println(stringz.IsCamelCase("typeScript"))
	fmt.Println(stringz.IsCamelCase("postgreSQL"))
	fmt.Println(stringz.IsCamelCase("goV1_21"))
	fmt.Println(stringz.IsCamelCase("MySQLV8"))
	// Output:
	// false
	// false
	// true
	// false
	// true
	// true
	// true
	// true
	// false
}

func ExampleIsPascalCase() {
	fmt.Println(stringz.IsPascalCase("snake_case"))
	fmt.Println(stringz.IsPascalCase("kebab-case"))
	fmt.Println(stringz.IsPascalCase("camelCase"))
	fmt.Println(stringz.IsPascalCase("PascalCase"))
	fmt.Println(stringz.IsPascalCase("Go"))
	fmt.Println(stringz.IsPascalCase("TypeScript"))
	fmt.Println(stringz.IsPascalCase("PostgreSQL"))
	fmt.Println(stringz.IsPascalCase("goV1_21"))
	fmt.Println(stringz.IsPascalCase("MySQLV8"))
	// Output:
	// false
	// false
	// false
	// true
	// true
	// true
	// true
	// false
	// true
}

func ExampleSplitSnakeCase() {
	fmt.Println(stringz.SplitSnakeCase("snake_case"))
	fmt.Println(stringz.SplitSnakeCase("kebab-case"))
	fmt.Println(stringz.SplitSnakeCase("camelCase"))
	fmt.Println(stringz.SplitSnakeCase("PascalCase"))
	fmt.Println(stringz.SplitSnakeCase("go"))
	fmt.Println(stringz.SplitSnakeCase("type_script"))
	fmt.Println(stringz.SplitSnakeCase("postgre_sql"))
	fmt.Println(stringz.SplitSnakeCase("goV1_21"))
	fmt.Println(stringz.SplitSnakeCase("MySQLV8"))
	// Output:
	// [snake case]
	// [kebab-case]
	// [camelCase]
	// [PascalCase]
	// [go]
	// [type script]
	// [postgre sql]
	// [goV1 21]
	// [MySQLV8]
}

func ExampleSplitKebabCase() {
	fmt.Println(stringz.SplitKebabCase("snake_case"))
	fmt.Println(stringz.SplitKebabCase("kebab-case"))
	fmt.Println(stringz.SplitKebabCase("camelCase"))
	fmt.Println(stringz.SplitKebabCase("PascalCase"))
	fmt.Println(stringz.SplitKebabCase("go"))
	fmt.Println(stringz.SplitKebabCase("type-script"))
	fmt.Println(stringz.SplitKebabCase("postgre-sql"))
	fmt.Println(stringz.SplitKebabCase("goV1_21"))
	fmt.Println(stringz.SplitKebabCase("MySQLV8"))
	// Output:
	// [snake_case]
	// [kebab case]
	// [camelCase]
	// [PascalCase]
	// [go]
	// [type script]
	// [postgre sql]
	// [goV1_21]
	// [MySQLV8]
}

func ExampleSplitCamelCase() {
	fmt.Println(stringz.SplitCamelCase("snake_case"))
	fmt.Println(stringz.SplitCamelCase("kebab-case"))
	fmt.Println(stringz.SplitCamelCase("camelCase"))
	fmt.Println(stringz.SplitCamelCase("PascalCase"))
	fmt.Println(stringz.SplitCamelCase("go"))
	fmt.Println(stringz.SplitCamelCase("typeScript"))
	fmt.Println(stringz.SplitCamelCase("postgreSQL"))
	fmt.Println(stringz.SplitCamelCase("goV1_21"))
	fmt.Println(stringz.SplitCamelCase("MySQLV8"))
	// Output:
	// [snake_case]
	// [kebab-case]
	// [camel Case]
	// [Pascal Case]
	// [go]
	// [type Script]
	// [postgre SQL]
	// [go V1_21]
	// [My SQL V8]
}

func ExampleSplitPascalCase() {
	fmt.Println(stringz.SplitPascalCase("snake_case"))
	fmt.Println(stringz.SplitPascalCase("kebab-case"))
	fmt.Println(stringz.SplitPascalCase("camelCase"))
	fmt.Println(stringz.SplitPascalCase("PascalCase"))
	fmt.Println(stringz.SplitPascalCase("Go"))
	fmt.Println(stringz.SplitPascalCase("TypeScript"))
	fmt.Println(stringz.SplitPascalCase("PostgreSQL"))
	fmt.Println(stringz.SplitPascalCase("goV1_21"))
	fmt.Println(stringz.SplitPascalCase("MySQLV8"))
	// Output:
	// [snake_case]
	// [kebab-case]
	// [camel Case]
	// [Pascal Case]
	// [Go]
	// [Type Script]
	// [Postgre SQL]
	// [go V1_21]
	// [My SQL V8]
}

func ExampleSplitCase() {
	fmt.Println(stringz.SplitCase("snake_case"))
	fmt.Println(stringz.SplitCase("kebab-case"))
	fmt.Println(stringz.SplitCase("camelCase"))
	fmt.Println(stringz.SplitCase("PascalCase"))
	fmt.Println(stringz.SplitCase("go"))
	fmt.Println(stringz.SplitCase("typeScript"))
	fmt.Println(stringz.SplitCase("postgreSQL"))
	fmt.Println(stringz.SplitCase("goV1_21"))
	fmt.Println(stringz.SplitCase("MySQLV8"))
	fmt.Println(stringz.SplitCase("A - 3"))
	// Output:
	// [snake case]
	// [kebab case]
	// [camel Case]
	// [Pascal Case]
	// [go]
	// [type Script]
	// [postgre SQL]
	// [go V1_21]
	// [My SQL V8]
	// [A - 3]
}
