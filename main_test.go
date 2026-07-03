package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	
	"github.com/creack/pty"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func runCLI(t *testing.T, envs []string, input string, args ...string) (string, string) {
	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	
	// Start with current envs and append our test overrides
	env := os.Environ()
	cmd.Env = append(env, envs...)

	if input != "" {
		ptmx, err := pty.Start(cmd)
		if err != nil {
			t.Fatalf("pty.Start failed: %v", err)
		}
		
		go func() {
			ptmx.WriteString(input)
		}()

		var buf bytes.Buffer
		io.Copy(&buf, ptmx)
		
		err = cmd.Wait()
		ptmx.Close()
		
		if err != nil && !strings.Contains(err.Error(), "exit status 1") {
			t.Logf("Command failed: %v\nOutput: %s", err, buf.String())
		}
		
		return buf.String(), buf.String()
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Only log if it's not a generic exit status 1 which is often expected
		if !strings.Contains(err.Error(), "exit status 1") {
			t.Logf("Command failed: %v\nStderr: %s\nStdout: %s", err, stderr.String(), stdout.String())
		}
	}
	return stdout.String(), stderr.String()
}

func getTestDBName() string {
	return os.Getenv("GO_UPTIME_MONITOR_TEST_DATABASE_NAME")
}

func setupTestDB(t *testing.T) {
	// Connect to default postgres db to create the test db
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to postgres: %v", err)
	}
	defer db.Close()
	
	testDB := getTestDBName()
	_, err = db.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, testDB))
	if err != nil {
		t.Fatalf("Failed to drop old test database: %v", err)
	}

	_, err = db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, testDB))
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
}

func cleanupTestDB(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		t.Errorf("Failed to connect to postgres for cleanup: %v", err)
		return
	}
	defer db.Close()

	testDB := getTestDBName()

	// Terminate any active connections
	_, _ = db.Exec(fmt.Sprintf(`
		SELECT pg_terminate_backend(pg_stat_activity.pid)
		FROM pg_stat_activity
		WHERE pg_stat_activity.datname = '%s'
		  AND pid <> pg_backend_pid();
	`, testDB))

	_, err = db.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, testDB))
	if err != nil {
		t.Errorf("Failed to drop test database: %v", err)
	}
}

func TestCLICommands(t *testing.T) {
	godotenv.Load() // load .env if available

	if os.Getenv("GO_UPTIME_MONITOR_TEST_DATABASE_NAME") == "" {
		t.Skip("GO_UPTIME_MONITOR_TEST_DATABASE_NAME is not set, skipping tests")
	}

	setupTestDB(t)
	defer cleanupTestDB(t)

	envs := []string{
		fmt.Sprintf("GO_UPTIME_MONITOR_DATABASE_NAME=%s", getTestDBName()),
		"GO_UPTIME_MONITOR_SESSION_SECRET=testsecret",
	}

	t.Run("Migrate", func(t *testing.T) {
		stdout, stderr := runCLI(t, envs, "", "migrate")
		if strings.Contains(stderr, "level\":\"fatal") || strings.Contains(stderr, "level\":\"error") {
			t.Errorf("Migration failed: %s", stderr)
		}
		_ = stdout
	})

	t.Run("Seed Users", func(t *testing.T) {
		// Answer yes to the confirmation prompt
		stdout, _ := runCLI(t, envs, "y\n", "users-seed")
		if !strings.Contains(stdout, "Created user: id=1 username=admin") {
			t.Errorf("Expected user to be created, got: %s", stdout)
		}

		// Try seeding again, should say already exists
		stdout, _ = runCLI(t, envs, "y\n", "users-seed")
		if !strings.Contains(stdout, "already exists") {
			t.Errorf("Expected already exists message, got: %s", stdout)
		}
		
		// Try aborting
		stdout, _ = runCLI(t, envs, "n\n", "users-seed")
		if !strings.Contains(stdout, "Aborted") {
			t.Errorf("Expected Aborted message, got: %s", stdout)
		}
	})

	t.Run("Show Users", func(t *testing.T) {
		stdout, _ := runCLI(t, envs, "", "users-show")
		if !strings.Contains(stdout, "admin") {
			t.Errorf("Expected admin in users list, got: %s", stdout)
		}
	})

	t.Run("Create User", func(t *testing.T) {
		input := "testcli\npass123\npass123\n"
		stdout, _ := runCLI(t, envs, input, "users-create")
		if !strings.Contains(stdout, "username=testcli") {
			t.Errorf("Expected testcli to be created, got: %s", stdout)
		}

		stdout, _ = runCLI(t, envs, "", "users-show")
		if !strings.Contains(stdout, "testcli") {
			t.Errorf("Expected testcli in users list, got: %s", stdout)
		}
	})

	t.Run("Create User Validation Error", func(t *testing.T) {
		input := "bad\npass\npass123\n"
		_, stderr := runCLI(t, envs, input, "users-create")
		if !strings.Contains(stderr, "Invalid input") {
			t.Errorf("Expected validation error in stderr, got: %s", stderr)
		}
	})

	t.Run("Delete All Users", func(t *testing.T) {
		// Abort
		stdout, _ := runCLI(t, envs, "n\n", "users-delete-all")
		if !strings.Contains(stdout, "Aborted") {
			t.Errorf("Expected Aborted message, got: %s", stdout)
		}

		// Confirm
		stdout, _ = runCLI(t, envs, "y\n", "users-delete-all")
		if !strings.Contains(stdout, "Deleted 2 user(s)") {
			t.Errorf("Expected Deleted 2 user(s) message, got: %s", stdout)
		}

		stdout, _ = runCLI(t, envs, "", "users-show")
		if !strings.Contains(stdout, "No users found") {
			t.Errorf("Expected No users found message, got: %s", stdout)
		}
	})
	
	t.Run("Unknown Command", func(t *testing.T) {
		_, stderr := runCLI(t, envs, "", "unknown-cmd")
		if !strings.Contains(stderr, "Unknown command: unknown-cmd") {
			t.Errorf("Expected Unknown command message, got: %s", stderr)
		}
	})
	
	t.Run("No arguments", func(t *testing.T) {
		stdout, _ := runCLI(t, envs, "")
		if !strings.Contains(stdout, "Usage: uptime-monitor") {
			t.Errorf("Expected usage message, got: %s", stdout)
		}
	})
}
