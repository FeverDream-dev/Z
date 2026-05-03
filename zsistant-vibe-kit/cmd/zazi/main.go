package main

import (
  "context"
  "encoding/json"
  "errors"
  "fmt"
  "os"
  "os/signal"
  "path/filepath"
  "strings"
  "time"

  "syscall"

  zkconfig "github.com/FeverDream-dev/zsistant/internal/config"
  az "github.com/FeverDream-dev/zsistant/internal/agents"
  "github.com/FeverDream-dev/zsistant/internal/bus"
  "github.com/FeverDream-dev/zsistant/internal/channels"
  "github.com/FeverDream-dev/zsistant/internal/skills"
  "github.com/FeverDream-dev/zsistant/internal/jobs"
  "github.com/FeverDream-dev/zsistant/internal/llm"
  "github.com/FeverDream-dev/zsistant/internal/server"
  "github.com/FeverDream-dev/zsistant/internal/trainer"
)

// Build-time versioning. Defaults can be overridden with -ldflags.
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
  case "agent":
    return runAgent(args[1:])
  case "--help", "-h":
    printUsage()
    return 0
  case "version":
    fmt.Printf("zazi version %s (commit %s)\n", zaziVersion, zaziCommit)
    return 0
  case "doctor":
    runDoctor()
    return 0
  case "init":
    if err := runInit(); err != nil {
      fmt.Fprintf(os.Stderr, "init error: %v\n", err)
      return 1
    }
    return 0
  case "chat":
    return runChat(args[1:])
  case "telegram":
    return runTelegram(args[1:])
  case "discord":
    return runDiscord(args[1:])
  case "whatsapp":
    return runWhatsApp(args[1:])
  case "channel":
    return runChannel(args[1:])
  case "job":
    return runJob(args[1:])
  case "train":
    return runTrain(args[1:])
  case "persona":
    return runTrain(args[1:])
  case "acl":
    return runACL(args[1:])
  case "skill":
    return runSkill(args[1:])
  case "validate":
    return runValidate(args[1:])
  case "release":
    return runRelease(args[1:])
  case "provider":
    return runProvider(args[1:])
  case "serve", "daemon":
    home, err := homeDir()
    if err != nil {
      fmt.Fprintln(os.Stderr, "can't determine home:", err)
      return 1
    }
    base := filepath.Join(home, ".zazi")
    // Load config for port, default to 8080
    cfg, err := zkconfig.Load("")
    if err != nil {
      cfg = zkconfig.NewDefault()
    }
    addr := fmt.Sprintf("0.0.0.0:%d", cfg.ServerPort)
    srv := server.New(addr, base)
    if err := srv.Run(); err != nil {
      fmt.Fprintf(os.Stderr, "server error: %v\n", err)
      return 1
    }
    return 0
  default:
    printUsage()
    return 1
  }
}

func printUsage() {
  fmt.Fprintln(os.Stdout, "zazi CLI - a minimal Zsistant helper")
  fmt.Fprintln(os.Stdout, "")
  fmt.Fprintln(os.Stdout, "Usage:")
  fmt.Fprintln(os.Stdout, "  zazi version           Print version information")
  fmt.Fprintln(os.Stdout, "  zazi doctor            Run diagnostics for configuration")
  fmt.Fprintln(os.Stdout, "  zazi init              Initialize user config and data directories")
  fmt.Fprintln(os.Stdout, "  zazi serve             Run web server")
  fmt.Fprintln(os.Stdout, "  zazi chat              Chat with an agent")
  fmt.Fprintln(os.Stdout, "  zazi telegram          Telegram adapter commands")
  fmt.Fprintln(os.Stdout, "  zazi discord           Discord adapter commands")
  fmt.Fprintln(os.Stdout, "  zazi whatsapp          WhatsApp adapter commands")
  fmt.Fprintln(os.Stdout, "  zazi channel           Channel adapter commands (alias for telegram/discord/whatsapp)")
  fmt.Fprintln(os.Stdout, "  zazi job               Job queue commands")
  fmt.Fprintln(os.Stdout, "  zazi train             Persona trainer commands")
  fmt.Fprintln(os.Stdout, "  zazi persona           Persona trainer commands (alias)")
  fmt.Fprintln(os.Stdout, "  zazi acl               Inter-agent ACL commands")
  fmt.Fprintln(os.Stdout, "  zazi skill             Skill analyzer commands")
  fmt.Fprintln(os.Stdout, "  zazi validate          UI validation commands")
  fmt.Fprintln(os.Stdout, "  zazi release           Release management commands")
  fmt.Fprintln(os.Stdout, "  zazi provider          LLM provider catalog commands")
  fmt.Fprintln(os.Stdout, "  zazi --help, -h        Show this help")
  fmt.Fprintln(os.Stdout, "")
  fmt.Fprintln(os.Stdout, "Examples:")
  fmt.Fprintln(os.Stdout, "  zazi version")
  fmt.Fprintln(os.Stdout, "  zazi init")
}

func runAgent(args []string) int {
  // Determine base registry path under HOME/.zazi
  home, err := homeDir()
  if err != nil {
    fmt.Fprintln(os.Stderr, "can't determine home:", err)
    return 1
  }
  base := filepath.Join(home, ".zazi")
  reg := az.New(base)
  if len(args) == 0 {
    printAgentUsage()
    return 0
  }
  switch args[0] {
  case "create":
    if len(args) < 3 {
      fmt.Fprintln(os.Stderr, "usage: zazi agent create <id> --name=<name> --role=<role>")
      return 1
    }
    id := args[1]
    var name, role string
    for _, a := range args[2:] {
      if strings.HasPrefix(a, "--name=") {
        name = strings.TrimPrefix(a, "--name=")
      } else if strings.HasPrefix(a, "--role=") {
        role = strings.TrimPrefix(a, "--role=")
      }
    }
    if name == "" || role == "" {
      fmt.Fprintln(os.Stderr, "missing --name or --role parameter")
      return 1
    }
    if _, err := reg.Create(id, name, role); err != nil {
      fmt.Fprintln(os.Stderr, "failed to create agent:", err)
      return 1
    }
    fmt.Println("agent created:", id)
    return 0
  case "list":
    list, err := reg.List()
    if err != nil {
      fmt.Fprintln(os.Stderr, "failed to list agents:", err)
      return 1
    }
    for _, a := range list {
      b, _ := json.MarshalIndent(a, "  ", "  ")
      fmt.Println(string(b))
    }
    return 0
  case "show":
    if len(args) < 2 {
      fmt.Fprintln(os.Stderr, "usage: zazi agent show <id>")
      return 1
    }
    id := args[1]
    a, err := reg.Get(id)
    if err != nil {
      fmt.Fprintln(os.Stderr, "failed to get agent:", err)
      return 1
    }
    b, _ := json.MarshalIndent(a, "  ", "  ")
    fmt.Println(string(b))
    return 0
  case "delete":
    if len(args) < 2 {
      fmt.Fprintln(os.Stderr, "usage: zazi agent delete <id>")
      return 1
    }
    id := args[1]
    if err := reg.Delete(id); err != nil {
      fmt.Fprintln(os.Stderr, "failed to delete agent:", err)
      return 1
    }
    fmt.Println("agent deleted:", id)
    return 0
  default:
    printAgentUsage()
    return 0
  }
}

func printAgentUsage() {
  fmt.Fprintln(os.Stdout, "Agent subcommands:")
  fmt.Fprintln(os.Stdout, "  zazi agent create <id> --name=<name> --role=<role>")
  fmt.Fprintln(os.Stdout, "  zazi agent list")
  fmt.Fprintln(os.Stdout, "  zazi agent show <id>")
  fmt.Fprintln(os.Stdout, "  zazi agent delete <id>")
}

