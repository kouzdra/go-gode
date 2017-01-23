package ide

import "log"
import "os"
import "os/exec"

func (ide *IDE) Make () {
	log.Printf ("++++++++++++++++ BUILD ++++++++++++++++++")
	cmd := os.ExpandEnv ("make -C $GOPATH/src/github.com/kouzdra/go-gode")
	log.Printf ("%s\n", cmd)
	exec.Command (cmd).Run ()
	log.Printf ("---------------- BUILD ------------------")
}
