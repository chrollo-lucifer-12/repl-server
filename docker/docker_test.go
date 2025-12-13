package docker

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

func TestDockerClientWithNodeJSProject(t *testing.T) {
	// Create Docker client
	client := NewDockerClient()
	defer client.Stop()

	ctx := context.Background()

	// Start container
	t.Log("Starting container...")
	var startBuf bytes.Buffer
	containerID := client.StartContainer(ctx, &startBuf)
	if containerID == "" {
		t.Fatal("Failed to start container")
	}
	t.Logf("Container started: %s", containerID)
	t.Logf("Container start output:\n%s", startBuf.String())

	// Cleanup container at the end
	defer func() {
		t.Log("Cleaning up container...")
		if err := client.DeleteContainer(ctx, containerID); err != nil {
			t.Logf("Warning: Failed to delete container: %v", err)
		}
	}()

	// Wait a bit for container to be fully ready
	time.Sleep(2 * time.Second)

	// Test 1: Check Node.js is installed
	t.Run("CheckNodeInstalled", func(t *testing.T) {
		var buf bytes.Buffer
		if err := client.ExecCommand(ctx, containerID, []string{"node", "--version"}, &buf); err != nil {
			t.Fatalf("Node.js not installed: %v", err)
		}
		version := strings.TrimSpace(buf.String())
		t.Logf("Node.js version: %s", version)
		t.Logf("Command output:\n%s", buf.String())
		if !strings.HasPrefix(version, "v") {
			t.Error("Unexpected Node.js version format")
		}
	})

	// Test 2: Check NPM is installed
	t.Run("CheckNPMInstalled", func(t *testing.T) {
		var buf bytes.Buffer
		if err := client.ExecCommand(ctx, containerID, []string{"npm", "--version"}, &buf); err != nil {
			t.Fatalf("NPM not installed: %v", err)
		}
		version := strings.TrimSpace(buf.String())
		t.Logf("NPM version: %s", version)
		t.Logf("Command output:\n%s", buf.String())
	})

	// Test 3: Check if /usr/local/bin/node exists
	t.Run("CheckNodeBinaryExists", func(t *testing.T) {
		var buf bytes.Buffer
		cmd := []string{"sh", "-c", "test -f /usr/local/bin/node && echo 'exists' || echo 'not found'"}
		if err := client.ExecCommand(ctx, containerID, cmd, &buf); err != nil {
			t.Fatalf("Failed to check node binary: %v", err)
		}
		output := strings.TrimSpace(buf.String())
		t.Logf("Node binary check result: %s", output)
		t.Logf("Command output:\n%s", buf.String())
		if !strings.Contains(output, "exists") {
			t.Error("Node binary not found at /usr/local/bin/node")
		}
	})

	// Test 4: Check working directory exists
	t.Run("CheckWorkingDirectory", func(t *testing.T) {
		var buf bytes.Buffer
		cmd := []string{"sh", "-c", "test -d /home/hi && echo 'exists' || echo 'not found'"}
		if err := client.ExecCommand(ctx, containerID, cmd, &buf); err != nil {
			t.Fatalf("Failed to check working directory: %v", err)
		}
		output := strings.TrimSpace(buf.String())
		t.Logf("Working directory check result: %s", output)
		t.Logf("Command output:\n%s", buf.String())
		if !strings.Contains(output, "exists") {
			t.Error("Working directory /home/hi not found")
		}
	})

	// Test 5: Check if package managers are accessible
	t.Run("CheckPackageManagerCommands", func(t *testing.T) {
		commands := []string{"node", "npm", "npx"}
		for _, cmd := range commands {
			var buf bytes.Buffer
			checkCmd := []string{"sh", "-c", "which " + cmd}
			if err := client.ExecCommand(ctx, containerID, checkCmd, &buf); err != nil {
				t.Errorf("Command %s not found in PATH: %v", cmd, err)
				continue
			}
			path := strings.TrimSpace(buf.String())
			t.Logf("%s found at: %s", cmd, path)
			t.Logf("Command output:\n%s", buf.String())
		}
	})

	// Test 6: Test Node.js can execute simple JavaScript
	t.Run("ExecuteSimpleJavaScript", func(t *testing.T) {
		var buf bytes.Buffer
		cmd := []string{"node", "-e", "console.log('Hello from Node.js')"}
		if err := client.ExecCommand(ctx, containerID, cmd, &buf); err != nil {
			t.Fatalf("Failed to execute JavaScript: %v", err)
		}
		output := strings.TrimSpace(buf.String())
		t.Logf("JavaScript execution output: %s", output)
		t.Logf("Full command output:\n%s", buf.String())
		if !strings.Contains(output, "Hello from Node.js") {
			t.Error("JavaScript execution failed")
		}
	})

	// Test 7: Check Node.js built-in modules
	t.Run("CheckNodeBuiltinModules", func(t *testing.T) {
		var buf bytes.Buffer
		cmd := []string{"node", "-e", "const fs = require('fs'); const path = require('path'); console.log('Modules OK')"}
		if err := client.ExecCommand(ctx, containerID, cmd, &buf); err != nil {
			t.Fatalf("Failed to load built-in modules: %v", err)
		}
		output := strings.TrimSpace(buf.String())
		t.Logf("Built-in modules check: %s", output)
		t.Logf("Command output:\n%s", buf.String())
		if !strings.Contains(output, "Modules OK") {
			t.Error("Built-in modules not working")
		}
	})

	// Test 8: Verify shell is available
	t.Run("CheckShellAvailable", func(t *testing.T) {
		var buf bytes.Buffer
		cmd := []string{"sh", "-c", "echo 'Shell works'"}
		if err := client.ExecCommand(ctx, containerID, cmd, &buf); err != nil {
			t.Fatalf("Shell not available: %v", err)
		}
		output := strings.TrimSpace(buf.String())
		t.Logf("Shell test output: %s", output)
		t.Logf("Command output:\n%s", buf.String())
		if !strings.Contains(output, "Shell works") {
			t.Error("Shell execution failed")
		}
	})

	// Test 9: Check if standard Unix utilities exist
	t.Run("CheckUnixUtilities", func(t *testing.T) {
		utilities := []string{"ls", "cat", "mkdir", "rm", "mv", "grep", "echo"}
		for _, util := range utilities {
			var buf bytes.Buffer
			cmd := []string{"sh", "-c", "which " + util}
			if err := client.ExecCommand(ctx, containerID, cmd, &buf); err != nil {
				t.Errorf("Utility %s not found: %v", util, err)
				t.Logf("Command output:\n%s", buf.String())
				continue
			}
			path := strings.TrimSpace(buf.String())
			t.Logf("Utility %s found at: %s", util, path)
			t.Logf("Command output:\n%s", buf.String())
		}
	})

	// Test 10: Check Node.js can create and read files
	t.Run("NodeCanCreateReadFiles", func(t *testing.T) {
		var buf bytes.Buffer
		script := "const fs = require('fs'); fs.writeFileSync('/tmp/test.txt', 'test content'); const content = fs.readFileSync('/tmp/test.txt', 'utf8'); console.log('File content:', content)"
		cmd := []string{"node", "-e", script}
		if err := client.ExecCommand(ctx, containerID, cmd, &buf); err != nil {
			t.Fatalf("Node.js file operations failed: %v", err)
		}
		output := buf.String()
		t.Logf("File operations output:\n%s", output)
		if !strings.Contains(output, "test content") {
			t.Error("Node.js couldn't read file content")
		}
	})

	// Test 11: Verify file system is writable
	t.Run("CheckFileSystemWritable", func(t *testing.T) {
		var buf bytes.Buffer
		cmd := []string{"sh", "-c", "touch /tmp/testfile && test -f /tmp/testfile && echo 'File created successfully' || echo 'Failed to create file'"}
		if err := client.ExecCommand(ctx, containerID, cmd, &buf); err != nil {
			t.Fatalf("Failed to test file system: %v", err)
		}
		output := strings.TrimSpace(buf.String())
		t.Logf("File system writability test: %s", output)
		t.Logf("Command output:\n%s", buf.String())
		if !strings.Contains(output, "created successfully") {
			t.Error("File system is not writable")
		}
	})

	// Test 12: Check environment variables
	t.Run("CheckEnvironmentVariables", func(t *testing.T) {
		var buf bytes.Buffer
		cmd := []string{"sh", "-c", "echo PATH=$PATH"}
		if err := client.ExecCommand(ctx, containerID, cmd, &buf); err != nil {
			t.Fatalf("Failed to read environment: %v", err)
		}
		path := strings.TrimSpace(buf.String())
		t.Logf("Environment PATH: %s", path)
		t.Logf("Command output:\n%s", buf.String())
		if !strings.Contains(path, "PATH=") {
			t.Error("PATH environment variable is not set properly")
		}
	})

	// Test 13: List contents of working directory
	t.Run("ListWorkingDirectory", func(t *testing.T) {
		var buf bytes.Buffer
		cmd := []string{"ls", "-la", "/home/hi"}
		if err := client.ExecCommand(ctx, containerID, cmd, &buf); err != nil {
			t.Fatalf("Failed to list working directory: %v", err)
		}
		t.Logf("Working directory contents:\n%s", buf.String())
	})

	// Test 14: Check disk space
	t.Run("CheckDiskSpace", func(t *testing.T) {
		var buf bytes.Buffer
		cmd := []string{"df", "-h", "/home/hi"}
		if err := client.ExecCommand(ctx, containerID, cmd, &buf); err != nil {
			t.Fatalf("Failed to check disk space: %v", err)
		}
		t.Logf("Disk space information:\n%s", buf.String())
	})

	// Test 15: Test complex Node.js operation
	t.Run("ComplexNodeJSOperation", func(t *testing.T) {
		var buf bytes.Buffer
		script := `
const fs = require('fs');
const path = require('path');

// Create a directory
const testDir = '/tmp/node-test';
if (!fs.existsSync(testDir)) {
  fs.mkdirSync(testDir);
  console.log('Created directory:', testDir);
}

// Write multiple files
const files = ['file1.txt', 'file2.txt', 'file3.txt'];
files.forEach(file => {
  const filePath = path.join(testDir, file);
  fs.writeFileSync(filePath, 'Content of ' + file);
  console.log('Created file:', filePath);
});

// List files
const dirContents = fs.readdirSync(testDir);
console.log('Directory contents:', dirContents.join(', '));

// Read and verify
files.forEach(file => {
  const content = fs.readFileSync(path.join(testDir, file), 'utf8');
  console.log(file + ':', content);
});

console.log('All operations completed successfully!');
`
		cmd := []string{"node", "-e", script}
		if err := client.ExecCommand(ctx, containerID, cmd, &buf); err != nil {
			t.Fatalf("Complex Node.js operation failed: %v", err)
		}
		output := buf.String()
		t.Logf("Complex Node.js operation output:\n%s", output)

		if !strings.Contains(output, "All operations completed successfully") {
			t.Error("Complex operation did not complete successfully")
		}
		if !strings.Contains(output, "file1.txt") {
			t.Error("Failed to create file1.txt")
		}
	})

	// Test 16: Install a simple NPM package and use it
	t.Run("InstallAndUseNPMPackage", func(t *testing.T) {
		// Create a project directory
		var buf bytes.Buffer
		cmd := []string{"mkdir", "-p", "/home/hi/test-project"}
		if err := client.ExecCommand(ctx, containerID, cmd, &buf); err != nil {
			t.Fatalf("Failed to create project directory: %v", err)
		}
		t.Logf("Create directory output:\n%s", buf.String())

		// Initialize npm project (create package.json)
		buf.Reset()
		initCmd := []string{"sh", "-c", "cd /home/hi/test-project && npm init -y"}
		if err := client.ExecCommand(ctx, containerID, initCmd, &buf); err != nil {
			t.Fatalf("Failed to initialize npm project: %v", err)
		}
		t.Logf("NPM init output:\n%s", buf.String())

		// Install a simple package (chalk - for colored console output)
		buf.Reset()
		t.Log("Installing 'chalk' package (this may take a moment)...")
		installCmd := []string{"sh", "-c", "cd /home/hi/test-project && npm install chalk@4.1.2"}
		if err := client.ExecCommand(ctx, containerID, installCmd, &buf); err != nil {
			t.Fatalf("Failed to install chalk package: %v", err)
		}
		t.Logf("NPM install output:\n%s", buf.String())

		// Check if node_modules directory exists
		buf.Reset()
		checkCmd := []string{"sh", "-c", "test -d /home/hi/test-project/node_modules && echo 'node_modules exists' || echo 'node_modules NOT found'"}
		if err := client.ExecCommand(ctx, containerID, checkCmd, &buf); err != nil {
			t.Fatalf("Failed to check node_modules: %v", err)
		}
		output := buf.String()
		t.Logf("node_modules check output:\n%s", output)
		if !strings.Contains(output, "exists") {
			t.Error("node_modules directory was not created")
		}

		// List contents of node_modules
		buf.Reset()
		listCmd := []string{"ls", "-la", "/home/hi/test-project/node_modules"}
		if err := client.ExecCommand(ctx, containerID, listCmd, &buf); err != nil {
			t.Fatalf("Failed to list node_modules: %v", err)
		}
		t.Logf("node_modules contents:\n%s", buf.String())

		// Check if chalk package exists
		buf.Reset()
		chalkCheckCmd := []string{"sh", "-c", "test -d /home/hi/test-project/node_modules/chalk && echo 'chalk package exists' || echo 'chalk NOT found'"}
		if err := client.ExecCommand(ctx, containerID, chalkCheckCmd, &buf); err != nil {
			t.Fatalf("Failed to check chalk package: %v", err)
		}
		t.Logf("Chalk package check:\n%s", buf.String())
		if !strings.Contains(buf.String(), "exists") {
			t.Error("chalk package was not installed")
		}

		// Check package.json exists
		buf.Reset()
		packageCheckCmd := []string{"sh", "-c", "test -f /home/hi/test-project/package.json && cat /home/hi/test-project/package.json"}
		if err := client.ExecCommand(ctx, containerID, packageCheckCmd, &buf); err != nil {
			t.Fatalf("Failed to read package.json: %v", err)
		}
		t.Logf("package.json contents:\n%s", buf.String())

		// Check package-lock.json exists
		buf.Reset()
		lockCheckCmd := []string{"sh", "-c", "test -f /home/hi/test-project/package-lock.json && echo 'package-lock.json exists' || echo 'package-lock.json NOT found'"}
		if err := client.ExecCommand(ctx, containerID, lockCheckCmd, &buf); err != nil {
			t.Fatalf("Failed to check package-lock.json: %v", err)
		}
		t.Logf("package-lock.json check:\n%s", buf.String())

		// Use the installed package
		buf.Reset()
		useChalkScript := `
const chalk = require('chalk');
console.log(chalk.blue('Hello from chalk!'));
console.log(chalk.green('Chalk is working!'));
console.log(chalk.red('Success!'));
console.log('Package test completed');
`
		useCmd := []string{"sh", "-c", "cd /home/hi/test-project && node -e \"" + useChalkScript + "\""}
		if err := client.ExecCommand(ctx, containerID, useCmd, &buf); err != nil {
			t.Fatalf("Failed to use chalk package: %v", err)
		}
		output = buf.String()
		t.Logf("Chalk usage output:\n%s", output)

		if !strings.Contains(output, "Package test completed") {
			t.Error("Failed to use chalk package properly")
		}
		if !strings.Contains(output, "Chalk is working") {
			t.Error("Chalk output not found")
		}
	})

	// Test 17: Verify node_modules structure
	t.Run("VerifyNodeModulesStructure", func(t *testing.T) {
		// Count files in node_modules
		var buf bytes.Buffer
		countCmd := []string{"sh", "-c", "find /home/hi/test-project/node_modules -type f | wc -l"}
		if err := client.ExecCommand(ctx, containerID, countCmd, &buf); err != nil {
			t.Fatalf("Failed to count node_modules files: %v", err)
		}
		t.Logf("Total files in node_modules:\n%s", buf.String())

		// Count directories in node_modules
		buf.Reset()
		countDirCmd := []string{"sh", "-c", "find /home/hi/test-project/node_modules -type d | wc -l"}
		if err := client.ExecCommand(ctx, containerID, countDirCmd, &buf); err != nil {
			t.Fatalf("Failed to count node_modules directories: %v", err)
		}
		t.Logf("Total directories in node_modules:\n%s", buf.String())

		// List top-level packages
		buf.Reset()
		listPackagesCmd := []string{"sh", "-c", "ls -1 /home/hi/test-project/node_modules"}
		if err := client.ExecCommand(ctx, containerID, listPackagesCmd, &buf); err != nil {
			t.Fatalf("Failed to list packages: %v", err)
		}
		t.Logf("Top-level packages in node_modules:\n%s", buf.String())

		// Check .bin directory
		buf.Reset()
		binCheckCmd := []string{"sh", "-c", "test -d /home/hi/test-project/node_modules/.bin && ls -la /home/hi/test-project/node_modules/.bin || echo '.bin directory not found'"}
		if err := client.ExecCommand(ctx, containerID, binCheckCmd, &buf); err != nil {
			t.Logf("Failed to check .bin directory: %v", err)
		}
		t.Logf(".bin directory contents:\n%s", buf.String())
	})

	t.Log("All tests completed successfully!")
}

