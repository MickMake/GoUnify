module github.com/MickMake/GoUnify/cmdVersion

go 1.18

replace github.com/MickMake/GoUnify/Only => ../Only

replace github.com/MickMake/GoUnify/cmdHelp => ../cmdHelp

replace github.com/MickMake/GoUnify/cmdLog => ../cmdLog

require (
	github.com/MickMake/GoUnify/Only v0.0.0-00010101000000-000000000000
	github.com/MickMake/GoUnify/cmdHelp v0.0.0-00010101000000-000000000000
	github.com/MickMake/GoUnify/cmdLog v0.0.0-00010101000000-000000000000
	github.com/blang/semver v3.5.1+incompatible
	github.com/google/go-github/v30 v30.1.0
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/rhysd/go-github-selfupdate v1.2.3
	github.com/spf13/cobra v1.5.0
	github.com/spf13/pflag v1.0.5
	github.com/tcnksm/go-gitconfig v0.1.2
	golang.org/x/oauth2 v0.0.0-20220909003341-f21342109be1
)

require (
	github.com/MichaelMure/go-term-markdown v0.1.4 // indirect
	github.com/MichaelMure/go-term-text v0.3.1 // indirect
	github.com/alecthomas/chroma v0.7.1 // indirect
	github.com/danwakefield/fnmatch v0.0.0-20160403171240-cbb64ac3d964 // indirect
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/dlclark/regexp2 v1.1.6 // indirect
	github.com/eliukblau/pixterm/pkg/ansimage v0.0.0-20191210081756-9fb6cf8c2f75 // indirect
	github.com/fatih/color v1.9.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gomarkdown/markdown v0.0.0-20191123064959-2c17d62f5098 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/inconshreveable/go-update v0.0.0-20160112193335-8152e7eb6ccf // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/kyokomi/emoji/v2 v2.2.8 // indirect
	github.com/lucasb-eyer/go-colorful v1.0.3 // indirect
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.11 // indirect
	github.com/mattn/go-runewidth v0.0.12 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/rivo/uniseg v0.1.0 // indirect
	github.com/ulikunitz/xz v0.5.9 // indirect
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/image v0.0.0-20191206065243-da761ea9ff43 // indirect
	golang.org/x/net v0.0.0-20220624214902-1bab6f366d9e // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
