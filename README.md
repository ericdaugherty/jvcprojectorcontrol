# JVC Projector Remote Control Library

This package provides the ability to remotely control a JVC Projector via TCP/IP.

It provides both a Go library that can be used by other applications, as well as a stand alone command line tool that can be used directly.

## JVC Projector Setup

The JVC Projector must be connected to your local network using an Ethernet cable.

In the projector, you must enable the LAN communication and (if applicable) set the Communication Terminal to LAN (only for projectors with RS232 ports)

For NZ projectors, a password must be set that is between 8 and 10 characters in length and must be used when connecting to the projector. 

The NZ500/700 requires the password to be hashed.

NX projectors do not require a password.

## Commands

Currently Supported Commands:
- NULL - Test command to verify connectivity.
- OFF - Turn Projector off
- ON - Turn Projector on
- Input1 - Select Input 1
- Input2 - Select Input 2

Additional commands can be added easily once you have the correct bytes required from the documentation. 

Docs:
- [NZ500, NZ700, RS1200, RS2200](https://manuals.jvckenwood.com/download/files/B5A-4685-11.pdf) Page 75

## Command Line

The command line tool takes several arguments. 

`-i 192.160.0.1` Specify the IP address of the projector (if known)

`-s` Scan subnet for projectors. Will send the command if one found, or return a list of IPs if multiple are found.

Command: (`-c NULL`)
- NULL
- OFF
- ON
- INPUT1
- INPUT2

`-p <PASSWORD>` Specify the password (NZ projectors)

Password Hash: (`-h NONE`)
- NONE - For NX projectors
- JVCKW - For most NZ projectors
- JVCKWPJ For NZ500/NZ700 and later projectors

`-d` enabled debug mode to see all traffic to/from the projector.

## Resources

Similar Projects
- [Python Library](https://github.com/bezmi/jvc_projector)
- [Python Library](https://github.com/iloveicedgreentea/pyjvcprojector)
- [Android App](https://github.com/LaUs3r/JVCProjectorRemoteControl)
