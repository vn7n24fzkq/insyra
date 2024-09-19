package lpgen

import (
	"fmt"
	"os"

	"github.com/HazelnutParadise/insyra"
)

type LPModel struct {
	Objective     string   // 目標函數
	ObjectiveType string   // 目標函數類型（Maximize or Minimize）
	Constraints   []string // 約束條件
	Bounds        []string // 變數邊界
	BinaryVars    []string // 二進制變數
	IntegerVars   []string // 整數變數
}

func NewLPModel() *LPModel {
	return &LPModel{}
}

// 設定目標函數（最大化或最小化）
func (lp *LPModel) SetObjective(objType, obj string) {
	lp.ObjectiveType = objType
	lp.Objective = obj
}

// 新增約束條件
func (lp *LPModel) AddConstraint(constr string) {
	lp.Constraints = append(lp.Constraints, constr)
}

// 新增變數邊界
func (lp *LPModel) AddBound(bound string) {
	lp.Bounds = append(lp.Bounds, bound)
}

// 新增二進制變數
func (lp *LPModel) AddBinaryVar(varName string) {
	// 確保變數不會同時是整數和二進制
	for i, intVar := range lp.IntegerVars {
		if intVar == varName {
			lp.IntegerVars = append(lp.IntegerVars[:i], lp.IntegerVars[i+1:]...)
			break
		}
	}
	lp.BinaryVars = append(lp.BinaryVars, varName)
}

// 新增整數變數
func (lp *LPModel) AddIntegerVar(varName string) {
	// 確保變數不會同時是整數和二進制
	for i, binVar := range lp.BinaryVars {
		if binVar == varName {
			lp.BinaryVars = append(lp.BinaryVars[:i], lp.BinaryVars[i+1:]...)
			break
		}
	}
	lp.IntegerVars = append(lp.IntegerVars, varName)
}

// 生成 LP 檔案
func (lp *LPModel) GenerateLPFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		insyra.LogFatal("lpgen.GenerateLPFile: Failed to create LP file: %v", err)
	}
	defer file.Close()

	// 設定最大化或最小化
	file.WriteString(lp.ObjectiveType + "\n")
	file.WriteString("  " + lp.Objective + "\n")

	// 添加約束條件
	file.WriteString("Subject To\n")
	for _, constr := range lp.Constraints {
		file.WriteString("  " + fmt.Sprintf("%s\n", constr))
	}

	// 添加變數邊界
	if len(lp.Bounds) > 0 {
		file.WriteString("Bounds\n")
		for _, bound := range lp.Bounds {
			file.WriteString("  " + bound + "\n")
		}
	}

	// 添加整數變數
	if len(lp.IntegerVars) > 0 {
		file.WriteString("General\n")
		for _, intVar := range lp.IntegerVars {
			file.WriteString("  " + intVar + "\n")
		}
	}

	// 添加二進制變數
	if len(lp.BinaryVars) > 0 {
		file.WriteString("Binary\n")
		for _, binVar := range lp.BinaryVars {
			file.WriteString("  " + binVar + "\n")
		}
	}

	// LP 文件結尾
	file.WriteString("End\n")
}