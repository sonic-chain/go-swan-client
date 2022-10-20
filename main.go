package main

import (
	"fmt"
	"github.com/filswan/go-swan-client/command"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
	"os"
	"strings"
)

func main() {
	app := &cli.App{
		Name:                 "swan-client",
		Usage:                "A PiB level data onboarding tool for Filecoin Network",
		Version:              command.VERSION,
		EnableBashCompletion: true,
		After: func(context *cli.Context) error {
			if r := recover(); r != nil {
				panic(r)
			}
			return nil
		},
		Commands:        []*cli.Command{toolsCmd, uploadCmd, taskCmd, dealCmd, autoCmd},
		HideHelpCommand: true,
	}
	if err := app.Run(os.Args); err != nil {
		var phe *PrintHelpErr
		fmt.Fprintf(os.Stderr, "ERROR: %s\n\n", err)
		if xerrors.As(err, &phe) {
			_ = cli.ShowCommandHelp(phe.Ctx, phe.Ctx.Command.Name)
		}
		os.Exit(1)
	}
}

var uploadCmd = &cli.Command{
	Name:      "upload",
	Usage:     "Upload CAR file to ipfs server",
	ArgsUsage: "[inputPath]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "input-dir",
			Aliases: []string{"i"},
			Usage:   "directory where source files are in.",
		},
	},
	Action: func(ctx *cli.Context) error {
		inputDir := ctx.String("input-dir")
		if inputDir == "" {
			return errors.New("input-dir is required")
		}
		_, err := command.UploadCarFilesByConfig(inputDir)
		if err != nil {
			logs.GetLogger().Error(err)
			return err
		}
		return nil
	},
}

var taskCmd = &cli.Command{
	Name:  "task",
	Usage: "Send task to swan",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "name",
			Usage: "task name",
		},
		&cli.StringFlag{
			Name:    "input-dir",
			Aliases: []string{"i"},
			Usage:   "absolute path where the json or csv format source files",
		},
		&cli.StringFlag{
			Name:    "out-dir",
			Aliases: []string{"o"},
			Usage:   "directory where target files will in",
			Value:   "/tmp/tasks",
		},
		&cli.BoolFlag{
			Name:  "auto",
			Usage: "automatically send the deal after the task is created",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "manual",
			Usage: "manually send the deal after the task is created",
			Value: false,
		},
		&cli.StringFlag{
			Name:  "miner",
			Usage: "target miner ID",
		},
		&cli.StringFlag{
			Name:  "dataset",
			Usage: "curated dataset",
		},
		&cli.StringFlag{
			Name:    "description",
			Aliases: []string{"d"},
			Usage:   "task description",
		},
		&cli.IntFlag{
			Name:    "max-copy-number",
			Aliases: []string{"max"},
			Usage:   "max copy number you want to send",
			Value:   8,
		},
	},
	Action: func(ctx *cli.Context) error {
		inputDir := ctx.String("input-dir")
		if inputDir == "" {
			return errors.New("input-dir is required")
		}
		if !strings.HasSuffix(inputDir, "csv") && !strings.HasSuffix(inputDir, "json") {
			return errors.New("inputDir must be json or csv format file")
		}
		logs.GetLogger().Info("your input source file as: ", inputDir)

		auto := ctx.Bool("auto")
		manual := ctx.Bool("manual")
		minerId := ctx.String("miner")

		if auto && minerId != "" {
			return errors.New("miner cannot set when auto value is true")
		}

		if manual && minerId != "" {
			return errors.New("miner cannot set when manual value is true")
		}

		if !auto && !manual && minerId == "" {
			return errors.New("auto, manual, miner have at least one setting value")
		}

		if auto && manual {
			return errors.New("auto and manual cannot be set at the same time")
		}

		outputDir := ctx.String("out-dir")
		_, fileDesc, _, total, err := command.CreateTaskByConfig(inputDir, &outputDir, ctx.String("name"), ctx.String("miner"),
			ctx.String("dataset"), ctx.String("description"), auto, manual, ctx.Int("max-copy-number"))
		if err != nil {
			logs.GetLogger().Error(err)
			return err
		}

		if auto {
			taskId := fileDesc[0].Uuid
			exitCh := make(chan interface{})
			go func() {
				defer func() {
					exitCh <- struct{}{}
				}()
				command.GetCmdAutoDeal(&outputDir).SendAutoBidDealsBySwanClientSourceId(inputDir, taskId, total)
			}()
			<-exitCh
		}
		return nil
	},
}

