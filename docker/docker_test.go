package docker

import (
	"bytes"
	"context"
	"testing"
)

/*
------------------------------------
Helpers
------------------------------------
*/

func setupContainer(t *testing.T) (*DockerClient, context.Context, string) {
	t.Helper()

	d := NewDockerClient()
	ctx := context.Background()

	containerID := d.StartContainer(ctx, nil)
	if containerID == "" {
		t.Fatal("failed to start container")
	}

	t.Cleanup(func() {
		_ = d.DeleteContainer(ctx, containerID)
		_ = d.Stop()
	})

	return d, ctx, containerID
}

/*
------------------------------------
Tests
------------------------------------
*/

// func TestListFiles(t *testing.T) {
// 	d, ctx, containerID := setupContainer(t)
// 	testDir := "/home/hi/testdir"
// 	var buf bytes.Buffer

// 	// 1. Create test directory
// 	if err := d.CreateDir(ctx, containerID, testDir, &buf); err != nil {
// 		t.Fatalf("CreateDir failed: %v", err)
// 	}

// 	// 2. Create some files and directories
// 	files := map[string]string{
// 		"file1.txt": "Hello",
// 		"file2.txt": "World",
// 	}
// 	dirs := []string{"subdir1", "subdir2"}

// 	for name, content := range files {
// 		if err := d.WriteFile(ctx, containerID, testDir+"/"+name, content, &buf); err != nil {
// 			t.Fatalf("WriteFile %s failed: %v", name, err)
// 		}
// 	}

// 	for _, dir := range dirs {
// 		if err := d.CreateDir(ctx, containerID, testDir+"/"+dir, &buf); err != nil {
// 			t.Fatalf("CreateDir %s failed: %v", dir, err)
// 		}
// 	}

// 	// 3. List files
// 	buf.Reset()
// 	if err := d.ListFiles(ctx, containerID, testDir, &buf); err != nil {
// 		t.Fatalf("ListFiles failed: %v", err)
// 	}

// 	// 4. Parse JSON output
// 	var listed []FileInfo
// 	if err := json.Unmarshal(buf.Bytes(), &listed); err != nil {
// 		t.Fatalf("Failed to parse ListFiles output: %v", err)
// 	}

// 	// 5. Verify files exist
// 	foundFiles := make(map[string]bool)
// 	foundDirs := make(map[string]bool)
// 	for _, f := range listed {
// 		if f.Type == "file" {
// 			foundFiles[f.Name] = true
// 		} else if f.Type == "dir" {
// 			foundDirs[f.Name] = true
// 		}
// 	}

// 	for name := range files {
// 		if !foundFiles[name] {
// 			t.Errorf("Expected file %s not found", name)
// 		}
// 	}

// 	for _, dir := range dirs {
// 		if !foundDirs[dir] {
// 			t.Errorf("Expected dir %s not found", dir)
// 		}
// 	}

// 	t.Logf("ListFiles output: %s", buf.String())
// }

// func TestCreateDir(t *testing.T) {
// 	d, ctx, containerID := setupContainer(t)

// 	var buf bytes.Buffer
// 	err := d.CreateDir(ctx, containerID, "/home/hi/testdir", &buf)
// 	if err != nil {
// 		t.Fatalf("CreateDir failed: %v", err)
// 	}

// 	buf.Reset()
// 	err = d.ExecCommand(ctx, containerID, []string{"ls", "/home/hi"}, &buf)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if !bytes.Contains(buf.Bytes(), []byte("testdir")) {
// 		t.Fatalf("directory not created, got: %s", buf.String())
// 	}
// }

// func TestWriteAndReadFile(t *testing.T) {
// 	d, ctx, containerID := setupContainer(t)

// 	var buf bytes.Buffer
// 	path := "/home/hi/hello.txt"
// 	content := "hello docker"

// 	err := d.WriteFile(ctx, containerID, path, content, &buf)
// 	if err != nil {
// 		t.Fatalf("WriteFile failed: %v", err)
// 	}

// 	buf.Reset()
// 	err = d.ReadFile(ctx, containerID, path, &buf)
// 	if err != nil {
// 		t.Fatalf("ReadFile failed: %v", err)
// 	}

// 	if buf.String() != content {
// 		t.Fatalf("unexpected content: %q", buf.String())
// 	}
// }

// // func TestListFiles(t *testing.T) {
// // 	d, ctx, containerID := setupContainer(t)

// // 	_ = d.WriteFile(ctx, containerID, "/home/hi/a.txt", "a", nil)
// // 	_ = d.WriteFile(ctx, containerID, "/home/hi/b.txt", "b", nil)