func runChannel(args []string) int {
  if len(args) == 0 {
    printChannelUsage()
    return 0
  }

  channelArgs := args[1:]
  switch args[0] {
  case "telegram":
    return runTelegram(channelArgs)
  case "discord":
    return runDiscord(channelArgs)
  case "whatsapp":
    return runWhatsApp(channelArgs)
  default:
    printChannelUsage()
    return 1
  }
}

func printChannelUsage() {
  fmt.Fprintln(os.Stdout, "Channel adapter commands:")
  fmt.Fprintln(os.Stdout, "  zazi channel telegram setup --agent=<id> [--token=<token>]")
  fmt.Fprintln(os.Stdout, "  zazi channel telegram test --agent=<id> --chat=<chat_id> --message=<msg>")
  fmt.Fprintln(os.Stdout, "  zazi channel telegram send --agent=<id> --chat=<chat_id> --message=<msg>")
  fmt.Fprintln(os.Stdout, "  zazi channel telegram listen --agent=<id>")
  fmt.Fprintln(os.Stdout, "  zazi channel discord setup --agent=<id> [--token=<token>]")
  fmt.Fprintln(os.Stdout, "  zazi channel discord test --agent=<id> --channel=<channel_id> --message=<msg>")
  fmt.Fprintln(os.Stdout, "  zazi channel whatsapp setup --agent=<id> [--phone=<phone_id>] [--token=<token>] [--verify=<verify_token>]")
  fmt.Fprintln(os.Stdout, "  zazi channel whatsapp test --agent=<id> --phone=<phone_number> --message=<msg>")
  fmt.Fprintln(os.Stdout, "  zazi channel whatsapp verify --mode=<mode> --token=<token> --challenge=<challenge>")
}

func runTelegram(args []string) int {
  if len(args) == 0 {
    printTelegramUsage()
    return 0
  }
  home, err := homeDir()
  if err != nil {
    fmt.Fprintln(os.Stderr, "can't determine home:", err)
    return 1
  }
  base := filepath.Join(home, ".zazi")

  switch args[0] {
  case "setup":
    return runTelegramSetup(args[1:], base)
  case "test":
    return runTelegramTest(args[1:], base)
  case "send":
    return runTelegramSend(args[1:], base)
  case "listen":
    return runTelegramListen(args[1:], base)
  default:
    printTelegramUsage()
    return 0
  }
}

func printTelegramUsage() {
  fmt.Fprintln(os.Stdout, "Telegram subcommands:")
  fmt.Fprintln(os.Stdout, "  zazi telegram setup --agent=<id> [--token=<token>]")
  fmt.Fprintln(os.Stdout, "  zazi telegram test --agent=<id> --chat=<chat_id> --message=<msg>")
  fmt.Fprintln(os.Stdout, "  zazi telegram send --agent=<id> --chat=<chat_id> --message=<msg>")
  fmt.Fprintln(os.Stdout, "  zazi telegram listen --agent=<id>")
}

func runDiscord(args []string) int {
  if len(args) == 0 {
    printDiscordUsage()
    return 0
  }
  home, err := homeDir()
  if err != nil {
    fmt.Fprintln(os.Stderr, "can't determine home:", err)
    return 1
  }
  base := filepath.Join(home, ".zazi")

  switch args[0] {
  case "setup":
    return runDiscordSetup(args[1:], base)
  case "test":
    return runDiscordTest(args[1:], base)
  default:
    printDiscordUsage()
    return 0
  }
}

func printDiscordUsage() {
  fmt.Fprintln(os.Stdout, "Discord subcommands:")
  fmt.Fprintln(os.Stdout, "  zazi discord setup --agent=<id> [--token=<token>]")
  fmt.Fprintln(os.Stdout, "  zazi discord test --agent=<id> --channel=<channel_id> --message=<msg>")
}

func runWhatsApp(args []string) int {
  if len(args) == 0 {
    printWhatsAppUsage()
    return 0
  }
  home, err := homeDir()
  if err != nil {
    fmt.Fprintln(os.Stderr, "can't determine home:", err)
    return 1
  }
  base := filepath.Join(home, ".zazi")

  switch args[0] {
  case "setup":
    return runWhatsAppSetup(args[1:], base)
  case "test":
    return runWhatsAppTest(args[1:], base)
  case "verify":
    return runWhatsAppVerify(args[1:])
  default:
    printWhatsAppUsage()
    return 0
  }
}

func printWhatsAppUsage() {
  fmt.Fprintln(os.Stdout, "WhatsApp subcommands:")
  fmt.Fprintln(os.Stdout, "  zazi whatsapp setup --agent=<id> [--phone=<phone_id>] [--token=<token>] [--verify=<verify_token>]")
  fmt.Fprintln(os.Stdout, "  zazi whatsapp test --agent=<id> --phone=<phone_number> --message=<msg>")
  fmt.Fprintln(os.Stdout, "  zazi whatsapp verify --mode=<mode> --token=<token> --challenge=<challenge>")
}

func runJob(args []string) int {
  if len(args) == 0 {
    printJobUsage()
    return 0
  }
  home, err := homeDir()
  if err != nil {
    fmt.Fprintln(os.Stderr, "can't determine home:", err)
    return 1
  }
  base := filepath.Join(home, ".zazi")

  switch args[0] {
  case "list":
    return runJobList(args[1:], base)
  case "pause":
    return runJobPause(args[1:], base)
  case "resume":
    return runJobResume(args[1:], base)
  case "cancel":
    return runJobCancel(args[1:], base)
  default:
    printJobUsage()
    return 0
  }
}

func printJobUsage() {
  fmt.Fprintln(os.Stdout, "Job subcommands:")
  fmt.Fprintln(os.Stdout, "  zazi job list --agent=<id>")
  fmt.Fprintln(os.Stdout, "  zazi job pause --agent=<id> --job=<job_id> [--reason=<reason>]")
  fmt.Fprintln(os.Stdout, "  zazi job resume --agent=<id> --job=<job_id>")
  fmt.Fprintln(os.Stdout, "  zazi job cancel --agent=<id> --job=<job_id>")
}

func runJobList(args []string, base string) int {
  var agentID string
  for _, a := range args {
    if strings.HasPrefix(a, "--agent=") {
      agentID = strings.TrimPrefix(a, "--agent=")
    }
  }
  if agentID == "" {
    fmt.Fprintln(os.Stderr, "usage: zazi job list --agent=<id>")
    return 1
  }

  reg := az.New(base)
  a, err := reg.Get(agentID)
  if err != nil {
    fmt.Fprintf(os.Stderr, "agent not found: %v\n", err)
    return 1
  }

  q := jobs.NewQueue(a.WorkspaceRoot)
  list, err := q.List()
  if err != nil {
    fmt.Fprintf(os.Stderr, "list jobs: %v\n", err)
    return 1
  }

  if len(list) == 0 {
    fmt.Println("No jobs found.")
    return 0
  }

  fmt.Printf("Jobs for agent %s:\n", agentID)
  for _, j := range list {
    status := j.Status
    if j.PausedReason != "" {
      status += " (" + j.PausedReason + ")"
    }
    fmt.Printf("  %s | %s | %s | retries: %d/%d\n", j.ID, j.Objective, status, j.RetryCount, j.MaxRetries)
  }
  return 0
}

func runJobPause(args []string, base string) int {
  var agentID, jobID, reason string
  for _, a := range args {
    if strings.HasPrefix(a, "--agent=") {
      agentID = strings.TrimPrefix(a, "--agent=")
    } else if strings.HasPrefix(a, "--job=") {
      jobID = strings.TrimPrefix(a, "--job=")
    } else if strings.HasPrefix(a, "--reason=") {
      reason = strings.TrimPrefix(a, "--reason=")
    }
  }
  if agentID == "" || jobID == "" {
    fmt.Fprintln(os.Stderr, "usage: zazi job pause --agent=<id> --job=<job_id> [--reason=<reason>]")
    return 1
  }
  if reason == "" {
    reason = "manual pause"
  }

  reg := az.New(base)
  a, err := reg.Get(agentID)
  if err != nil {
    fmt.Fprintf(os.Stderr, "agent not found: %v\n", err)
    return 1
  }

  q := jobs.NewQueue(a.WorkspaceRoot)
  paused, err := q.Pause(jobID, reason)
  if err != nil {
    fmt.Fprintf(os.Stderr, "pause job: %v\n", err)
    return 1
  }
  fmt.Printf("Job %s paused: %s\n", paused.ID, paused.PausedReason)
  return 0
}

