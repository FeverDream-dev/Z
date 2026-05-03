package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/FeverDream-dev/zsistant/internal/assistant"
	cfgpkg "github.com/FeverDream-dev/zsistant/internal/config"
	"github.com/FeverDream-dev/zsistant/internal/server"
)

var zaziVersion = "v0.0.1-dev"
var zaziCommit = "unknown"

func main() {
	code := run(os.Args[1:])
	if code != 0 {
		os.Exit(code)
	}
}

func run(args []string) int {
	if len(args) == 0 {
		printUsage()
		return 0
	}
	switch args[0] {
	case "assistant":
		return runAssistant(args[1:])
	case "serve", "server", "daemon":
		return runServe(args[1:])
	case "version":
		fmt.Printf("zazi version %s (commit %s)\n", zaziVersion, zaziCommit)
		return 0
	case "doctor":
		return runDoctor()
	case "init":
		if err := runInit(); err != nil {
			fmt.Fprintf(os.Stderr, "init error: %v\n", err)
			return 1
		}
		return 0
	case "--help", "-h":
		printUsage()
		return 0
	default:
		printUsage()
		return 1
	}
}

func printUsage() {
	fmt.Println("zazi CLI - Zsistant assistant manager")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  zazi assistant create \u003cid\u003e --name=\u003cname\u003e            Create a new assistant")
	fmt.Println("  zazi assistant list                                   List all assistants")
	fmt.Println("  zazi assistant show \u003cid\u003e                              Show assistant details")
	fmt.Println("  zazi assistant delete \u003cid\u003e                            Delete an assistant")
	fmt.Println("  zazi serve [--addr=:8080] [--base=~/.zazi]            Start web server")
	fmt.Println("  zazi version                                          Print version")
	fmt.Println("  zazi doctor                                           Run diagnostics")
	fmt.Println("  zazi init                                             Initialize config and directories")
	fmt.Println("  zazi --help                                           Show this help")
}

func runAssistant(args []string) int {
	if len(args) == 0 {
		fmt.Println("Assistant subcommands: create, list, show, delete")
		return 1
	}

	switch args[0] {
	case "create":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "usage: zazi assistant create \u003cid\u003e --name=\u003cname\u003e")
			return 1
		}
		id := args[1]
		var name string
		for i := 2; i < len(args); i++ {
			if len(args[i]) > 7 && args[i][:7] == "--name=" {
				name = args[i][7:]
			}
		}
		if name == "" {
			name = id
		}
		reg := getReg()
		a, err := reg.Create(id, name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create assistant: %v\n", err)
			return 1
		}
		b, _ := json.MarshalIndent(a, "  ", "  ")
		fmt.Println(string(b))
		return 0

	case "list":
		reg := getReg()
		list, err := reg.List()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to list assistants: %v\n", err)
			return 1
		}
		for _, a := range list {
			fmt.Printf("- %s (%s) status=%s\n", a.ID, a.Name, a.Status)
		}
		return 0

	case "show":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "usage: zazi assistant show \u003cid\u003e")
			return 1
		}
		reg := getReg()
		a, err := reg.Get(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to get assistant: %v\n", err)
			return 1
		}
		b, _ := json.MarshalIndent(a, "  ", "  ")
		fmt.Println(string(b))
		return 0

	case "delete":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "usage: zazi assistant delete \u003cid\u003e")
			return 1
		}
		reg := getReg()
		if err := reg.Delete(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "failed to delete: %v\n", err)
			return 1
		}
		fmt.Println("assistant deleted:", args[1])
		return 0

	default:
		fmt.Println("Unknown assistant subcommand:", args[0])
		fmt.Println("Available: create, list, show, delete")
		return 1
	}
}

func runServe(args []string) int {
	addr := ":8080"
	basePath := ""
	for i := range args {
		if len(args[i]) > 7 && args[i][:7] == "--addr=" {
			addr = args[i][7:]
		}
		if len(args[i]) > 7 && args[i][:7] == "--base=" {
			basePath = args[i][7:]
		}
	}
	if basePath == "" {
		cfgPath, _ := cfgpkg.DefaultPath()
		cfg, _ := cfgpkg.Load(cfgPath)
		if cfg == nil {
			cfg = &cfgpkg.Config{}
			cfg.DataPath = "~/.zazi"
		}
		basePath = cfgpkg.ExpandPath(cfg.DataPath)
	}
	srv := server.New(addr, basePath)
	fmt.Printf("Starting Zsistant server on http://%s\n", addr)
	fmt.Printf("Data path: %s\n", basePath)
	if err := srv.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		return 1
	}
	return 0
}

func runInit() error {
	path, err := cfgpkg.DefaultPath()
	if err != nil {
		return err
	}
	if _, statErr := os.Stat(path); statErr == nil {
		fmt.Printf("Config already exists at %s. Skipping write.\n", path)
		return nil
	}
	cfg := cfgpkg.DefaultConfig()
	if err := cfgpkg.Save(path, &cfg); err != nil {
		return fmt.Errorf("saving default config: %w", err)
	}
	if err := cfgpkg.EnsureDirs(&cfg); err != nil {
		return fmt.Errorf("ensuring data directories: %w", err)
	}
	fmt.Printf("Initialized Zazi at %s\n", path)
	return nil
}

func runDoctor() int {
	path, _ := cfgpkg.DefaultPath()
	c, err := cfgpkg.Load(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		return 1
	}

	if _, statErr := os.Stat(path); statErr != nil {
		fmt.Printf("Config file: %s - MISSING\n", path)
	} else {
		fmt.Printf("Config file: %s - OK\n", path)
	}

	dataPath := cfgpkg.ExpandPath(c.DataPath)
	writable := true
	if st, err := os.Stat(dataPath); err != nil || !st.IsDir() {
		writable = false
	} else {
		testPath := filepath.Join(dataPath, ".doctor_test")
		if err := os.WriteFile(testPath, []byte("ok"), 0644); err != nil {
			writable = false
		} else {
			_ = os.Remove(testPath)
		}
	}

	if writable {
		fmt.Printf("Data path: %s - writable\n", dataPath)
	} else {
		fmt.Printf("Data path: %s - not writable\n", dataPath)
	}

	fmt.Printf("Config: DataPath=%q, LogLevel=%q, ServerPort=%d, Providers=%v\n",
		c.DataPath, c.LogLevel, c.ServerPort, c.Providers)

	if !writable {
		fmt.Println("Hint: Ensure the data directory exists and is writable, or run 'zazi init'.")
	}
	return 0
}

func getReg() *assistant.Registry {
	cfgPath, _ := cfgpkg.DefaultPath()
	cfg, _ := cfgpkg.Load(cfgPath)
	if cfg == nil {
		cfg = &cfgpkg.Config{}
		cfg.DataPath = "~/.zazi"
	}
	path := cfgpkg.ExpandPath(cfg.DataPath)
	return assistant.NewRegistry(path)
}

func homeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if home == "" {
		return "", errors.New("empty home directory")
	}
	return home, nil
}