var dealCmd = &cli.Command{
	Name:  "deal",
	Usage: "Send manual bid deal",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "csv",
			Usage: "the CSV file path of deal metadata",
		},
		&cli.StringFlag{
			Name:  "json",
			Usage: "the JSON file path of deal metadata",
		},
		&cli.StringFlag{
			Name:    "out-dir",
			Aliases: []string{"o"},
			Usage:   "directory where target files will in",
			Value:   "/tmp/tasks",
		},
		&cli.StringFlag{
			Name:  "miner",
			Usage: "target miner ID",
		},
	},
	Action: func(ctx *cli.Context) error {
		metadataJsonPath := ctx.String("json")
		metadataCsvPath := ctx.String("csv")

		if len(metadataJsonPath) == 0 && len(metadataCsvPath) == 0 {
			return errors.New("both metadataJsonPath and metadataCsvPath is nil")
		}

		if len(metadataJsonPath) > 0 && len(metadataCsvPath) > 0 {
			return errors.New("metadata file path is required, it cannot contain csv file path  or json file path  at the same time")
		}

		if len(metadataJsonPath) > 0 {
			logs.GetLogger().Info("Metadata json file:", metadataJsonPath)
		}
		if len(metadataCsvPath) > 0 {
			logs.GetLogger().Info("Metadata csv file:", metadataCsvPath)
		}

		_, err := command.SendDealsByConfig(ctx.String("out-dir"), ctx.String("miner"), metadataJsonPath, metadataCsvPath)
		if err != nil {
			logs.GetLogger().Error(err)
			return err
		}
		return nil
	},
}

var autoCmd = &cli.Command{
	Name:      "auto",
	Usage:     "Auto send bid deal",
	ArgsUsage: "[inputPath]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "out-dir",
			Aliases: []string{"o"},
			Usage:   "directory where target files will in.",
		},
	},
	Action: func(ctx *cli.Context) error {
		err := command.SendAutoBidDealsLoopByConfig(ctx.String("out-dir"))
		if err != nil {
			logs.GetLogger().Error(err)
			return err
		}
		return nil
	},
	Hidden: true,
}

var inPutFlag = cli.StringFlag{
	Name:    "input-dir",
	Aliases: []string{"i"},
	Usage:   "directory where source file(s) is(are) in.",
}
var importFlag = cli.BoolFlag{
	Name:  "import",
	Usage: "Whether to import lotus (default: false)",
	Value: false,
}

var outPutFlag = cli.StringFlag{
	Name:    "out-dir",
	Aliases: []string{"o"},
	Usage:   "directory where CAR file(s) will be generated.",
	Value:   "/tmp/tasks",
}

var lotusCarCmd = &cli.Command{
	Name:      "lotus",
	Usage:     "Use lotus api to generate CAR file",
	ArgsUsage: "[inputPath]",
	Flags: []cli.Flag{
		&inPutFlag,
		&outPutFlag,
		&importFlag,
	},
	Action: func(ctx *cli.Context) error {
		inputDir := ctx.String("input-dir")
		if inputDir == "" {
			return errors.New("input-dir is required")
		}
		outputDir := ctx.String("out-dir")
		if _, err := command.CreateCarFilesByConfig(inputDir, &outputDir); err != nil {
			logs.GetLogger().Error(err)
			return err
		}
		return nil
	},
}

var splitCarCmd = &cli.Command{
	Name:            "graphsplit",
	Usage:           "Use go-graphsplit tools",
	Subcommands:     []*cli.Command{generateCarCmd, carRestoreCmd},
	HideHelpCommand: true,
}