// // 	var buf bytes.Buffer
// // 	err := d.ListFiles(ctx, containerID, "/home/hi", &buf)
// // 	if err != nil {
// // 		t.Fatal(err)
// // 	}

// // 	var files []FileInfo
// // 	if err := json.Unmarshal(buf.Bytes(), &files); err != nil {
// // 		t.Fatal(err)
// // 	}

// // 	if len(files) < 2 {
// // 		t.Fatalf("expected files, got %+v", files)
// // 	}
// // }

// func TestStatFile(t *testing.T) {
// 	d, ctx, containerID := setupContainer(t)

// 	_ = d.WriteFile(ctx, containerID, "/home/hi/stat.txt", "data", nil)

// 	var buf bytes.Buffer
// 	err := d.StatFile(ctx, containerID, "/home/hi/stat.txt", &buf)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	var info FileInfo
// 	if err := json.Unmarshal(buf.Bytes(), &info); err != nil {
// 		t.Fatal(err)
// 	}

// 	if info.Size == 0 {
// 		t.Fatal("file size should not be zero")
// 	}
// }

// func TestSearchInFile(t *testing.T) {
// 	d, ctx, containerID := setupContainer(t)

// 	content := "hello\nworld\nhello again"
// 	_ = d.WriteFile(ctx, containerID, "/home/hi/search.txt", content, nil)

// 	var buf bytes.Buffer
// 	err := d.SearchInFile(ctx, containerID, "/home/hi/search.txt", "hello", &buf)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if !bytes.Contains(buf.Bytes(), []byte("hello")) {
// 		t.Fatal("search term not found")
// 	}
// }

// func TestRenameFileDir(t *testing.T) {
// 	d, ctx, containerID := setupContainer(t)

// 	_ = d.WriteFile(ctx, containerID, "/home/hi/old.txt", "x", nil)

// 	var buf bytes.Buffer
// 	err := d.RenameFileDir(ctx, containerID, "/home/hi/old.txt", "/home/hi/new.txt", &buf)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	buf.Reset()
// 	_ = d.ExecCommand(ctx, containerID, []string{"ls", "/home/hi"}, &buf)

// 	if !bytes.Contains(buf.Bytes(), []byte("new.txt")) {
// 		t.Fatal("file rename failed")
// 	}
// }

// func TestRemoveFile(t *testing.T) {
// 	d, ctx, containerID := setupContainer(t)

// 	_ = d.WriteFile(ctx, containerID, "/home/hi/delete.txt", "bye", nil)

// 	var buf bytes.Buffer
// 	err := d.RemoveFile(ctx, "/home/hi/delete.txt", containerID, &buf)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	buf.Reset()
// 	_ = d.ExecCommand(ctx, containerID, []string{"ls", "/home/hi"}, &buf)

// 	if bytes.Contains(buf.Bytes(), []byte("delete.txt")) {
// 		t.Fatal("file was not deleted")
// 	}
// }

