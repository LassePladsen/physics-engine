package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/LassePladsen/physics-engine/graphics"
	"github.com/LassePladsen/physics-engine/logger"
)

const (
	ticksPerSecond          = 60.0
	deltaTime               = 1 / ticksPerSecond
	gravityAcceleration     = 9.81 // m/s^2
	particleRestitution     = 0.2
	windowBoundsRestitution = 0.5
	friction                = 0.1
)

type demo struct {
	name string
	run  func()
}

var availableDemos = []demo{
	{name: "two-particles", run: runTwoParticlesDemo},
	{name: "many-particles", run: runManyParticlesDemo},
	{name: "orbit", run: runOrbitDemo},
	{name: "unstable-orbit", run: runUnstableOrbitDemo},
	{name: "merge", run: runMergeDemo},
}

func main() {
	os.Exit(runDemo(filepath.Base(os.Args[0]), os.Args[1:], os.Stdout, os.Stderr, availableDemos, logger.Init))
}

func runDemo(program string, args []string, stdout, stderr io.Writer, demos []demo, initLogger func()) int {
	if len(args) == 0 {
		printDemoUsage(stdout, program, demos)
		return 0
	}
	if len(args) != 1 {
		printDemoUsage(stderr, program, demos)
		return 2
	}

	for _, demo := range demos {
		if args[0] == demo.name {
			initLogger()
			graphics.InitGraphics(ticksPerSecond)
			demo.run()
			return 0
		}
	}

	printDemoUsage(stderr, program, demos)
	return 2
}

func printDemoUsage(w io.Writer, program string, demos []demo) {
	fmt.Fprintln(w, "Available demos:")
	for _, demo := range demos {
		fmt.Fprintf(w, "  %s\n", demo.name)
	}
	fmt.Fprintf(w, "Usage: %s <demo>\n", program)
}
