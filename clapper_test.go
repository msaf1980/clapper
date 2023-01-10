package clapper

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

/*----------------*/

// test unsupported flag
func TestUnsupportedAssignment(t *testing.T) {

	// options
	options := map[string][]string{
		"---version": []string{"---version"},
		"---v":       []string{"---v=1.0.0"},
		"-version":   []string{"-version"},
	}

	for flag, options := range options {
		// command
		cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			t.Fatalf("Error: %v, out: %q", err, string(output))
		} else {
			if !strings.Contains(fmt.Sprintf("%s", output), fmt.Sprintf(`error => clapper.ErrorUnsupportedFlag{Name:"%s"}`, flag)) {
				t.Fail()
			}
		}
	}
}

// test empty root command
func TestEmptyRootCommand(t *testing.T) {
	// command
	cmd := exec.Command("go", "run", "demo/cmd.go")

	// get output
	if output, err := cmd.Output(); err != nil {
		t.Fatalf("Error: %v, out: %q", err, string(output))
	} else {
		lines := []string{
			`sub-command => ""`,
			`argument(output) => &clapper.Arg{Name:"output", IsVariadic:false, Value:""}`,
			`flag(force) => &clapper.Flag{Name:"force", IsBoolean:true, Value:"false"}`,
			`flag(verbose) => &clapper.Flag{Name:"verbose", IsBoolean:true, Value:"false"}`,
			`flag(version) => &clapper.Flag{Name:"version", IsBoolean:false, Value:""}`,
			`flag(dir) => &clapper.Flag{Name:"dir", IsBoolean:false, Value:"/var/users"}`,
		}

		for _, line := range lines {
			if !strings.Contains(fmt.Sprintf("%s", output), line) {
				t.Fatalf("got\n%q\nwant line\n%q", output, line)
			}
		}
	}
}

// test root command when not registered
func TestUnregisteredRootCommand(t *testing.T) {
	// command
	cmd := exec.Command("go", "run", "demo/cmd.go")
	cmd.Env = append(os.Environ(), "NO_ROOT=TRUE")

	// get output
	if output, err := cmd.Output(); err != nil {
		t.Fatalf("Error: %v, out: %q", err, string(output))
	} else {
		lines := []string{
			`error => clapper.ErrorUnknownCommand{Name:""}`,
		}

		for _, line := range lines {
			if !strings.Contains(fmt.Sprintf("%s", output), line) {
				t.Fatalf("got\n%q\nwant line\n%q", output, line)
			}
		}
	}
}

// test an unregistered flag
func TestUnregisteredFlag(t *testing.T) {

	// flags
	flags := map[string][]string{
		"-d":          []string{"-V", "1.0.1", "-v", "--force", "-d", "./sub/dir"},
		"--m":         []string{"-V", "1.0.1", "-v", "--force", "--m", "./sub/dir"},
		"--directory": []string{"-V", "1.0.1", "-v", "--force", "--directory", "./sub/dir"},
	}

	for flag, options := range flags {
		// command
		cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			t.Fatalf("Error: %v, out: %q", err, string(output))
		} else {
			if !strings.Contains(fmt.Sprintf("%s", output), fmt.Sprintf(`error => clapper.ErrorUnknownFlag{Name:"%s"}`, flag)) {
				t.Fail()
			}
		}
	}

}

// test for valid inverted flag values
func TestValidInvertFlagValues(t *testing.T) {

	// options list
	optionsList := [][]string{
		[]string{"info", "student", "-V", "-v", "--output", "./opt/dir", "--no-clean"},
		[]string{"info", "student", "--version", "--no-clean", "--output", "./opt/dir", "--verbose"},
	}

	for _, options := range optionsList {
		// command
		cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			t.Fatalf("Error: %v, out: %q", err, string(output))
		} else {
			lines := []string{
				`sub-command => "info"`,
				`argument(category) => &clapper.Arg{Name:"category", IsVariadic:false, Value:"student"}`,
				`argument(username) => &clapper.Arg{Name:"username", IsVariadic:false, Value:""}`,
				`argument(subjects) => &clapper.Arg{Name:"subjects", IsVariadic:true, Value:""}`,
				`flag(version) => &clapper.Flag{Name:"version", IsBoolean:false, Value:"1.0.1"}`,
				`flag(output) => &clapper.Flag{Name:"output", IsBoolean:false, Value:"./opt/dir"}`,
				`flag(verbose) => &clapper.Flag{Name:"verbose", IsBoolean:true, Value:"true"}`,
				`flag(clean) => &clapper.Flag{Name:"clean", IsBoolean:true, Value:"false"}`,
			}

			for _, line := range lines {
				if !strings.Contains(fmt.Sprintf("%s", output), line) {
					t.Fatalf("got\n%q\nwant line\n%q", output, line)
				}
			}
		}
	}
}

