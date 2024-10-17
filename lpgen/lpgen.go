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

// set objective
func (lp *LPModel) SetObjective(objType, obj string) *LPModel {
	lp.ObjectiveType = objType
	lp.Objective = obj
	return lp
}

// add constraint
func (lp *LPModel) AddConstraint(constr string) *LPModel {
	lp.Constraints = append(lp.Constraints, constr)
	return lp
}

// add bound
func (lp *LPModel) AddBound(bound string) *LPModel {
	lp.Bounds = append(lp.Bounds, bound)
	return lp
}

// add binary var
func (lp *LPModel) AddBinaryVar(varName string) *LPModel {
	// 確保變數不會同時是整數和二進制
	for i, intVar := range lp.IntegerVars {
		if intVar == varName {
			lp.IntegerVars = append(lp.IntegerVars[:i], lp.IntegerVars[i+1:]...)
			break
		}
	}
	lp.BinaryVars = append(lp.BinaryVars, varName)
	return lp
}

// add integer var
func (lp *LPModel) AddIntegerVar(varName string) *LPModel {
	// 確保變數不會同時是整數和二進制
	for i, binVar := range lp.BinaryVars {
		if binVar == varName {
			lp.BinaryVars = append(lp.BinaryVars[:i], lp.BinaryVars[i+1:]...)
			break
		}
	}
	lp.IntegerVars = append(lp.IntegerVars, varName)
	return lp
}

// generate lp file
func (lp *LPModel) GenerateLPFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		insyra.LogWarning("lpgen.GenerateLPFile: Failed to create LP file: %v", err)
		return
	}
	defer file.Close()

	// 設定最大化或最小化
	_, err = file.WriteString(lp.ObjectiveType + "\n")
	if err != nil {
		insyra.LogWarning("lpgen.GenerateLPFile: Failed to write objective type: %v, returning", err)
		return
	}
	_, err = file.WriteString("  " + "obj: " + lp.Objective + "\n")
	if err != nil {
		insyra.LogWarning("lpgen.GenerateLPFile: Failed to write objective: %v, returning", err)
		return
	}

	// 添加約束條件
	_, err = file.WriteString("Subject To\n")
	if err != nil {
		insyra.LogWarning("lpgen.GenerateLPFile: Failed to write subject to: %v, returning", err)
		return
	}
	for i, constr := range lp.Constraints {
		_, err = file.WriteString("  " + fmt.Sprintf("c%d: %s\n", i+1, constr))
		if err != nil {
			insyra.LogWarning("lpgen.GenerateLPFile: Failed to write constraint: %v, returning", err)
			return
		}
	}

	// 添加變數邊界
	if len(lp.Bounds) > 0 {
		_, err = file.WriteString("Bounds\n")
		if err != nil {
			insyra.LogWarning("lpgen.GenerateLPFile: Failed to write bounds: %v, returning", err)
			return
		}
		for _, bound := range lp.Bounds {
			_, err = file.WriteString("  " + bound + "\n")
			if err != nil {
				insyra.LogWarning("lpgen.GenerateLPFile: Failed to write bound: %v, returning", err)
				return
			}
		}
	}

	// 添加整數變數
	if len(lp.IntegerVars) > 0 {
		file.WriteString("General\n")
		for _, intVar := range lp.IntegerVars {
			_, err = file.WriteString("  " + intVar + "\n")
			if err != nil {
				insyra.LogWarning("lpgen.GenerateLPFile: Failed to write integer var: %v, returning", err)
				return
			}
		}
	}

	// 添加二進制變數
	if len(lp.BinaryVars) > 0 {
		file.WriteString("Binary\n")
		for _, binVar := range lp.BinaryVars {
			_, err = file.WriteString("  " + binVar + "\n")
			if err != nil {
				insyra.LogWarning("lpgen.GenerateLPFile: Failed to write binary var: %v, returning", err)
				return
			}
		}
	}

	// LP 文件結尾
	_, err = file.WriteString("End\n")
	if err != nil {
		insyra.LogWarning("lpgen.GenerateLPFile: Failed to write end: %v, returning", err)
		return
	}
}