func runJobResume(args []string, base string) int {
  var agentID, jobID string
  for _, a := range args {
    if strings.HasPrefix(a, "--agent=") {
      agentID = strings.TrimPrefix(a, "--agent=")
    } else if strings.HasPrefix(a, "--job=") {
      jobID = strings.TrimPrefix(a, "--job=")
    }
  }
  if agentID == "" || jobID == "" {
    fmt.Fprintln(os.Stderr, "usage: zazi job resume --agent=<id> --job=<job_id>")
    return 1
  }

  reg := az.New(base)
  a, err := reg.Get(agentID)
  if err != nil {
    fmt.Fprintf(os.Stderr, "agent not found: %v\n", err)
    return 1
  }

  q := jobs.NewQueue(a.WorkspaceRoot)
  resumed, err := q.Resume(jobID)
  if err != nil {
    fmt.Fprintf(os.Stderr, "resume job: %v\n", err)
    return 1
  }
  fmt.Printf("Job %s resumed (status: %s)\n", resumed.ID, resumed.Status)
  return 0
}

func runRelease(args []string) int {
  if len(args) == 0 {
    printReleaseUsage()
    return 0
  }

  switch args[0] {
  case "prepare":
    return runReleasePrepare(args[1:])
  case "publish":
    return runReleasePublish(args[1:])
  default:
    printReleaseUsage()
    return 0
  }
}

func printReleaseUsage() {
  fmt.Fprintln(os.Stdout, "Release subcommands:")
  fmt.Fprintln(os.Stdout, "  zazi release prepare --version=vX.Y.Z")
  fmt.Fprintln(os.Stdout, "  zazi release publish --version=vX.Y.Z --approve=<approval text>")
}

func runReleasePrepare(args []string) int {
  var version string
  for _, a := range args {
    if strings.HasPrefix(a, "--version=") {
      version = strings.TrimPrefix(a, "--version=")
    }
  }
  if version == "" {
    fmt.Fprintln(os.Stderr, "usage: zazi release prepare --version=vX.Y.Z")
    return 1
  }

  fmt.Printf("=== Release Preparation: %s ===\n\n", version)

  // Run tests
  fmt.Println("Running tests...")
  // We can't exec go test from here easily, so we document it
  fmt.Println("  [INFO] Run: go test ./...")
  fmt.Println("  [INFO] Run: go build ./...")
  fmt.Println("  [INFO] Run: zazi validate ui")
  fmt.Println()

  // Generate changelog draft
  fmt.Println("Changelog Draft:")
  fmt.Println("  ## [Unreleased]")
  fmt.Printf("  - Release %s\n", version)
  fmt.Println("  - [Insert changes here]")
  fmt.Println()

  // Show checklist
  fmt.Println("Release Checklist:")
  fmt.Println("  [ ] All tests pass")
  fmt.Println("  [ ] Build succeeds")
  fmt.Println("  [ ] No secrets in source")
  fmt.Println("  [ ] Version bumped")
  fmt.Println("  [ ] CHANGELOG updated")
  fmt.Println("  [ ] UI validation passed")
  fmt.Println("  [ ] Documentation current")
  fmt.Println()

  fmt.Println("=== Preparation Complete ===")
  fmt.Println("Nothing has been published. Review the checklist and run:")
  fmt.Printf("  zazi release publish --version=%s --approve=\"I have reviewed...\"\n", version)
  return 0
}

func runReleasePublish(args []string) int {
  var version, approve string
  for _, a := range args {
    if strings.HasPrefix(a, "--version=") {
      version = strings.TrimPrefix(a, "--version=")
    } else if strings.HasPrefix(a, "--approve=") {
      approve = strings.TrimPrefix(a, "--approve=")
    }
  }
  if version == "" || approve == "" {
    fmt.Fprintln(os.Stderr, "usage: zazi release publish --version=vX.Y.Z --approve=<approval text>")
    return 1
  }

  if len(approve) < 10 {
    fmt.Fprintln(os.Stderr, "Approval text must be at least 10 characters")
    return 1
  }

  fmt.Printf("=== Release Publish: %s ===\n\n", version)
  fmt.Printf("Approval recorded: %q\n", approve)
  fmt.Println()
  fmt.Println("[DRY-RUN] Would execute:")
  fmt.Printf("  git tag %s\n", version)
  fmt.Printf("  git push origin %s\n", version)
  fmt.Println("  Create GitHub release with generated notes")
  fmt.Println()
  fmt.Println("=== Release workflow gated ===")
  fmt.Println("No actual push/tag was performed. Complete these steps manually after final review.")
  return 0
}

func runValidate(args []string) int {
  if len(args) == 0 {
    printValidateUsage()
    return 0
  }

  switch args[0] {
  case "ui":
    return runValidateUI(args[1:])
  default:
    printValidateUsage()
    return 0
  }
}

func printValidateUsage() {
  fmt.Fprintln(os.Stdout, "Validation subcommands:")
  fmt.Fprintln(os.Stdout, "  zazi validate ui [--url=http://localhost:8080]")
}

func runValidateUI(args []string) int {
  url := "http://localhost:8080"
  for _, a := range args {
    if strings.HasPrefix(a, "--url=") {
      url = strings.TrimPrefix(a, "--url=")
    }
  }

  home, err := homeDir()
  if err != nil {
    fmt.Fprintln(os.Stderr, "can't determine home:", err)
    return 1
  }
  base := filepath.Join(home, ".zazi")

  s := server.New("localhost:8080", base)
  results := s.Validate(url)
  fmt.Println(server.FormatReport(results))
  return 0
}

func runSkill(args []string) int {
  if len(args) == 0 {
    printSkillUsage()
    return 0
  }

  switch args[0] {
  case "analyze":
    return runSkillAnalyze(args[1:])
  default:
    printSkillUsage()
    return 0
  }
}

func printSkillUsage() {
  fmt.Fprintln(os.Stdout, "Skill subcommands:")
  fmt.Fprintln(os.Stdout, "  zazi skill analyze --path=<folder> | --text=<markdown>")
}

func runSkillAnalyze(args []string) int {
  var path, text string
  for _, a := range args {
    if strings.HasPrefix(a, "--path=") {
      path = strings.TrimPrefix(a, "--path=")
    } else if strings.HasPrefix(a, "--text=") {
      text = strings.TrimPrefix(a, "--text=")
    }
  }

  a := skills.NewAnalyzer()
  var report *skills.RiskReport
  var err error

  if path != "" {
    report, err = a.AnalyzeFolder(path)
    if err != nil {
      fmt.Fprintf(os.Stderr, "analyze failed: %v\n", err)
      return 1
    }
  } else if text != "" {
    report = a.AnalyzeText("inline-skill", text)
  } else {
    fmt.Fprintln(os.Stderr, "usage: zazi skill analyze --path=<folder> | --text=<markdown>")
    return 1
  }

  fmt.Println(report.String())
  return 0
}

func runACL(args []string) int {
  if len(args) == 0 {
    printACLUsage()
    return 0
  }
  home, err := homeDir()
  if err != nil {
    fmt.Fprintln(os.Stderr, "can't determine home:", err)
    return 1
  }
  base := filepath.Join(home, ".zazi")

  switch args[0] {
  case "allow":
    return runACLAllow(args[1:], base)
  case "revoke":
    return runACLRevoke(args[1:], base)
  case "list":
    return runACLList(args[1:], base)
  default:
    printACLUsage()
    return 0
  }
}

