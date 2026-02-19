package cli

import (
	"fmt"
)

// Códigos ANSI para colores (funcionan en Mac/Linux/Windows moderno)
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	Bold        = "\033[1m"
	ColorGray   = "\033[90m"
)

func PrintBanner() {
	fmt.Println(ColorPurple + Bold)
	fmt.Println(`
    ____                            __  ______
   / __ \_________  ____ ___  ____ / /_/ ____/
  / /_/ / ___/ __ \/ __  __ \/ __ / __/ /     
 / ____/ /  / /_/ / / / / / / /_/ / /_/ /___  
/_/   /_/   \____/_/ /_/ /_/ .___/\__/\____/  
                          /_/                 
	` + ColorReset)
	fmt.Println(ColorCyan + "   The Prompt Compiler for Engineering Excellence" + ColorReset)
	fmt.Println(ColorCyan + "   v0.1.0-alpha • by Cesar Rivas" + ColorReset)
	fmt.Println()
}

func PrintSuccess(msg string) {
	fmt.Printf("%s✔ %s%s\n", ColorGreen, msg, ColorReset)
}

func PrintError(msg string) {
	fmt.Printf("%s✖ %s%s\n", ColorRed, msg, ColorReset)
}

func PrintWarning(msg string) {
	fmt.Printf("%s⚠ %s%s\n", ColorYellow, msg, ColorReset)
}

func PrintInfo(msg string) {
	fmt.Printf("%sℹ %s%s\n", ColorCyan, msg, ColorReset)
}

func PrintSection(title string) {
	fmt.Println()
	fmt.Printf("%s%s%s\n", Bold+ColorCyan, title, ColorReset)
	fmt.Println(ColorGray + "--------------------------------------------------" + ColorReset)
}
