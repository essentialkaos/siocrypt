package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/essentialkaos/ek/v13/errors"
	"github.com/essentialkaos/ek/v13/fmtc"
	"github.com/essentialkaos/ek/v13/fsutil"
	"github.com/essentialkaos/ek/v13/options"
	"github.com/essentialkaos/ek/v13/secstr"
	"github.com/essentialkaos/ek/v13/support"
	"github.com/essentialkaos/ek/v13/support/deps"
	"github.com/essentialkaos/ek/v13/system/procname"
	"github.com/essentialkaos/ek/v13/terminal"
	"github.com/essentialkaos/ek/v13/terminal/tty"
	"github.com/essentialkaos/ek/v13/usage"
	"github.com/essentialkaos/ek/v13/usage/completion/bash"
	"github.com/essentialkaos/ek/v13/usage/completion/fish"
	"github.com/essentialkaos/ek/v13/usage/completion/zsh"
	"github.com/essentialkaos/ek/v13/usage/man"
	"github.com/essentialkaos/ek/v13/usage/update"

	"golang.org/x/crypto/scrypt"

	"github.com/essentialkaos/sio"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Basic utility info
const (
	APP  = "siocrypt"
	VER  = "0.1.0"
	DESC = "Tool for encrypting/decrypting arbitrary data streams"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Options
const (
	OPT_DECRYPT  = "D:decrypt"
	OPT_CIPHER   = "c:cipher"
	OPT_PASSWORD = "p:password"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_UPDATE       = "U:update"
	OPT_VERB_VER     = "vv:verbose-version"
	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
)

const (
	CIPHER_AES256   = "AES256"   // AES-256 GCM
	CIPHER_C20P1305 = "C20P1305" // ChaCha20 Poly1305
)

const ENV_PASSWORD = "SIOCRYPT_PASSWORD"

// ////////////////////////////////////////////////////////////////////////////////// //

// optMap contains information about all supported options
var optMap = options.Map{
	OPT_DECRYPT:  {Type: options.BOOL},
	OPT_CIPHER:   {Value: CIPHER_C20P1305},
	OPT_PASSWORD: {},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL},
	OPT_VER:      {Type: options.MIXED},

	OPT_UPDATE:       {Type: options.MIXED},
	OPT_VERB_VER:     {Type: options.BOOL},
	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

// color tags for app name and version
var colorTagApp, colorTagVer string

// password is secure string with password
var password *secstr.String

// removeOutputOnError is flag to remove output file on error
var removeOutputOnError bool

// ////////////////////////////////////////////////////////////////////////////////// //

// Run is main utility function
func Run(gitRev string, gomod []byte) {
	preConfigureUI()

	args, errs := options.Parse(optMap)

	if !errs.IsEmpty() {
		terminal.Error("Options parsing errors:")
		terminal.Error(errs.Error(" - "))
		os.Exit(1)
	}

	configureUI()

	switch {
	case options.Has(OPT_COMPLETION):
		os.Exit(printCompletion())
	case options.Has(OPT_GENERATE_MAN):
		printMan()
		os.Exit(0)
	case options.GetB(OPT_VER):
		genAbout(gitRev).Print(options.GetS(OPT_VER))
		os.Exit(0)
	case options.GetB(OPT_VERB_VER):
		support.Collect(APP, VER).
			WithRevision(gitRev).
			WithDeps(deps.Extract(gomod)).
			Print()
		os.Exit(0)
	case withSelfUpdate && options.GetB(OPT_UPDATE):
		os.Exit(updateBinary())
	case options.GetB(OPT_HELP) || (!hasStdinData() && len(args) == 0):
		genUsage().Print()
		os.Exit(0)
	}

	err := errors.Chain(
		processPassword,
		validateOptions,
	)

	if err != nil {
		terminal.Error(err)
		os.Exit(1)
	}

	err = process(args)

	if err != nil {
		terminal.Error(err)
		os.Exit(1)
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// preConfigureUI preconfigures UI based on information about user terminal
func preConfigureUI() {
	if !tty.IsTTY() {
		fmtc.DisableColors = true
	}

	switch {
	case fmtc.IsTrueColorSupported():
		colorTagApp, colorTagVer = "{*}{#87D787}", "{#87D787}"
	case fmtc.Is256ColorsSupported():
		colorTagApp, colorTagVer = "{*}{#114}", "{#114}"
	default:
		colorTagApp, colorTagVer = "{*}{g}", "{g}"
	}
}

// configureUI configures user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}
}

// processPassword tries to securily read password from options or env vars
func processPassword() error {
	var err error

	pwd := options.GetS(OPT_PASSWORD)

	if pwd != "" {
		err = procname.Replace(pwd, strings.Repeat("*", len(pwd)))

		if err != nil {
			return fmt.Errorf("Can't hide provided password: %v", err)
		}

		options.Delete(OPT_PASSWORD)
	} else if os.Getenv(ENV_PASSWORD) != "" {
		pwd = os.Getenv(ENV_PASSWORD)
		os.Setenv(ENV_PASSWORD, "")
	}

	if pwd == "" {
		return fmt.Errorf("Password is not set")
	}

	password, err = secstr.NewSecureString(pwd)

	if err != nil {
		return fmt.Errorf("Can't create secure string for password: %v", err)
	}

	return nil
}

// validateOptions validates options
func validateOptions() error {
	if options.Has(OPT_CIPHER) {
		switch strings.ToUpper(options.GetS(OPT_CIPHER)) {
		case CIPHER_AES256, CIPHER_C20P1305:
			// ok
		default:
			return fmt.Errorf("Unsupported cipher suite %q", options.GetS(OPT_CIPHER))
		}
	}

	return nil
}

// process starts arguments processing
func process(args options.Arguments) error {
	input, output, err := getIO(args)

	if err != nil {
		return err
	}

	key, err := deriveKey(input, output)

	if err != nil {
		return err
	}

	cfg := sio.Config{Key: key, CipherSuites: getCipherSuite()}

	if options.GetB(OPT_DECRYPT) {
		_, err = sio.Decrypt(output, input, cfg)

		if err != nil {
			cleanupOutput(output.Name())
			return fmt.Errorf("Error while decrypting data: %v", err)
		}
	} else {
		_, err = sio.Encrypt(output, input, cfg)

		if err != nil {
			cleanupOutput(output.Name())
			return fmt.Errorf("Error while encrypting data: %v", err)
		}
	}

	return nil
}

// getIO returns input and output targets
func getIO(args options.Arguments) (*os.File, *os.File, error) {
	var err error
	var input, output *os.File
	var inputFile, outputFile string

	switch len(args) {
	case 0:
		input = os.Stdin
		output = os.Stdout
	case 1:
		inputFile = args.Get(0).Clean().String()
		output = os.Stdout
	default:
		if args.Get(0).Is("--") {
			input = os.Stdin
		} else {
			inputFile = args.Get(0).Clean().String()
		}

		outputFile = args.Get(1).String()
	}

	if inputFile != "" {
		err = fsutil.ValidatePerms("FRS", inputFile)

		if err != nil {
			return nil, nil, fmt.Errorf("Can't use given input: %v", err)
		}

		input, err = os.OpenFile(inputFile, os.O_RDONLY, 0)

		if err != nil {
			return nil, nil, fmt.Errorf("Can't use given input: %v", err)
		}
	}

	if outputFile != "" {
		output, err = os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0600)

		if err != nil {
			return nil, nil, fmt.Errorf("Can't use given output: %v", err)
		}

		removeOutputOnError = true
	}

	return input, output, err
}

// deriveKey creates derived key from password and salt
func deriveKey(input, output *os.File) ([]byte, error) {
	salt := make([]byte, 32)

	if options.GetB(OPT_DECRYPT) {
		_, err := io.ReadFull(input, salt)

		if err != nil {
			return nil, fmt.Errorf("Can't read salt from encrypted data: %v", err)
		}
	} else {
		_, err := io.ReadFull(rand.Reader, salt)

		if err != nil {
			return nil, fmt.Errorf("Can't generate random salt: %v", err)
		}

		_, err = output.Write(salt)

		if err != nil {
			return nil, fmt.Errorf("Can't write salt to output: %v", err)
		}
	}

	key, err := scrypt.Key(password.Data, salt, 32768, 16, 1, 32)

	if err != nil {
		return nil, fmt.Errorf("Can't derive key from password: %v", err)
	}

	password.Destroy()

	return key, nil
}

// getCipherSuite returns target cipher suite
func getCipherSuite() []byte {
	switch strings.ToUpper(options.GetS(OPT_CIPHER)) {
	case CIPHER_AES256:
		return []byte{sio.AES_256_GCM}
	case CIPHER_C20P1305:
		return []byte{sio.CHACHA20_POLY1305}
	}

	return []byte{}
}

// hasStdinData return true if there is some data in stdin
func hasStdinData() bool {
	stdin, err := os.Stdin.Stat()

	if err != nil {
		return false
	}

	if stdin.Mode()&os.ModeCharDevice != 0 {
		return false
	}

	return true
}

// cleanupOutput removes output file on error
func cleanupOutput(file string) {
	if removeOutputOnError {
		os.Remove(file)
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// printCompletion prints completion for given shell
func printCompletion() int {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Print(bash.Generate(info, APP))
	case "fish":
		fmt.Print(fish.Generate(info, APP))
	case "zsh":
		fmt.Print(zsh.Generate(info, optMap, APP))
	default:
		return 1
	}

	return 0
}

// printMan prints man page
func printMan() {
	fmt.Println(man.Generate(genUsage(), genAbout("")))
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("", "?input", "?output")

	info.AppNameColorTag = colorTagApp

	info.Spoiler = `  Without an input file, tool reads data from standard input and writes to standard
  output if no output file is specified.

  Password provided via {g}--password{!}/{g}-p{!} option will be masked in process command. Also, 
  password can't be specified using {m}SIOCRYPT_PASSWORD{!} environment variable.`

	info.AddOption(OPT_DECRYPT, "Decrypt data")
	info.AddOption(OPT_PASSWORD, "Password for encrypting/decrypting", "password")
	info.AddOption(OPT_CIPHER, "Cipher to use {s}(AES256/{_}C20P1305{!_}){!}", "cipher")

	if withSelfUpdate {
		info.AddOption(OPT_UPDATE, "Update application to the latest version")
	}

	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddEnv(ENV_PASSWORD, "Password for encrypting/decrypting {s-}(String){!}")

	info.AddRawExample(
		"cat data.enc | siocrypt -p TeSt1234 | grep John",
		"Read data from standard input and grep over decrypted output",
	)

	info.AddRawExample(
		"cat data.enc | siocrypt -p TeSt1234 -- data.txt",
		"Read data from standard input and save decrypted data as data.txt",
	)

	info.AddExample(
		"-p TeSt1234 data.enc | grep John",
		"Read data from file data.enc and grep over decrypted output",
	)

	info.AddExample(
		"-p TeSt1234 data.enc data.txt",
		"Read data from file data.enc and save decrypted data as data.txt",
	)

	info.AddRawExample(
		"SIOCRYPT_PASSWORD=TeSt1234 siocrypt data.enc data.txt",
		"Use password from environment variable while decrypting data",
	)

	return info
}

// genAbout generates info about version
func genAbout(gitRev string) *usage.About {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2009,
		Owner:   "ESSENTIAL KAOS",

		AppNameColorTag: colorTagApp,
		VersionColorTag: colorTagVer,
		DescSeparator:   "{s}—{!}",

		License:       "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
		BugTracker:    "https://github.com/essentialkaos/siocrypt/issues",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/siocrypt", update.GitHubChecker},
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
	}

	return about
}

// ////////////////////////////////////////////////////////////////////////////////// //