func printACLUsage() {
  fmt.Fprintln(os.Stdout, "ACL subcommands:")
  fmt.Fprintln(os.Stdout, "  zazi acl allow --agent=<id> --peer=<peer_id> --perm=<perm>[,<perm>...] [--scope=<path>]")
  fmt.Fprintln(os.Stdout, "  zazi acl revoke --agent=<id> --peer=<peer_id> [--perm=<perm>[,<perm>...]]")
  fmt.Fprintln(os.Stdout, "  zazi acl list --agent=<id>")
  fmt.Fprintln(os.Stdout, "Permissions: summary, file_read, file_write, exec")
}

func runACLAllow(args []string, base string) int {
  var agentID, peerID, permStr, scope string
  for _, a := range args {
    if strings.HasPrefix(a, "--agent=") {
      agentID = strings.TrimPrefix(a, "--agent=")
    } else if strings.HasPrefix(a, "--peer=") {
      peerID = strings.TrimPrefix(a, "--peer=")
    } else if strings.HasPrefix(a, "--perm=") {
      permStr = strings.TrimPrefix(a, "--perm=")
    } else if strings.HasPrefix(a, "--scope=") {
      scope = strings.TrimPrefix(a, "--scope=")
    }
  }
  if agentID == "" || peerID == "" || permStr == "" {
    fmt.Fprintln(os.Stderr, "usage: zazi acl allow --agent=<id> --peer=<peer_id> --perm=<perm>[,<perm>...] [--scope=<path>]")
    return 1
  }

  reg := az.New(base)
  if _, err := reg.Get(agentID); err != nil {
    fmt.Fprintf(os.Stderr, "agent not found: %v\n", err)
    return 1
  }

  var perms []bus.Permission
  for _, p := range strings.Split(permStr, ",") {
    perms = append(perms, bus.Permission(p))
  }

  b := bus.NewBus(base)
  if err := b.Allow(agentID, peerID, perms, scope); err != nil {
    fmt.Fprintf(os.Stderr, "allow failed: %v\n", err)
    return 1
  }
  fmt.Printf("Allowed %s to %s: %v (scope: %s)\n", peerID, agentID, perms, scope)
  return 0
}

func runACLRevoke(args []string, base string) int {
  var agentID, peerID, permStr string
  for _, a := range args {
    if strings.HasPrefix(a, "--agent=") {
      agentID = strings.TrimPrefix(a, "--agent=")
    } else if strings.HasPrefix(a, "--peer=") {
      peerID = strings.TrimPrefix(a, "--peer=")
    } else if strings.HasPrefix(a, "--perm=") {
      permStr = strings.TrimPrefix(a, "--perm=")
    }
  }
  if agentID == "" || peerID == "" {
    fmt.Fprintln(os.Stderr, "usage: zazi acl revoke --agent=<id> --peer=<peer_id> [--perm=<perm>[,<perm>...]]")
    return 1
  }

  reg := az.New(base)
  if _, err := reg.Get(agentID); err != nil {
    fmt.Fprintf(os.Stderr, "agent not found: %v\n", err)
    return 1
  }

  var perms []bus.Permission
  if permStr != "" {
    for _, p := range strings.Split(permStr, ",") {
      perms = append(perms, bus.Permission(p))
    }
  }

  b := bus.NewBus(base)
  if err := b.Revoke(agentID, peerID, perms); err != nil {
    fmt.Fprintf(os.Stderr, "revoke failed: %v\n", err)
    return 1
  }
  if len(perms) == 0 {
    fmt.Printf("Revoked all permissions for %s from %s\n", peerID, agentID)
  } else {
    fmt.Printf("Revoked %v for %s from %s\n", perms, peerID, agentID)
  }
  return 0
}

func runACLList(args []string, base string) int {
  var agentID string
  for _, a := range args {
    if strings.HasPrefix(a, "--agent=") {
      agentID = strings.TrimPrefix(a, "--agent=")
    }
  }
  if agentID == "" {
    fmt.Fprintln(os.Stderr, "usage: zazi acl list --agent=<id>")
    return 1
  }

  reg := az.New(base)
  if _, err := reg.Get(agentID); err != nil {
    fmt.Fprintf(os.Stderr, "agent not found: %v\n", err)
    return 1
  }

  b := bus.NewBus(base)
  rules, err := b.List(agentID)
  if err != nil {
    fmt.Fprintf(os.Stderr, "list failed: %v\n", err)
    return 1
  }

  if len(rules) == 0 {
    fmt.Printf("No ACL rules for agent %s (deny by default)\n", agentID)
    return 0
  }

  fmt.Printf("ACL rules for agent %s:\n", agentID)
  for _, r := range rules {
    fmt.Printf("  %s -> %v (scope: %s)\n", r.PeerID, r.Perms, r.Scope)
  }
  return 0
}

func runTrain(args []string) int {
  if len(args) == 0 {
    printTrainUsage()
    return 0
  }
  home, err := homeDir()
  if err != nil {
    fmt.Fprintln(os.Stderr, "can't determine home:", err)
    return 1
  }
  base := filepath.Join(home, ".zazi")

  switch args[0] {
  case "observe":
    return runTrainObserve(args[1:], base)
  case "propose":
    return runTrainPropose(args[1:], base)
  case "apply":
    return runTrainApply(args[1:], base)
  default:
    printTrainUsage()
    return 0
  }
}

func printTrainUsage() {
  fmt.Fprintln(os.Stdout, "Trainer subcommands:")
  fmt.Fprintln(os.Stdout, "  zazi train observe --agent=<id> --message=<msg>")
  fmt.Fprintln(os.Stdout, "  zazi train propose --agent=<id> [--threshold=<n>]")
  fmt.Fprintln(os.Stdout, "  zazi train apply --agent=<id> --patch=<text>")
}

