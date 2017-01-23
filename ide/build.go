package ide

import "log"
//import "bytes"
import "strings"
import "os"
import "os/exec"

func (ide *IDE) Make () {
	log.Printf ("++++++++++++++++ BUILD ++++++++++++++++++")
	cmd := exec.Command ("make", "-C", os.ExpandEnv ("$GOPATH/src/github.com/kouzdra/go-gode"), "gode")
	cmd.Stdin  = strings.NewReader ("")
	cmd.Stdout = os.Stdout
	//var out bytes.Buffer
	//cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Printf("Error: %s", err)
	}
	//log.Printf("%s", out.String())
	log.Printf ("---------------- BUILD ------------------")
}
