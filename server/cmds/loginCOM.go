package cmds

import "flag"

const LoginCommandName string = "LOGIN"

var loginFlag *flag.FlagSet = flag.NewFlagSet(LoginCommandName, flag.ContinueOnError)

//////////////////////////////////////////////////////////////
//                                                          //
//////////////////////////////////////////////////////////////

type LoginModel struct {
	UserId string
}

//////////////////////////////////////////////////////////////
//                                                          //
//////////////////////////////////////////////////////////////

type LoginModelProvider struct {
}

//////////////////////////////////////////////////////////////
//                                                          //
//////////////////////////////////////////////////////////////

/*
func (cmd commandLogin) Execute(conn *Conn, param string) {

}
*/