func runTrainObserve(args []string, base string) int {
  var agentID, message string
  for _, a := range args {
    if strings.HasPrefix(a, "--agent=") {
      agentID = strings.TrimPrefix(a, "--agent=")
    } else if strings.HasPrefix(a, "--message=") {
      message = strings.TrimPrefix(a, "--message=")
    }
  }
  if agentID == "" || message == "" {
    fmt.Fprintln(os.Stderr, "usage: zazi train observe --agent=<id> --message=<msg>")
    return 1
  }

  reg := az.New(base)
  a, err := reg.Get(agentID)
  if err != nil {
    fmt.Fprintf(os.Stderr, "agent not found: %v\n", err)
    return 1
  }

  tr := trainer.NewTrainer()
  tr.Observe(message)
  tr.ObserveFeedback(message)

  // Persist profile to a simple signals file
  signalsPath := filepath.Join(a.WorkspaceRoot, "trainer_signals.txt")
  f, err := os.OpenFile(signalsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  if err != nil {
    fmt.Fprintf(os.Stderr, "failed to open signals file: %v\n", err)
    return 1
  }
  defer f.Close()

  top := tr.Profile().Top()
  for _, d := range top {
    fmt.Fprintf(f, "%s=%s (score:%d)\n", d.Name, d.Value, d.Score)
  }

  fmt.Printf("Observed %d style signals for agent %s\n", len(top), agentID)
  return 0
}

func runTrainPropose(args []string, base string) int {
  var agentID string
  threshold := 2
  for _, a := range args {
    if strings.HasPrefix(a, "--agent=") {
      agentID = strings.TrimPrefix(a, "--agent=")
    } else if strings.HasPrefix(a, "--threshold=") {
      fmt.Sscanf(strings.TrimPrefix(a, "--threshold="), "%d", &threshold)
    }
  }
  if agentID == "" {
    fmt.Fprintln(os.Stderr, "usage: zazi train propose --agent=<id> [--threshold=<n>]")
    return 1
  }

  reg := az.New(base)
  a, err := reg.Get(agentID)
  if err != nil {
    fmt.Fprintf(os.Stderr, "agent not found: %v\n", err)
    return 1
  }

  tr := trainer.NewTrainer()
  // Load existing signals
  signalsPath := filepath.Join(a.WorkspaceRoot, "trainer_signals.txt")
  if data, err := os.ReadFile(signalsPath); err == nil {
    lines := strings.Split(string(data), "\n")
    for _, line := range lines {
      // Parse "dimension=value (score:N)" roughly
      if strings.Contains(line, "=") {
        parts := strings.SplitN(line, "=", 2)
        dim := parts[0]
        rest := parts[1]
        var val string
        if idx := strings.Index(rest, " (score:"); idx > 0 {
          val = rest[:idx]
        } else {
          val = rest
        }
        tr.Profile().Apply(trainer.Signal{Dimension: dim, Value: val})
      }
    }
  }

  patch := tr.ProposePatch(threshold)
  if patch == "" {
    fmt.Println("No persona patch proposed (insufficient evidence).")
    return 0
  }

  if trainer.IsMajorChange(patch) {
    fmt.Println("=== PROPOSED MAJOR PERSONA PATCH ===")
    fmt.Println(patch)
    fmt.Println("=== This is a major change. Review before applying. ===")
  } else {
    fmt.Println("=== PROPOSED PERSONA PATCH ===")
    fmt.Println(patch)
  }
  return 0
}

func runTrainApply(args []string, base string) int {
  var agentID, patch string
  for _, a := range args {
    if strings.HasPrefix(a, "--agent=") {
      agentID = strings.TrimPrefix(a, "--agent=")
    } else if strings.HasPrefix(a, "--patch=") {
      patch = strings.TrimPrefix(a, "--patch=")
    }
  }
  if agentID == "" || patch == "" {
    fmt.Fprintln(os.Stderr, "usage: zazi train apply --agent=<id> --patch=<text>")
    return 1
  }

  reg := az.New(base)
  a, err := reg.Get(agentID)
  if err != nil {
    fmt.Fprintf(os.Stderr, "agent not found: %v\n", err)
    return 1
  }

  if trainer.IsMajorChange(patch) {
    fmt.Println("Warning: This patch involves a major change (tone/autonomy).")
    fmt.Println("Run with --force to apply without approval.")
    return 1
  }

  if err := trainer.ApplyPatch(a.WorkspaceRoot, patch); err != nil {
    fmt.Fprintf(os.Stderr, "apply patch: %v\n", err)
    return 1
  }
  fmt.Println("Patch applied successfully.")
  return 0
}

func runJobCancel(args []string, base string) int {
  var agentID, jobID string
  for _, a := range args {
    if strings.HasPrefix(a, "--agent=") {
      agentID = strings.TrimPrefix(a, "--agent=")
    } else if strings.HasPrefix(a, "--job=") {
      jobID = strings.TrimPrefix(a, "--job=")
    }
  }
  if agentID == "" || jobID == "" {
    fmt.Fprintln(os.Stderr, "usage: zazi job cancel --agent=<id> --job=<job_id>")
    return 1
  }

  reg := az.New(base)
  a, err := reg.Get(agentID)
  if err != nil {
    fmt.Fprintf(os.Stderr, "agent not found: %v\n", err)
    return 1
  }

  q := jobs.NewQueue(a.WorkspaceRoot)
  cancelled, err := q.Cancel(jobID)
  if err != nil {
    fmt.Fprintf(os.Stderr, "cancel job: %v\n", err)
    return 1
  }
  fmt.Printf("Job %s cancelled\n", cancelled.ID)
  return 0
}

func runWhatsAppSetup(args []string, base string) int {
  var agentID, phoneID, token, verifyToken string
  for _, a := range args {
    if strings.HasPrefix(a, "--agent=") {
      agentID = strings.TrimPrefix(a, "--agent=")
    } else if strings.HasPrefix(a, "--phone=") {
      phoneID = strings.TrimPrefix(a, "--phone=")
    } else if strings.HasPrefix(a, "--token=") {
      token = strings.TrimPrefix(a, "--token=")
    } else if strings.HasPrefix(a, "--verify=") {
      verifyToken = strings.TrimPrefix(a, "--verify=")
    }
  }
  if agentID == "" {
    fmt.Fprintln(os.Stderr, "usage: zazi whatsapp setup --agent=<id> [--phone=<phone_id>] [--token=<token>] [--verify=<verify_token>]")
    return 1
  }

  reg := az.New(base)
  if _, err := reg.Get(agentID); err != nil {
    fmt.Fprintf(os.Stderr, "agent not found: %v\n", err)
    return 1
  }

  cfgPath, _ := zkconfig.DefaultPath()
  cfg, err := zkconfig.Load(cfgPath)
  if err != nil {
    cfg = zkconfig.NewDefault()
  }
  if cfg.Secrets == nil {
    cfg.Secrets = make(map[string]string)
  }

  if token != "" {
    cfg.Secrets["whatsapp_"+agentID+"_access_token"] = token
  }
  if phoneID != "" {
    cfg.Secrets["whatsapp_"+agentID+"_phone_id"] = phoneID
  }
  if verifyToken != "" {
    cfg.Secrets["whatsapp_"+agentID+"_verify_token"] = verifyToken
  }

  if token != "" || phoneID != "" || verifyToken != "" {
    if err := zkconfig.Save(cfgPath, cfg); err != nil {
      fmt.Fprintf(os.Stderr, "failed to save config: %v\n", err)
      return 1
    }
    fmt.Printf("WhatsApp configuration saved for agent %s\n", agentID)
    if token != "" {
      fmt.Printf("  Access token: %s\n", channels.RedactToken(token))
    }
    if phoneID != "" {
      fmt.Printf("  Phone number ID: %s\n", phoneID)
    }
    if verifyToken != "" {
      fmt.Printf("  Verify token: %s\n", channels.RedactToken(verifyToken))
    }
  } else {
    existingToken := cfg.Secrets["whatsapp_"+agentID+"_access_token"]
    if existingToken != "" {
      fmt.Printf("Agent %s has WhatsApp configured (token: %s)\n", agentID, channels.RedactToken(existingToken))
    } else {
      fmt.Printf("Agent %s has no WhatsApp configuration (dry-run mode only)\n", agentID)
    }
  }
  return 0
}

func runWhatsAppTest(args []string, base string) int {
  var agentID, phoneNumber, message string
  for _, a := range args {
    if strings.HasPrefix(a, "--agent=") {
      agentID = strings.TrimPrefix(a, "--agent=")
    } else if strings.HasPrefix(a, "--phone=") {
      phoneNumber = strings.TrimPrefix(a, "--phone=")
    } else if strings.HasPrefix(a, "--message=") {
      message = strings.TrimPrefix(a, "--message=")
    }
  }
  if agentID == "" || phoneNumber == "" || message == "" {
    fmt.Fprintln(os.Stderr, "usage: zazi whatsapp test --agent=<id> --phone=<phone_number> --message=<msg>")
    return 1
  }

  reg := az.New(base)
  a, err := reg.Get(agentID)
  if err != nil {
    fmt.Fprintf(os.Stderr, "agent not found: %v\n", err)
    return 1
  }

  adapter := channels.NewWhatsAppAdapter("", "", "")
  adapter.BindPhone(phoneNumber, agentID)

  msg := adapter.TestEvent(phoneNumber, "test-user", message)
  fmt.Printf("[TEST] Received WhatsApp message for agent %s: %s\n", msg.AgentID, msg.Content)

  job := jobs.NewJob(agentID, "whatsapp", message)
  evt := jobs.JobEvent{
    JobID:     job.ID,
    AgentID:   agentID,
    EventType: "job.created",
    Message:   "WhatsApp message from " + phoneNumber,
    CreatedAt: job.CreatedAt,
  }
  _ = jobs.AppendEvent(a.WorkspaceRoot, evt)

  provider := llm.NewMockProvider()
  response, _ := provider.Complete(message)

  job.Status = "completed"
  job.Result = response
  now := time.Now()
  job.UpdatedAt = now
  job.CompletedAt = &now

  evt = jobs.JobEvent{
    JobID:     job.ID,
    AgentID:   agentID,
    EventType: "job.completed",
    Message:   response,
    CreatedAt: now,
  }
  _ = jobs.AppendEvent(a.WorkspaceRoot, evt)

  _ = adapter.SendMessage(phoneNumber, response)

  fmt.Printf("[TEST] Response: %s\n", response)
  fmt.Println("Test completed successfully (dry-run mode, no real API calls)")
  return 0
}

func runWhatsAppVerify(args []string) int {
  var mode, token, challenge, verifyToken string
  for _, a := range args {
    if strings.HasPrefix(a, "--mode=") {
      mode = strings.TrimPrefix(a, "--mode=")
    } else if strings.HasPrefix(a, "--token=") {
      token = strings.TrimPrefix(a, "--token=")
    } else if strings.HasPrefix(a, "--challenge=") {
      challenge = strings.TrimPrefix(a, "--challenge=")
    } else if strings.HasPrefix(a, "--verify=") {
      verifyToken = strings.TrimPrefix(a, "--verify=")
    }
  }
  if mode == "" || token == "" || challenge == "" || verifyToken == "" {
    fmt.Fprintln(os.Stderr, "usage: zazi whatsapp verify --mode=<mode> --token=<token> --challenge=<challenge> --verify=<verify_token>")
    return 1
  }

  adapter := channels.NewWhatsAppAdapter("", "", verifyToken)
  result, err := adapter.VerifyWebhook(mode, token, challenge)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Webhook verification failed: %v\n", err)
    return 1
  }
  fmt.Printf("Webhook verified successfully. Challenge response: %s\n", result)
  return 0
}

