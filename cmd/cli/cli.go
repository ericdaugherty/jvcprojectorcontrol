package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ericdaugherty/jvcprojectorcontrol"
)

var cmds = map[string]jvcprojectorcontrol.Command{
	"NULL":   jvcprojectorcontrol.NullCommand,
	"OFF":    jvcprojectorcontrol.OffCommand,
	"ON":     jvcprojectorcontrol.OnCommand,
	"INPUT1": jvcprojectorcontrol.Input1Command,
	"INPUT2": jvcprojectorcontrol.Input2Command,
}

var hash = map[string]jvcprojectorcontrol.HashMode{
	"NONE":    jvcprojectorcontrol.HashNone,
	"JVCKW":   jvcprojectorcontrol.HashJVCKW,
	"JVCKWPJ": jvcprojectorcontrol.HashJVCKWPJ,
}

func main() {
	var ipAddress, password, hashMode, command string
	var scan, debug bool

	flag.StringVar(&ipAddress, "i", "", "IP address of the projector")
	flag.StringVar(&password, "p", "", "Password for the projector")
	flag.StringVar(&hashMode, "h", "NONE", "Hash mode for the password (NONE, JVCKW, JVCKWPJ)")
	flag.StringVar(&command, "c", "NULL", "Command to send to the projector")
	flag.BoolVar(&scan, "s", false, "Scan for projectors on the local network")
	flag.BoolVar(&debug, "d", false, "Enable debug mode")
	flag.Parse()

	cmd, exists := cmds[command]
	if !exists {
		fmt.Printf("Error: Unknown command '%s'\n", command)
		os.Exit(1)
	}

	hash, exists := hash[hashMode]
	if !exists {
		fmt.Printf("Error: Unknown hash mode '%s'\n", hashMode)
		os.Exit(1)
	}

	if scan {
		fmt.Println("Scanning for projectors...")
		projectors := jvcprojectorcontrol.ScanForProjectors(debug)

		if len(projectors) == 0 {
			fmt.Println("No projectors found")
			return
		} else if len(projectors) > 1 {
			fmt.Printf("Found %d projectors:\n", len(projectors))
			for _, proj := range projectors {
				fmt.Printf("  - %s\n", proj)
			}
			fmt.Println("Please specify an IP address using -i to send a command")
			return
		} else {
			fmt.Printf("Found projector at %s.\n", projectors[0])
			ipAddress = projectors[0]
		}
		return
	}

	if ipAddress == "" {
		fmt.Println("Error: IP address is required (use -i)")
		fmt.Println("       Or use -s to scan for projectors")
		flag.Usage()
		os.Exit(1)
	}

	proj := jvcprojectorcontrol.NewProjector(ipAddress, password, hash, debug)
	err := proj.SendCommand(cmd)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Command Successful\n")
}