var generateCarCmd = &cli.Command{
	Name:      "car",
	Usage:     "Generate CAR files of the specified size",
	ArgsUsage: "[inputPath]",
	Flags: []cli.Flag{
		&inPutFlag,
		&outPutFlag,
		&importFlag,
		&cli.IntFlag{
			Name:  "parallel",
			Usage: "number goroutines run when building ipld nodes",
			Value: 5,
		},
		&cli.Int64Flag{
			Name:    "slice-size",
			Aliases: []string{"size"},
			Usage:   "GiB of each piece (default: 16GiB)",
			Value:   17179869184,
		},
		&cli.BoolFlag{
			Name:  "parent-path",
			Usage: "specify graph parent path",
			Value: false,
		},
	},
	Action: func(ctx *cli.Context) error {
		inputDir := ctx.String("input-dir")
		if inputDir == "" {
			return errors.New("input-dir is required")
		}
		outputDir := ctx.String("out-dir")
		if _, err := command.CreateGoCarFilesByConfig(inputDir, &outputDir, ctx.Int("parallel"), ctx.Int64("slice-size"), ctx.Bool("parent-path")); err != nil {
			logs.GetLogger().Error(err)
			return err
		}
		return nil
	},
}

var carRestoreCmd = &cli.Command{
	Name:      "restore",
	Usage:     "Restore files from CAR files",
	ArgsUsage: "[inputPath]",
	Flags: []cli.Flag{
		&outPutFlag,
		&cli.StringFlag{
			Name:     "input-dir",
			Aliases:  []string{"i"},
			Usage:    "specify source CAR path, directory or file",
			Required: true,
		},
		&cli.Int64Flag{
			Name:  "parallel",
			Usage: "number goroutines run when building ipld nodes",
			Value: 5,
		},
	},
	Action: func(ctx *cli.Context) error {
		inputDir := ctx.String("input-dir")
		if inputDir == "" {
			return errors.New("input-dir is required")
		}
		outputDir := ctx.String("out-dir")
		if err := command.RestoreCarFilesByConfig(inputDir, &outputDir, ctx.Int("parallel")); err != nil {
			logs.GetLogger().Error(err)
			return err
		}
		return nil
	},
}

var ipfsCarCmd = &cli.Command{
	Name:      "ipfs",
	Usage:     "Use ipfs api to generate CAR file",
	ArgsUsage: "[inputPath]",
	Flags: []cli.Flag{
		&inPutFlag,
		&outPutFlag,
		&importFlag,
	},
	Action: func(ctx *cli.Context) error {
		inputDir := ctx.String("input-dir")
		if inputDir == "" {
			return errors.New("input-dir is required")
		}
		outputDir := ctx.String("out-dir")
		if _, err := command.CreateIpfsCarFilesByConfig(inputDir, &outputDir); err != nil {
			logs.GetLogger().Error(err)
			return err
		}
		return nil
	},
}

var toolsCmd = &cli.Command{
	Name:            "generate-car",
	Usage:           "Generate CAR files from a file or directory",
	Subcommands:     []*cli.Command{splitCarCmd, lotusCarCmd, ipfsCarCmd, ipfsCmdCarCmd},
	HideHelpCommand: true,
}

var ipfsCmdCarCmd = &cli.Command{
	Name:      "ipfs-car",
	Usage:     "use the ipfs-car command to generate the CAR file",
	ArgsUsage: "[inputPath]",
	Flags: []cli.Flag{
		&inPutFlag,
		&outPutFlag,
		&importFlag,
	},
	Action: func(ctx *cli.Context) error {
		inputDir := ctx.String("input-dir")
		if inputDir == "" {
			return errors.New("input-dir is required")
		}
		outputDir := ctx.String("out-dir")
		if _, err := command.CreateIpfsCmdCarFilesByConfig(inputDir, &outputDir); err != nil {
			logs.GetLogger().Error(err)
			return err
		}
		return nil
	},
}

type PrintHelpErr struct {
	Err error
	Ctx *cli.Context
}

func (e *PrintHelpErr) Error() string {
	return e.Err.Error()
}

func (e *PrintHelpErr) Unwrap() error {
	return e.Err
}

func (e *PrintHelpErr) Is(o error) bool {
	_, ok := o.(*PrintHelpErr)
	return ok
}

func ShowHelp(cctx *cli.Context, err error) error {
	return &PrintHelpErr{Err: err, Ctx: cctx}
}

func WithCategory(cat string, cmd *cli.Command) *cli.Command {
	cmd.Category = strings.ToUpper(cat)
	return cmd
}