// Additional test for container lifecycle
func TestContainerLifecycle(t *testing.T) {
	client := NewDockerClient()
	defer client.Stop()

	ctx := context.Background()

	// Start container
	var buf bytes.Buffer
	containerID := client.StartContainer(ctx, &buf)
	if containerID == "" {
		t.Fatal("Failed to start container")
	}
	t.Logf("Container started: %s", containerID)
	t.Logf("Start output:\n%s", buf.String())

	// Give it time to start
	time.Sleep(2 * time.Second)

	// Test basic command
	buf.Reset()
	if err := client.ExecCommand(ctx, containerID, []string{"echo", "test"}, &buf); err != nil {
		t.Fatalf("Failed to execute command: %v", err)
	}

	output := buf.String()
	t.Logf("Echo command output:\n%s", output)
	if !strings.Contains(output, "test") {
		t.Errorf("Expected 'test' in output, got: %s", output)
	}

	// Check container info
	buf.Reset()
	if err := client.ExecCommand(ctx, containerID, []string{"sh", "-c", "hostname && pwd"}, &buf); err != nil {
		t.Fatalf("Failed to get container info: %v", err)
	}
	t.Logf("Container info:\n%s", buf.String())

	// Stop and remove container
	if err := client.DeleteContainer(ctx, containerID); err != nil {
		t.Fatalf("Failed to delete container: %v", err)
	}

	t.Log("Container lifecycle test passed")
}
