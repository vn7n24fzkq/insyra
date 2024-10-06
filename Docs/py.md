# `py` Package Documentation

## Overview

The `py` package allows Golang programs to execute Python code seamlessly and interactively. It provides functionality to pass Go variables into Python scripts, and execute the Python code. Results from the Python script can be sent back to the Go program automatically.

## Functions

### `RunCode`

```go
func RunCode(code string) map[string]interface{}
```

This function is used to execute arbitrary Python code. It appends default Python code required for communication and runs the code.

#### Parameters

- `code` (string): The Python code to be executed.

#### Returns

- `map[string]interface{}`: A map representing the result received from the Python server. This map will contain the data returned from Python through the `insyra_return` function.

#### Example

```go
result := RunCode(`
print("Hello from Python")
insyra_return({"message": "Hello from Python", "value": 123})
`)
fmt.Println(result)
```

---

### `RunCodef`

```go
func RunCodef(codeTemplate string, args ...interface{}) map[string]interface{}
```

This function is used to execute Python code with variables passed from Go. The function automatically serializes the Go variables into JSON and makes them available in the Python code through the `vars` dictionary.

In the Python code template, use `$v1`, `$v2`, `$v3`, etc., as placeholders for the Go variables. These placeholders are replaced by corresponding variables from Go before execution.

#### Parameters

- `codeTemplate` (string): The Python code template where Go variables are passed.
- `args` (`...interface{}`): A variable-length argument list of Go variables to be passed to Python.

#### Returns

- `map[string]interface{}`: A map representing the result received from the Python server. This map will contain the data returned from Python through the `insyra_return` function.

#### Example

```go
message := "Hello from Go"
numbers := []int{1, 2, 3, 4}

result := RunCodef(`
print(f"Message: {$v1}")
print(f"Numbers: {$v2}")
insyra_return({"message": $v1, "value": $v2})
`, message, numbers)

fmt.Println(result)
```

## Auto-installed Dependencies

- **Python Environment**: Insyra automatically installs Python environment in the `.insyra_py_env` directory in your project root.
- **Python Libraries**: Insyra automatically installs following Python libraries:
  - `requests`
  - `json`