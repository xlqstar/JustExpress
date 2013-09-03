/*package main

import (
	"fmt"
)

func main() {
	fmt.Print("Input Name:")
	var name string
	// for i := 0; i < 4; i++ {
	n, err := fmt.Scanf("%q", &name)
	// n, err := fmt.Scan(&name)
	fmt.Println(n, err, name)
	// }
}*/

/*package main

import (
	"fmt"
)

func main() {
	var arg0, arg1, arg2, arg3, arg4 string
	fmt.Print("Input Name:")
	for i := 0; i < 3; i++ {
		n, err := fmt.Scanln(&arg0, &arg1, &arg2, &arg3, &arg4)
		fmt.Println(n, err, arg0, arg1, arg2, arg3, arg4)
	}
}*/
/*
package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	log.Println("haha")
	data, _, _ := reader.ReadLine()
	command := string(data)
	log.Println("command", command)
}*/

package main

import (
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("explorer.exe", "/select,F:\\kuaipan\\Projects\\liqun3.me\\wahaha", "blog@2013-9-1\\article.md")
	log.Println(cmd.Args, len(cmd.Args))
	cmd.Run()
}
