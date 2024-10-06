package py

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/HazelnutParadise/insyra"
)

// 主要邏輯
func init() {
	go startServer()

	// 如果目錄不存在，自動創建
	if _, err := os.Stat(installDir); os.IsNotExist(err) {
		os.MkdirAll(installDir, os.ModePerm)
	} else {
		// 檢查 Python 執行檔是否已存在
		if _, err := os.Stat(pyPath); err == nil {
			insyra.LogDebug("Python installation already exists!")
			err = installDependencies()
			if err != nil {
				insyra.LogFatal("Failed to install dependencies: %v", err)
			}
			insyra.LogInfo("Dependencies installation completed successfully!")
			return
		}
	}

	// 安裝 Python
	err := installPython(pythonVersion)
	if err != nil {
		insyra.LogFatal("py.init: Failed to install Python: %v", err)
	}
	insyra.LogInfo("py.init: Python installation completed successfully!")
	err = installDependencies()
	if err != nil {
		insyra.LogFatal("py.init: Failed to install dependencies: %v", err)
	}
	insyra.LogInfo("py.init: Dependencies installation completed successfully!")
}

// 下載並安裝 Python 的邏輯
func installPython(version string) error {

	if runtime.GOOS == "windows" {
		return installPythonOnWindows(version, absInstallDir)
	}
	return installPythonOnUnix(version, absInstallDir)
}

// Windows 平台安裝
func installPythonOnWindows(version string, installDir string) error {
	downloadURL := fmt.Sprintf("https://www.python.org/ftp/python/%s/python-%s-amd64.exe", version, version)
	installerPath := filepath.Join(os.TempDir(), fmt.Sprintf("python-%s-installer.exe", version))

	fmt.Println("Downloading Python installer for Windows...")
	err := downloadFile(installerPath, downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download installer: %w", err)
	}

	fmt.Println("Running Python installer...")
	cmd := exec.Command(installerPath, "/quiet", "InstallAllUsers=1", fmt.Sprintf("TargetDir=%s", installDir))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Unix-like 平台安裝 (Linux/macOS)
func installPythonOnUnix(version string, installDir string) error {
	downloadURL := fmt.Sprintf("https://www.python.org/ftp/python/%s/Python-%s.tgz", version, version)
	pythonTar := filepath.Join(os.TempDir(), fmt.Sprintf("Python-%s.tgz", version))

	fmt.Println("Downloading Python for Unix-like systems...")
	err := downloadFile(pythonTar, downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download Python: %w", err)
	}

	fmt.Println("Extracting Python files...")
	err = extractTar(pythonTar, os.TempDir())
	if err != nil {
		return fmt.Errorf("failed to extract Python: %w", err)
	}

	pythonSrcDir := filepath.Join(os.TempDir(), fmt.Sprintf("Python-%s", version))

	fmt.Println("Configuring and installing Python...")
	return installPythonFromSource(pythonSrcDir, installDir)
}

// 下載檔案
func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// 解壓縮 .tgz（適用於 Unix-like 系統）
func extractTar(filepath string, destDir string) error {
	cmd := exec.Command("tar", "-xvzf", filepath, "-C", destDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// 從原始碼編譯並安裝 Python（適用於 Unix-like 系統）
func installPythonFromSource(srcDir string, installDir string) error {
	configureCmd := exec.Command("./configure", fmt.Sprintf("--prefix=%s", installDir))
	configureCmd.Dir = srcDir
	configureCmd.Stdout = os.Stdout
	configureCmd.Stderr = os.Stderr
	if err := configureCmd.Run(); err != nil {
		return err
	}

	makeCmd := exec.Command("make")
	makeCmd.Dir = srcDir
	makeCmd.Stdout = os.Stdout
	makeCmd.Stderr = os.Stderr
	if err := makeCmd.Run(); err != nil {
		return err
	}

	installCmd := exec.Command("make", "install")
	installCmd.Dir = srcDir
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	return installCmd.Run()
}

func installDependencies() error {
	cmd := exec.Command(pyPath, "-m", "pip", "install", "requests")
	cmd.Dir = absInstallDir
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	return cmd.Run()
}