func TestNodeJSProjectRun(t *testing.T) {
	d, ctx, containerID := setupContainer(t)

	projectDir := "/home/hi/nodeapp"
	var buf bytes.Buffer

	logFiles := func(msg string, path string) {
		var filesBuf bytes.Buffer
		if err := d.ListFiles(ctx, containerID, path, &filesBuf); err != nil {
			t.Logf("Error listing files: %v", err)
		} else {
			t.Logf("%s (%s):\n%s", msg, path, filesBuf.String())
		}
	}

	// 1. Create project directory
	if err := d.CreateDir(ctx, containerID, projectDir, &buf); err != nil {
		t.Fatalf("CreateDir failed: %v", err)
	}
	logFiles("After CreateDir", projectDir)

	// 2. npm init -y
	buf.Reset()
	err := d.ExecCommand(ctx, containerID,
		[]string{"sh", "-c", "cd " + projectDir + " && npm init -y"},
		&buf,
	)
	if err != nil {
		t.Fatalf("npm init failed: %v\n%s", err, buf.String())
	}
	logFiles("After npm init", projectDir)

	// Verify package.json exists
	buf.Reset()
	err = d.ExecCommand(ctx, containerID,
		[]string{"sh", "-c", "test -f " + projectDir + "/package.json && echo ok"},
		&buf,
	)
	if err != nil || !bytes.Contains(buf.Bytes(), []byte("ok")) {
		t.Fatalf("package.json not created: %v\n%s", err, buf.String())
	}
	logFiles("After verifying package.json", projectDir)

	// 3. npm install
	buf.Reset()
	err = d.ExecCommand(ctx, containerID,
		[]string{"sh", "-c", "cd " + projectDir + " && npm install"},
		&buf,
	)
	if err != nil {
		t.Fatalf("npm install failed: %v\n%s", err, buf.String())
	}
	logFiles("After npm install", projectDir)

	// 4. Install lodash
	buf.Reset()
	err = d.ExecCommand(ctx, containerID,
		[]string{"sh", "-c", "cd " + projectDir + " && npm install lodash"},
		&buf,
	)
	if err != nil {
		t.Fatalf("npm install lodash failed: %v\n%s", err, buf.String())
	}
	logFiles("After npm install lodash", projectDir)

	// Verify node_modules exists
	// buf.Reset()
	// err = d.ExecCommand(ctx, containerID,
	// 	[]string{"sh", "-c", "test -d " + projectDir + "/node_modules && echo ok"},
	// 	&buf,
	// )
	// if err != nil || !bytes.Contains(buf.Bytes(), []byte("ok")) {
	// 	t.Fatalf("node_modules not created: %v\n%s", err, buf.String())
	// }
	// logFiles("After verifying node_modules", projectDir)

	// 5. create index.js
	jsCode := `
const _ = require("lodash");
console.log("RESULT:", _.map([1,2,3,4], n => n*2).join(","));
`
	buf.Reset()
	err = d.WriteFile(ctx, containerID, projectDir+"/index.js", jsCode, &buf)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	logFiles("After writing index.js", projectDir)

	// 6. run node index.js
	buf.Reset()
	err = d.ExecCommand(ctx, containerID,
		[]string{"sh", "-c", "cd " + projectDir + " && node index.js"},
		&buf,
	)
	if err != nil {
		t.Fatalf("node run failed: %v\n%s", err, buf.String())
	}
	logFiles("After running node index.js", projectDir)

	// 7. verify output
	output := buf.String()
	expected := "RESULT: 2,4,6,8"
	if !bytes.Contains([]byte(output), []byte(expected)) {
		t.Fatalf("unexpected output:\n%s\nexpected to contain:\n%s", output, expected)
	}

	t.Log("Node app output:", output)
}

func TestNodeHTTPServer(t *testing.T) {
	d, ctx, containerID := setupContainer(t)
	defer d.DeleteContainer(ctx, containerID)

	projectDir := "/home/hi/nodeapp"
	var buf bytes.Buffer

	// 1. Create project directory
	if err := d.CreateDir(ctx, containerID, projectDir, &buf); err != nil {
		t.Fatalf("CreateDir failed: %v", err)
	}

	// 2. Create package.json
	jsPackage := `{
  "name": "nodeapp",
  "version": "1.0.0",
  "main": "index.js",
  "dependencies": {}
}`
	if err := d.WriteFile(ctx, containerID, projectDir+"/package.json", jsPackage, &buf); err != nil {
		t.Fatalf("Write package.json failed: %v", err)
	}

	// 3. Create index.js with HTTP server
	jsCode := `
const http = require('http');

const server = http.createServer((req, res) => {
  res.end("Hello from Docker Node HTTP server");
});

server.listen(3000, () => {
  console.log("Server running on port 3000");
});
`
	if err := d.WriteFile(ctx, containerID, projectDir+"/index.js", jsCode, &buf); err != nil {
		t.Fatalf("Write index.js failed: %v", err)
	}

	// 4. Start the server in background
	buf.Reset()
	err := d.ExecCommand(ctx, containerID,
		[]string{"sh", "-c", "cd " + projectDir + " && node index.js &"},
		&buf,
	)
	if err != nil {
		t.Fatalf("Starting Node server failed: %v\n%s", err, buf.String())
	}

	// 5. Install curl to test server
	buf.Reset()
	err = d.ExecCommand(ctx, containerID,
		[]string{"sh", "-c", "apk add --no-cache curl"},
		&buf,
	)
	if err != nil {
		t.Fatalf("Installing curl failed: %v\n%s", err, buf.String())
	}

	buf.Reset()
	err = d.ExecCommand(ctx, containerID,
		[]string{"sh", "-c", "curl -s http://localhost:3000"},
		&buf,
	)
	if err != nil {
		t.Fatalf("Curl request failed: %v\n%s", err, buf.String())
	}

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("Hello from Docker Node HTTP server")) {
		t.Fatalf("Unexpected output:\n%s", output)
	}

	t.Log("HTTP server response:", output)
}
