package docker

import (
	"fmt"
	"os/exec"
)

func PostgresRun() {
	// Command to start the PostgreSQL container
	cmd := exec.Command("docker", "run", "-d", 
    "--name", "postgres-ew", 
		"-e", "POSTGRES_USER=postgres", 
		"-e", "POSTGRES_PASSWORD=password", 
		"-e", "POSTGRES_DB=postgres", 
		"-p", "5432:5432", 
		"postgres:latest")

	// Run the command
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error starting container:", err)
	} else {
		fmt.Println("PostgreSQL container started successfully")
	}
}

func PostgresStart() {
	// Command to start the PostgreSQL container
	cmd := exec.Command("docker", "start", "-d", "--name", "my-postgres", 
		"-e", "POSTGRES_USER=myuser", 
		"-e", "POSTGRES_PASSWORD=mypassword", 
		"-e", "POSTGRES_DB=mydb", 
		"-p", "5432:5432", 
		"postgres:latest")

	// Run the command
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error starting container:", err)
	} else {
		fmt.Println("PostgreSQL container started successfully")
	}
}