// test for invalid flag error when an inverted flag is used without `--no-` prefix
func TestErrorUnknownFlagForInvertFlags(t *testing.T) {

	// options list
	optionsList := map[string][]string{
		"--clean":   []string{"info", "student", "-V", "-v", "--output", "./opt/dir", "--clean"},
		"--no-dump": []string{"info", "student", "--version", "--no-dump", "--output", "./opt/dir", "--verbose"},
	}

	for flag, options := range optionsList {
		// command
		cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			t.Fatalf("Error: %v, out: %q", err, string(output))
		} else {
			if !strings.Contains(fmt.Sprintf("%s", output), fmt.Sprintf(`error => clapper.ErrorUnknownFlag{Name:"%s"}`, flag)) {
				t.Fail()
			}
		}
	}
}

// test `--flag=value` syntax
func TestFlagAssignmentSyntax(t *testing.T) {

	// options list
	optionsList := [][]string{
		[]string{"info", "student", "-v", "--version=2.0.0", "thatisuday"},
		[]string{"info", "student", "thatisuday", "-v", "-V=2.0.0"},
	}

	for _, options := range optionsList {
		// command
		cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			t.Fatalf("Error: %v, out: %q", err, string(output))
		} else {
			lines := []string{
				`sub-command => "info"`,
				`argument(category) => &clapper.Arg{Name:"category", IsVariadic:false, Value:"student"}`,
				`argument(username) => &clapper.Arg{Name:"username", IsVariadic:false, Value:"thatisuday"}`,
				`argument(subjects) => &clapper.Arg{Name:"subjects", IsVariadic:true, Value:""}`,
				`flag(version) => &clapper.Flag{Name:"version", IsBoolean:false, Value:"2.0.0"}`,
				`flag(output) => &clapper.Flag{Name:"output", IsBoolean:false, Value:"./"}`,
				`flag(verbose) => &clapper.Flag{Name:"verbose", IsBoolean:true, Value:"true"}`,
			}

			for _, line := range lines {
				if !strings.Contains(fmt.Sprintf("%s", output), line) {
					t.Fatalf("got\n%q\nwant line\n%q", output, line)
				}
			}
		}
	}
}

// test for valid variadic argument values
func TestValidVariadicArgumentValues(t *testing.T) {

	// options list
	optionsList := [][]string{
		[]string{"info", "student", "thatisuday", "-V", "-v", "--output", "./opt/dir", "--no-clean", "math", "science", "physics"},
		[]string{"info", "student", "--version", "--no-clean", "thatisuday", "--output", "./opt/dir", "math", "science", "--verbose", "physics"},
	}

	for _, options := range optionsList {
		// command
		cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			t.Fatalf("Error: %v, out: %q", err, string(output))
		} else {
			lines := []string{
				`sub-command => "info"`,
				`argument(category) => &clapper.Arg{Name:"category", IsVariadic:false, Value:"student"}`,
				`argument(username) => &clapper.Arg{Name:"username", IsVariadic:false, Value:"thatisuday"}`,
				`argument(subjects) => &clapper.Arg{Name:"subjects", IsVariadic:true, Value:"math,science,physics"}`,
				`flag(version) => &clapper.Flag{Name:"version", IsBoolean:false, Value:"1.0.1"}`,
				`flag(output) => &clapper.Flag{Name:"output", IsBoolean:false, Value:"./opt/dir"}`,
				`flag(verbose) => &clapper.Flag{Name:"verbose", IsBoolean:true, Value:"true"}`,
				`flag(clean) => &clapper.Flag{Name:"clean", IsBoolean:true, Value:"false"}`,
			}

			for _, line := range lines {
				if !strings.Contains(fmt.Sprintf("%s", output), line) {
					t.Fatalf("got\n%q\nwant line\n%q", output, line)
				}
			}
		}
	}
}

/*-------------------*/

// test root command with options
func TestRootCommandWithOptions(t *testing.T) {

	// options list
	optionsList := [][]string{
		[]string{"userinfo", "-V", "1.0.1", "-v", "--force", "--dir", "./sub/dir"},
		[]string{"-V", "1.0.1", "--verbose", "--force", "userinfo", "--dir", "./sub/dir"},
		[]string{"-V", "1.0.1", "-v", "--force", "--dir", "./sub/dir", "userinfo"},
		[]string{"--version", "1.0.1", "--verbose", "--force", "--dir", "./sub/dir", "userinfo"},
	}

	for _, options := range optionsList {
		// command
		cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			t.Fatalf("Error: %v, out: %q", err, string(output))
		} else {
			lines := []string{
				`sub-command => ""`,
				`argument(output) => &clapper.Arg{Name:"output", IsVariadic:false, Value:"userinfo"}`,
				`flag(force) => &clapper.Flag{Name:"force", IsBoolean:true, Value:"true"}`,
				`flag(verbose) => &clapper.Flag{Name:"verbose", IsBoolean:true, Value:"true"}`,
				`flag(version) => &clapper.Flag{Name:"version", IsBoolean:false, Value:"1.0.1"}`,
				`flag(dir) => &clapper.Flag{Name:"dir", IsBoolean:false, Value:"./sub/dir"}`,
			}

			for _, line := range lines {
				if !strings.Contains(fmt.Sprintf("%s", output), line) {
					t.Fatalf("got\n%q\nwant line\n%q", output, line)
				}
			}
		}
	}
}

