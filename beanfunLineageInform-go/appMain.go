package main

import (
	"bufio"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func main() {
	// find Lineage launcher Process
	isFindWindow, windowHandle := findWindow("Lineage launcher")

	if isFindWindow {
		// get window Thread Process Id
		windowProcessThreadID := getWindowThreadProcessID(windowHandle)
		fmt.Printf("Lineage launcher PID: %d\n", windowProcessThreadID)
		// Open Process
		isOpenSuccess, handle := openProcess(windowProcessThreadID)

		if isOpenSuccess {
			var address uint32 = 0x5574a6 // data address
			var offset uint32 = 0x02      // offset
			fmt.Println()
			fmt.Print("Beanfun Account/Password/Server: ")
			// Read data
			for {
				// Read Process Memory
				data := readProcessMemory(handle, address+offset)
				fmt.Print(string(data))
				// end of data
				if data == 0 {
					fmt.Println()
					break
				}
				// get another byte
				offset += 0x01
			}
			fmt.Println()
			fmt.Println("Developed by L1J-TW")
		} else {
			fmt.Println("Could not get handle!")
			fmt.Println("請使用系統管理員啟動!")
		}
	} else {
		fmt.Println("Window not found!")
	}
	// Wait for Enter
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func readProcessMemory(handle uintptr, address uint32) byte {
	var data byte
	var kernelModule = syscall.NewLazyDLL("kernel32.dll")
	var readProcessMemoryProcess = kernelModule.NewProc("ReadProcessMemory")
	readProcessMemoryProcess.Call(uintptr(handle), uintptr(address), uintptr(unsafe.Pointer(&data)), unsafe.Sizeof(data), 0)
	return data
}

func openProcess(processID int) (bool, uintptr) {
	const processAllAccess = 0x1f0fff
	var kernelModule = syscall.NewLazyDLL("kernel32.dll")
	var openProcessProcess = kernelModule.NewProc("OpenProcess")
	hwnd, _, _ := openProcessProcess.Call(uintptr(processAllAccess), uintptr(0), uintptr(processID))
	return hwnd != 0, hwnd
}

func getWindowThreadProcessID(windowHandle uintptr) int {
	// get window Thread Process id
	var windowProcessID int
	// Load dll
	var userModule = syscall.NewLazyDLL("user32.dll")
	// Create Process
	var getWindowThreadProcessIDProcess = userModule.NewProc("GetWindowThreadProcessId")
	getWindowThreadProcessIDProcess.Call(uintptr(windowHandle), uintptr(unsafe.Pointer(&windowProcessID)))
	return windowProcessID
}

func findWindow(windowName string) (bool, uintptr) {
	// Load dll
	var userModule = syscall.NewLazyDLL("user32.dll")
	// Create Process
	var findWindowProcess = userModule.NewProc("FindWindowW")
	// Call FindWindowW windows api
	result, _, _ := findWindowProcess.Call(0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(windowName))))
	return result != 0, result
}