func runDiscordSetup(args []string, base string) int {
  var agentID, token string
  for _, a := range args {
    if strings.HasPrefix(a, "--agent=") {
      agentID = strings.TrimPrefix(a, "--agent=")
    } else if strings.HasPrefix(a, "--token=") {
      token = strings.TrimPrefix(a, "--token=")
    }
  }
  if agentID == "" {
    fmt.Fprintln(os.Stderr, "usage: zazi discord setup --agent=<id> [--token=<token>]")
    return 1
  }

  reg := az.New(base)
  if _, err := reg.Get(agentID); err != nil {
    fmt.Fprintf(os.Stderr, "agent not found: %v\n", err)
    return 1
  }

  cfgPath, _ := zkconfig.DefaultPath()
  cfg, err := zkconfig.Load(cfgPath)
  if err != nil {
    cfg = zkconfig.NewDefault()
  }
  if cfg.Secrets == nil {
    cfg.Secrets = make(map[string]string)
  }

  if token != "" {
    if err := channels.ValidateDiscordToken(token); err != nil {
      fmt.Fprintf(os.Stderr, "invalid token: %v\n", err)
      return 1
    }
    cfg.Secrets["discord_"+agentID+"_token"] = token
    if err := zkconfig.Save(cfgPath, cfg); err != nil {
      fmt.Fprintf(os.Stderr, "failed to save config: %v\n", err)
      return 1
    }
    fmt.Printf("Discord token configured for agent %s (redacted: %s)\n", agentID, channels.RedactToken(token))
  } else {
    existing := cfg.Secrets["discord_"+agentID+"_token"]
    if existing != "" {
      fmt.Printf("Agent %s has a Discord token configured (%s)\n", agentID, channels.RedactToken(existing))
    } else {
      fmt.Printf("Agent %s has no Discord token (dry-run mode only)\n", agentID)
    }
  }
  return 0
}

func runDiscordTest(args []string, base string) int {
  var agentID, channelID, message string
  for _, a := range args {
    if strings.HasPrefix(a, "--agent=") {
      agentID = strings.TrimPrefix(a, "--agent=")
    } else if strings.HasPrefix(a, "--channel=") {
      channelID = strings.TrimPrefix(a, "--channel=")
    } else if strings.HasPrefix(a, "--message=") {
      message = strings.TrimPrefix(a, "--message=")
    }
  }
  if agentID == "" || channelID == "" || message == "" {
    fmt.Fprintln(os.Stderr, "usage: zazi discord test --agent=<id> --channel=<channel_id> --message=<msg>")
    return 1
  }

  reg := az.New(base)
  a, err := reg.Get(agentID)
  if err != nil {
    fmt.Fprintf(os.Stderr, "agent not found: %v\n", err)
    return 1
  }

  adapter := channels.NewDiscordAdapter("")
  adapter.BindChannel(channelID, agentID)

  msg := adapter.TestEvent(channelID, "test-user", message)
  fmt.Printf("[TEST] Received Discord message for agent %s: %s\n", msg.AgentID, msg.Content)

  job := jobs.NewJob(agentID, "discord", message)
  evt := jobs.JobEvent{
    JobID:     job.ID,
    AgentID:   agentID,
    EventType: "job.created",
    Message:   "Discord message from channel " + channelID,
    CreatedAt: job.CreatedAt,
  }
  _ = jobs.AppendEvent(a.WorkspaceRoot, evt)

  provider := llm.NewMockProvider()
  response, _ := provider.Complete(message)

  job.Status = "completed"
  job.Result = response
  now := time.Now()
  job.UpdatedAt = now
  job.CompletedAt = &now

  evt = jobs.JobEvent{
    JobID:     job.ID,
    AgentID:   agentID,
    EventType: "job.completed",
    Message:   response,
    CreatedAt: now,
  }
  _ = jobs.AppendEvent(a.WorkspaceRoot, evt)

  _ = adapter.SendMessage(channelID, response)

  fmt.Printf("[TEST] Response: %s\n", response)
  fmt.Println("Test completed successfully (dry-run mode, no real API calls)")
  return 0
}

func runTelegramSetup(args []string, base string) int {
  var agentID, token string
  for _, a := range args {
    if strings.HasPrefix(a, "--agent=") {
      agentID = strings.TrimPrefix(a, "--agent=")
    } else if strings.HasPrefix(a, "--token=") {
      token = strings.TrimPrefix(a, "--token=")
    }
  }
  if agentID == "" {
    fmt.Fprintln(os.Stderr, "usage: zazi telegram setup --agent=<id> [--token=<token>]")
    return 1
  }

  reg := az.New(base)
  if _, err := reg.Get(agentID); err != nil {
    fmt.Fprintf(os.Stderr, "agent not found: %v\n", err)
    return 1
  }

  // Load config and update secrets
  cfgPath, _ := zkconfig.DefaultPath()
  cfg, err := zkconfig.Load(cfgPath)
  if err != nil {
    cfg = zkconfig.NewDefault()
  }
  if cfg.Secrets == nil {
    cfg.Secrets = make(map[string]string)
  }

  if token != "" {
    if err := channels.ValidateToken(token); err != nil {
      fmt.Fprintf(os.Stderr, "invalid token: %v\n", err)
      return 1
    }
    cfg.Secrets["telegram_"+agentID+"_token"] = token
    if err := zkconfig.Save(cfgPath, cfg); err != nil {
      fmt.Fprintf(os.Stderr, "failed to save config: %v\n", err)
      return 1
    }
    fmt.Printf("Telegram token configured for agent %s (redacted: %s)\n", agentID, channels.RedactToken(token))
  } else {
    // Show current status
    existing := cfg.Secrets["telegram_"+agentID+"_token"]
    if existing != "" {
      fmt.Printf("Agent %s has a Telegram token configured (%s)\n", agentID, channels.RedactToken(existing))
    } else {
      fmt.Printf("Agent %s has no Telegram token (dry-run mode only)\n", agentID)
    }
  }
  return 0
}