// test sub-command with options
func TestSubCommandWithOptions(t *testing.T) {

	// options list
	optionsList := [][]string{
		[]string{"info", "student", "-V", "-v", "--output", "./opt/dir"},
		[]string{"info", "student", "--version", "--output", "./opt/dir", "--verbose"},
	}

	for _, options := range optionsList {
		// command
		cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			t.Fatalf("Error: %v, out: %q", err, string(output))
		} else {
			lines := []string{
				`sub-command => "info"`,
				`argument(category) => &clapper.Arg{Name:"category", IsVariadic:false, Value:"student"}`,
				`argument(username) => &clapper.Arg{Name:"username", IsVariadic:false, Value:""}`,
				`argument(subjects) => &clapper.Arg{Name:"subjects", IsVariadic:true, Value:""}`,
				`flag(version) => &clapper.Flag{Name:"version", IsBoolean:false, Value:"1.0.1"}`,
				`flag(output) => &clapper.Flag{Name:"output", IsBoolean:false, Value:"./opt/dir"}`,
				`flag(verbose) => &clapper.Flag{Name:"verbose", IsBoolean:true, Value:"true"}`,
				`flag(clean) => &clapper.Flag{Name:"clean", IsBoolean:true, Value:"true"}`,
			}

			for _, line := range lines {
				if !strings.Contains(fmt.Sprintf("%s", output), line) {
					t.Fatalf("got\n%q\nwant line\n%q", output, line)
				}
			}
		}
	}
}

// test sub-command with valid and extra arguments
func TestSubCommandWithArguments(t *testing.T) {

	// options list
	optionsList := [][]string{
		[]string{"info", "-v", "student", "-V", "2.0.0", "thatisuday"},
		[]string{"info", "student", "-v", "thatisuday", "--version", "2.0.0"},
	}

	for _, options := range optionsList {
		// command
		cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			t.Fatalf("Error: %v, out: %q", err, string(output))
		} else {
			lines := []string{
				`sub-command => "info"`,
				`argument(category) => &clapper.Arg{Name:"category", IsVariadic:false, Value:"student"}`,
				`argument(username) => &clapper.Arg{Name:"username", IsVariadic:false, Value:"thatisuday"}`,
				`argument(subjects) => &clapper.Arg{Name:"subjects", IsVariadic:true, Value:""}`,
				`flag(version) => &clapper.Flag{Name:"version", IsBoolean:false, Value:"2.0.0"}`,
				`flag(output) => &clapper.Flag{Name:"output", IsBoolean:false, Value:"./"}`,
				`flag(verbose) => &clapper.Flag{Name:"verbose", IsBoolean:true, Value:"true"}`,
			}

			for _, line := range lines {
				if !strings.Contains(fmt.Sprintf("%s", output), line) {
					t.Fatalf("got\n%q\nwant line\n%q", output, line)
				}
			}
		}
	}
}

// test validate arg
func TestInvalidArg(t *testing.T) {
	// options
	options := []string{"info", "worker", "-V", "-v", "2.0.0"}

	// command
	cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

	// get output
	if output, err := cmd.Output(); err != nil {
		t.Fatalf("Error: %v, out: %q", err, string(output))
	} else {
		out := string(output)
		want := "error => clapper.ErrorUnsupportedValue{Name:\"category\", Value:\"worker\"}\n"
		if out != want {
			t.Fatalf("got\n%q\nwant\n%q", out, want)
		}
	}
}

// test validate flag
func TestInvalidFlag(t *testing.T) {
	// options
	options := []string{"info", "student", "-V", "2.0.1", "-v"}

	// command
	cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

	// get output
	if output, err := cmd.Output(); err != nil {
		t.Fatalf("Error: %v, out: %q", err, string(output))
	} else {
		out := string(output)
		want := "error => clapper.ErrorUnsupportedValue{Name:\"version\", Value:\"2.0.1\"}\n"
		if out != want {
			t.Fatalf("got\n%q\nwant\n%q", out, want)
		}
	}
}
