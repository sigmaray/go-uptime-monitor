package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"go-uptime-monitor/models"

	"github.com/rs/zerolog/log"
	"golang.org/x/term"
	"gorm.io/gorm"
)

func Run(db *gorm.DB, args []string) {
	if len(args) == 0 {
		PrintUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "users-seed":
		usersSeed(db)
	case "users-create":
		usersCreate(db)
	case "users-show":
		usersShow(db)
	case "users-delete-all":
		usersDeleteAll(db)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", args[0])
		PrintUsage()
		os.Exit(1)
	}
}

func PrintUsage() {
	fmt.Println(`Usage: uptime-monitor <command>

Server:
  s, server               Start the HTTP server

Database:
  migrate                 Apply pending database migrations

User commands:
  users-seed              Create admin user (admin/admin)
  users-create            Create a user interactively
  users-show              List all users
  users-delete-all        Delete all users`)
}

func usersSeed(db *gorm.DB) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("WARNING: Creating user with login \"admin\" and password \"admin\" is insecure and dangerous in production.")
	if !confirmYes(reader, "Do you want to continue? [y/N]: ") {
		fmt.Println("Aborted.")
		return
	}

	existing, err := models.FindUserByUsername(db, "admin")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to check admin")
	}
	if existing != nil {
		fmt.Println("User 'admin' already exists")
		return
	}

	input := models.CreateUserInput{Username: "admin", Password: "admin", ConfirmPassword: "admin"}
	if err := input.Validate(); err != nil {
		log.Fatal().Str("error", models.FormatValidationError(err)).Msg("Invalid input")
	}

	user, err := models.CreateUser(db, input.Username, input.Password)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create admin user")
	}
	fmt.Printf("Created user: id=%d username=%s\n", user.ID, user.Username)
}

func usersCreate(db *gorm.DB) {
	reader := bufio.NewReader(os.Stdin)

	username := readLine(reader, "Login: ")
	password := readPassword(reader, "Password: ")
	confirm := readPassword(reader, "Confirm password: ")

	input := models.CreateUserInput{Username: username, Password: password, ConfirmPassword: confirm}
	if err := input.Validate(); err != nil {
		log.Fatal().Str("error", models.FormatValidationError(err)).Msg("Invalid input")
	}

	user, err := models.CreateUser(db, input.Username, input.Password)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create user")
	}

	fmt.Printf("Created user: id=%d username=%s\n", user.ID, user.Username)
}

func usersShow(db *gorm.DB) {
	var users []models.User
	if err := db.Order("id asc").Find(&users).Error; err != nil {
		log.Fatal().Err(err).Msg("Failed to fetch users")
	}

	if len(users) == 0 {
		fmt.Println("No users found.")
		return
	}

	fmt.Printf("%-5s %-20s %-20s\n", "ID", "Username", "Created At")
	fmt.Println(strings.Repeat("-", 47))
	for _, user := range users {
		fmt.Printf("%-5d %-20s %-20s\n", user.ID, user.Username, user.CreatedAt.Format("2006-01-02 15:04"))
	}
}

func usersDeleteAll(db *gorm.DB) {
	reader := bufio.NewReader(os.Stdin)

	if !confirmYes(reader, "This will permanently delete ALL users. Are you sure? [y/N]: ") {
		fmt.Println("Aborted.")
		return
	}

	result := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&models.User{})
	if result.Error != nil {
		log.Fatal().Err(result.Error).Msg("Failed to delete users")
	}
	fmt.Printf("Deleted %d user(s).\n", result.RowsAffected)
}

func confirmYes(reader *bufio.Reader, prompt string) bool {
	answer := strings.ToLower(readLine(reader, prompt))
	return answer == "y" || answer == "yes"
}

func readLine(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read input")
	}
	return strings.TrimSpace(line)
}

func readPassword(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println() // print newline after password
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read input")
	}
	return strings.TrimSpace(string(bytePassword))
}