func runTelegramTest(args []string, base string) int {
  var agentID, chatID, message string
  for _, a := range args {
    if strings.HasPrefix(a, "--agent=") {
      agentID = strings.TrimPrefix(a, "--agent=")
    } else if strings.HasPrefix(a, "--chat=") {
      chatID = strings.TrimPrefix(a, "--chat=")
    } else if strings.HasPrefix(a, "--message=") {
      message = strings.TrimPrefix(a, "--message=")
    }
  }
  if agentID == "" || chatID == "" || message == "" {
    fmt.Fprintln(os.Stderr, "usage: zazi telegram test --agent=<id> --chat=<chat_id> --message=<msg>")
    return 1
  }

  reg := az.New(base)
  a, err := reg.Get(agentID)
  if err != nil {
    fmt.Fprintf(os.Stderr, "agent not found: %v\n", err)
    return 1
  }

  // Create adapter in dry-run mode (no token)
  adapter := channels.NewTelegramAdapter("")
  adapter.BindChat(chatID, agentID)

  // Simulate inbound message
  msg := adapter.TestMessage(chatID, "test-user", message)
  fmt.Printf("[TEST] Received message for agent %s: %s\n", msg.AgentID, msg.Content)

  // Create job and process
  job := jobs.NewJob(agentID, "telegram", message)
  evt := jobs.JobEvent{
    JobID:     job.ID,
    AgentID:   agentID,
    EventType: "job.created",
    Message:   "Telegram message from chat " + chatID,
    CreatedAt: job.CreatedAt,
  }
  _ = jobs.AppendEvent(a.WorkspaceRoot, evt)

  // Use mock provider
  provider := llm.NewMockProvider()
  response, _ := provider.Complete(message)

  job.Status = "completed"
  job.Result = response
  now := time.Now()
  job.UpdatedAt = now
  job.CompletedAt = &now

  evt = jobs.JobEvent{
    JobID:     job.ID,
    AgentID:   agentID,
    EventType: "job.completed",
    Message:   response,
    CreatedAt: now,
  }
  _ = jobs.AppendEvent(a.WorkspaceRoot, evt)

  // Simulate sending response back
  _ = adapter.SendMessage(chatID, response)

  fmt.Printf("[TEST] Response: %s\n", response)
  fmt.Println("Test completed successfully (dry-run mode, no real API calls)")
  return 0
}

func runTelegramSend(args []string, base string) int {
  var agentID, chatID, message string
  for _, a := range args {
    if strings.HasPrefix(a, "--agent=") {
      agentID = strings.TrimPrefix(a, "--agent=")
    } else if strings.HasPrefix(a, "--chat=") {
      chatID = strings.TrimPrefix(a, "--chat=")
    } else if strings.HasPrefix(a, "--message=") {
      message = strings.TrimPrefix(a, "--message=")
    }
  }
  if agentID == "" || chatID == "" || message == "" {
    fmt.Fprintln(os.Stderr, "usage: zazi telegram send --agent=<id> --chat=<chat_id> --message=<msg>")
    return 1
  }

  // Load token from config
  cfgPath, _ := zkconfig.DefaultPath()
  cfg, err := zkconfig.Load(cfgPath)
  if err != nil {
    cfg = zkconfig.NewDefault()
  }
  token := cfg.Secrets["telegram_"+agentID+"_token"]
  if token == "" {
    fmt.Fprintln(os.Stderr, "No Telegram token configured for this agent.")
    fmt.Fprintln(os.Stderr, "Run: zazi telegram setup --agent="+agentID+" --token=<YOUR_TOKEN>")
    fmt.Fprintln(os.Stderr, "Or use test mode: zazi telegram test --agent="+agentID+" --chat="+chatID+" --message="+message)
    return 1
  }

  adapter := channels.NewTelegramAdapter(token)
  if err := adapter.SendMessage(chatID, message); err != nil {
    fmt.Fprintf(os.Stderr, "failed to send: %v\n", err)
    return 1
  }
  fmt.Println("Message sent successfully")
  return 0
}

func runTelegramListen(args []string, base string) int {
  var agentID string
  for _, a := range args {
    if strings.HasPrefix(a, "--agent=") {
      agentID = strings.TrimPrefix(a, "--agent=")
    }
  }
  if agentID == "" {
    fmt.Fprintln(os.Stderr, "usage: zazi telegram listen --agent=<id>")
    return 1
  }

  reg := az.New(base)
  a, err := reg.Get(agentID)
  if err != nil {
    fmt.Fprintf(os.Stderr, "agent not found: %v\n", err)
    return 1
  }

  cfgPath, _ := zkconfig.DefaultPath()
  cfg, _ := zkconfig.Load(cfgPath)
  token := cfg.Secrets["telegram_"+agentID+"_token"]
  if token == "" {
    fmt.Fprintln(os.Stderr, "No Telegram token configured for this agent.")
    fmt.Fprintln(os.Stderr, "Run: zazi telegram setup --agent="+agentID+" --token=<YOUR_TOKEN>")
    return 1
  }

  ollamaKey := cfg.Secrets["ollama_api_key"]
  router := llm.NewRouter(60 * time.Second)
  if ollamaKey != "" {
    router.Register(llm.TaskCheap, llm.NewOllamaProvider(ollamaKey, ""))
  }
  router.Register(llm.TaskCheap, llm.NewMockProvider())

  adapter := channels.NewTelegramAdapter(token)
  adapter.BindChat("any", agentID)

  fmt.Printf("[telegram] Listening for messages for agent %s...\n", agentID)
  fmt.Println("Press Ctrl+C to stop.")

  ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
  defer stop()

  handler := func(chatID int, fromName, text string) string {
    fmt.Printf("[telegram] Message from %s: %s\n", fromName, text)

    job := jobs.NewJob(agentID, "telegram", text)
    q := jobs.NewQueue(a.WorkspaceRoot)
    _ = q.Enqueue(job)

    evt := jobs.JobEvent{
      JobID:     job.ID,
      AgentID:   agentID,
      EventType: "job.created",
      Message:   text,
      CreatedAt: job.CreatedAt,
    }
    _ = jobs.AppendEvent(a.WorkspaceRoot, evt)

    response, providerName, err := router.Route(llm.TaskCheap, text)
    if err != nil {
      fmt.Fprintf(os.Stderr, "[telegram] LLM error: %v\n", err)
      job.Status = "failed"
      job.LastError = err.Error()
      job.UpdatedAt = time.Now()
      _ = q.Update(job)
      return "Sorry, I couldn't process that."
    }

    job.Status = "completed"
    job.Result = response
    now := time.Now()
    job.UpdatedAt = now
    job.CompletedAt = &now
    _ = q.Update(job)

    evt = jobs.JobEvent{
      JobID:     job.ID,
      AgentID:   agentID,
      EventType: "job.completed",
      Message:   response,
      Metadata:  map[string]string{"provider": providerName},
      CreatedAt: now,
    }
    _ = jobs.AppendEvent(a.WorkspaceRoot, evt)

    fmt.Printf("[%s] %s\n", providerName, response)
    return response
  }

  if err := adapter.Listen(ctx, handler); err != nil {
    fmt.Fprintf(os.Stderr, "[telegram] listen stopped: %v\n", err)
    return 1
  }
  return 0
}

