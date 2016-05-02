package main

import (
	"bufio"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const PROCESS_ALL_ACCESS = 0x1f0fff

func main() {
	// Load dll
	var userModule = syscall.NewLazyDLL("user32.dll")
	var kernelModule = syscall.NewLazyDLL("kernel32.dll")

	// Call windows api
	// find Lineage launcher Process
	var findWindowProcess = userModule.NewProc("FindWindowW")
	hwnd, _, _ := findWindowProcess.Call(0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Lineage launcher"))))

	if hwnd <= 0 {
		fmt.Println("Window not found!")
	} else {
		// get Process id
		var processId int
		var getWindowThreadProcessIDProcess = userModule.NewProc("GetWindowThreadProcessId")
		getWindowThreadProcessIDProcess.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&processId)))
		fmt.Printf("Lineage launcher PID: %d\n", processId)

		// OpenProcess
		var openProcessProcess = kernelModule.NewProc("OpenProcess")
		hwnd, _, _ := openProcessProcess.Call(uintptr(PROCESS_ALL_ACCESS), uintptr(0), uintptr(processId))

		if hwnd <= 0 {
			fmt.Println("Could not get handle!")
			fmt.Println("請使用系統管理員啟動!")
		} else {
			var data byte
			var address uint32 = 0x5574a6 // data address
			var offset uint32 = 0x02      // offset
			fmt.Println()
			fmt.Print("Beanfun Account/Password/Server: ")
			// Read data
			for {
				// Read Process Memory
				var readProcessMemoryProcess = kernelModule.NewProc("ReadProcessMemory")
				readProcessMemoryProcess.Call(uintptr(hwnd), uintptr(address+offset), uintptr(unsafe.Pointer(&data)), unsafe.Sizeof(data), 0)
				fmt.Print(string(data))
				// get another byte
				offset += 0x01
				// end of data
				if data == 0 {
					fmt.Println()
					break
				}
			}
			fmt.Println()
			fmt.Println("Developed by L1J-TW")
		}
	}
	// Wait for Enter
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
