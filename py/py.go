package py

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Run the Python code and return the result.
func RunCode(code string) map[string]interface{} {
	code = generateDefaultPyCode() + code

	// 執行 Python 代碼
	pythonCmd := exec.Command(pyPath, "-c", code)
	pythonCmd.Dir = installDir
	pythonCmd.Stdout = os.Stdout
	pythonCmd.Stderr = os.Stderr
	err := pythonCmd.Run()
	if err != nil {
		log.Fatalf("Failed to run Python code: %v", err)
	}

	// 從 server 接收資料
	return pyResult
}

// Run the Python code with the given Golang variables and return the result.
// The variables should be passed by $v1, $v2, $v3... in the codeTemplate.
func RunCodef(codeTemplate string, args ...interface{}) map[string]interface{} {
	// 產生包含 insyra_return 函數的預設 Python 代碼
	codeTemplate = generateDefaultPyCode() + codeTemplate

	// 將所有 Go 變數轉換為一個 JSON 字典傳遞給 Python
	dataMap := make(map[string]interface{})
	for i, arg := range args {
		pythonVarName := fmt.Sprintf("var%d", i+1)
		dataMap[pythonVarName] = arg
	}

	// 將字典序列化為 JSON
	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		log.Fatalf("Failed to serialize variables: %v", err)
	}

	// 替換 codeTemplate 中的 $v1, $v2... 佔位符為對應的 vars['var1'], vars['var2']...
	for i := range args {
		placeholder := fmt.Sprintf("$v%d", i+1)
		codeTemplate = strings.ReplaceAll(codeTemplate, placeholder, fmt.Sprintf("vars['var%d']", i+1))
	}

	// 在 Python 中生成變數賦值語句
	pythonVarCode := fmt.Sprintf("vars = json.loads('%s')", string(jsonData))

	// 構建完整的 Python 代碼，並確保正確導入 json 模組
	fullCode := fmt.Sprintf(`
import json
%s
%s
`, pythonVarCode, codeTemplate)

	// 執行 Python 代碼
	pythonCmd := exec.Command(pyPath, "-c", fullCode)
	pythonCmd.Stdout = os.Stdout
	pythonCmd.Stderr = os.Stderr

	// 執行 Python 命令
	err = pythonCmd.Run()
	if err != nil {
		log.Fatalf("Failed to run Python code: %v", err)
	}

	// 從 server 接收資料
	return pyResult
}

func generateDefaultPyCode() string {
	return fmt.Sprintf(`
import requests
import json

def insyra_return(data, url="http://localhost:%v/pyresult"):
    response = requests.post(url, json=data)
    if response.status_code != 200:
        print(f"Failed to send result: {response.status_code}")

`, port)
}