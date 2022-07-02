package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	author  string = "seaung"
	version string = "1.0.x-dev"
)

func checkSudo() {
	if os.Geteuid() != 0 {
		New().LoggerError("This program need to have root permission to execute for now!")
		os.Exit(1)
	}
}

func showBanner() {
	name := fmt.Sprintf("Nox (v.%s)", version)
	banner := `
         ,--.                      
       ,--.'|                      
   ,--,:  : |                      
,'--.''|  ' :   ,---.              
|   :  :  | |  '   ,'\ ,--,  ,--,  
:   |   \ | : /   /   ||'. \/ .''|  
|   : '  '; |.   ; ,. :'  \/  / ;  
'   ' ;.    ;'   | |: : \  \.' /   
|   | | \   |'   | .; :  \  ;  ;   
'   : |  ; .'|   :    | / \  \  \  
|   | ''--'   \   \  /./__;   ;  \ 
'   : |        ''----' |   :/\  \ ; 
;   |.'                '---'  '--'  
'---'
	`

	lines := strings.Split(banner, "\n")
	width := len(lines[1])

	fmt.Println(banner)
	color.Green(fmt.Sprintf("%[1]*s", -width, fmt.Sprintf("%[1]*s", (width+len(name))/2, name)))
	color.Blue(fmt.Sprintf("%[1]*s", -width, fmt.Sprintf("%[1]*s", (width+len(author))/2, author)))
	fmt.Println()
}

func InitConsole() {
	checkSudo()
	showBanner()
}