func runChat(args []string) int {
  var agentID, message string
  for _, a := range args {
    if strings.HasPrefix(a, "--agent=") {
      agentID = strings.TrimPrefix(a, "--agent=")
    } else if strings.HasPrefix(a, "--message=") {
      message = strings.TrimPrefix(a, "--message=")
    }
  }
  if agentID == "" || message == "" {
    fmt.Fprintln(os.Stderr, "usage: zazi chat --agent=<id> --message=<msg>")
    return 1
  }

  home, err := homeDir()
  if err != nil {
    fmt.Fprintln(os.Stderr, "can't determine home:", err)
    return 1
  }
  base := filepath.Join(home, ".zazi")
  reg := az.New(base)

  a, err := reg.Get(agentID)
  if err != nil {
    fmt.Fprintf(os.Stderr, "agent not found: %v\n", err)
    return 1
  }

  job := jobs.NewJob(agentID, "cli", message)

  // Persist job to queue
  q := jobs.NewQueue(a.WorkspaceRoot)
  _ = q.Enqueue(job)

  // Log job created event
  evt := jobs.JobEvent{
    JobID:     job.ID,
    AgentID:   agentID,
    EventType: "job.created",
    Message:   "Job created for chat",
    CreatedAt: job.CreatedAt,
  }
  _ = jobs.AppendEvent(a.WorkspaceRoot, evt)

  // Use LLM router with fallback support
  // Try Ollama Cloud first if configured, fallback to mock
  router := llm.NewRouter(60 * time.Second)
  cfgPath, _ := zkconfig.DefaultPath()
  cfg, _ := zkconfig.Load(cfgPath)
  ollamaKey := cfg.Secrets["ollama_api_key"]
  if ollamaKey != "" {
    router.Register(llm.TaskCheap, llm.NewOllamaProvider(ollamaKey, ""))
    router.Register(llm.TaskCheap, llm.NewMockProvider())
    router.Register(llm.TaskStrong, llm.NewOllamaProvider(ollamaKey, ""))
    router.Register(llm.TaskStrong, llm.NewMockProvider())
    router.Register(llm.TaskCoding, llm.NewOllamaProvider(ollamaKey, ""))
    router.Register(llm.TaskCoding, llm.NewMockProvider())
  } else {
    router.Register(llm.TaskCheap, llm.NewMockProvider())
    router.Register(llm.TaskStrong, llm.NewMockProvider())
    router.Register(llm.TaskCoding, llm.NewMockProvider())
  }

  response, providerName, err := router.Route(llm.TaskCheap, message)
  if err != nil {
    job.Status = "failed"
    job.LastError = err.Error()
    job.UpdatedAt = time.Now()
    _ = q.Update(job)

    // Check for loop
    ld := jobs.NewLoopDetector(3)
    isLoop, reason := ld.Check([]jobs.JobEvent{evt, {EventType: "job.failed", Message: err.Error()}})
    if isLoop {
      _, _ = q.Pause(job.ID, reason)
      loopEvt := jobs.JobEvent{
        JobID:     job.ID,
        AgentID:   agentID,
        EventType: "job.paused",
        Message:   reason,
        CreatedAt: time.Now(),
      }
      _ = jobs.AppendEvent(a.WorkspaceRoot, loopEvt)
    }

    fmt.Fprintf(os.Stderr, "provider error: %v\n", err)
    return 1
  }

  job.Status = "completed"
  job.Result = response
  now := time.Now()
  job.UpdatedAt = now
  job.CompletedAt = &now
  _ = q.Update(job)

  // Log completion event with provider info
  evt = jobs.JobEvent{
    JobID:     job.ID,
    AgentID:   agentID,
    EventType: "job.completed",
    Message:   response,
    Metadata:  map[string]string{"provider": providerName},
    CreatedAt: now,
  }
  _ = jobs.AppendEvent(a.WorkspaceRoot, evt)

  fmt.Printf("[%s] %s\n", providerName, response)
  return 0
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

func runInit() error {
  // Use the new config loader to initialize defaults and directories.
  home, err := homeDir()
  if err != nil {
    return err
  }
  // Default config path
  path, err := zkconfig.DefaultPath()
  if err != nil {
    return err
  }
  if _, statErr := os.Stat(path); statErr == nil {
    fmt.Printf("Config already exists at %s. Skipping write.\n", path)
    return nil
  }
  // Create default config and ensure data directories
  cfg := zkconfig.NewDefault()
  if err := zkconfig.Save(path, cfg); err != nil {
    return fmt.Errorf("saving default config: %w", err)
  }
  if err := zkconfig.EnsureDirs(cfg); err != nil {
    return fmt.Errorf("ensuring data directories: %w", err)
  }
  fmt.Printf("Initialized Zazi at %s\n", path)
  _ = home // silence if unused in this path; home is not strictly needed
  return nil
}

func runDoctor() {
  // Load current config using the loader and verify basic health.
  path, _ := zkconfig.DefaultPath()
  c, err := zkconfig.Load(path)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
    return
  }

  // Config file existence
  if _, statErr := os.Stat(path); statErr != nil {
    fmt.Printf("Config file: %s - MISSING\n", path)
  } else {
    fmt.Printf("Config file: %s - OK\n", path)
  }

  // Data path writability check
  dataPath := zkconfig.ExpandPath(c.DataPath)
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
    fmt.Printf("Data path: %s - not writable (check permissions)\n", dataPath)
  }

  // Print non-secret config values for verification
  fmt.Printf("Config: DataPath=%q, LogLevel=%q, ServerPort=%d, LLMProviders=%v\n",
    c.DataPath, c.LogLevel, c.ServerPort, c.LLMProviders)

  if !writable {
    fmt.Println("Hint: Ensure the data directory exists and is writable, or run 'zazi init'.")
  }
}

func runProvider(args []string) int {
  if len(args) == 0 {
    printProviderUsage()
    return 0
  }
  switch args[0] {
  case "list":
    return runProviderList(args[1:])
  case "show":
    return runProviderShow(args[1:])
  default:
    printProviderUsage()
    return 0
  }
}

func printProviderUsage() {
  fmt.Fprintln(os.Stdout, "Provider subcommands:")
  fmt.Fprintln(os.Stdout, "  zazi provider list [--tier=cheap|strong|coding|vision]")
  fmt.Fprintln(os.Stdout, "  zazi provider show --name=<provider>")
}

func runProviderList(args []string) int {
  tier := ""
  for _, a := range args {
    if strings.HasPrefix(a, "--tier=") {
      tier = strings.TrimPrefix(a, "--tier=")
    }
  }

  var providers []llm.ProviderInfo
  if tier != "" {
    providers = llm.ProvidersByTier(tier)
  } else {
    providers = llm.BuiltInProviders
  }

  fmt.Printf("Providers (%d total, %d models)\n", llm.TotalProviderCount(), llm.TotalModelCount())
  for _, p := range providers {
    fmt.Printf("  %-20s %s (%d models, default: %s)\n", p.Name, p.BaseURL, len(p.Models), p.DefaultModel)
  }
  return 0
}

func runProviderShow(args []string) int {
  var name string
  for _, a := range args {
    if strings.HasPrefix(a, "--name=") {
      name = strings.TrimPrefix(a, "--name=")
    }
  }
  if name == "" {
    fmt.Fprintln(os.Stderr, "usage: zazi provider show --name=<provider>")
    return 1
  }

  p, err := llm.FindProvider(name)
  if err != nil {
    fmt.Fprintf(os.Stderr, "%v\n", err)
    return 1
  }

  fmt.Printf("Provider: %s\n", p.Name)
  fmt.Printf("Base URL: %s\n", p.BaseURL)
  fmt.Printf("Default Model: %s\n", p.DefaultModel)
  fmt.Printf("Models (%d):\n", len(p.Models))
  for _, m := range p.Models {
    caps := strings.Join(m.Capabilities, ", ")
    if caps == "" {
      caps = "-"
    }
    fmt.Printf("  %-40s ctx=%-8d tier=%-8s caps=%s\n", m.ID, m.ContextWindow, m.Tier, caps)
  }
  return 0
}
