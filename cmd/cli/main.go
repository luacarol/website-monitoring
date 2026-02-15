package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitorings = 3
const delay = 5

func main() {
	showIntroduction()

	for {
		showMenu()
		command := readCommand()

		// if command == 1 {
		// 	fmt.Println("Monitorando...")
		// } else if command == 2 {
		// 	fmt.Println("Exibindo Logs...")
		// } else if command == 0 {
		// 	fmt.Println("Saindo do Programa...")
		// } else {
		// 	fmt.Println("I don't know this command")
		// }

		switch command {
		case 1:
			initMonitoring()
		case 2:
			printLogs()
		case 0:
			fmt.Println("Saindo do Programa...")
			os.Exit(0)
		default:
			fmt.Println("I don't know this command")
			os.Exit(-1)
		}
	}
}

func showIntroduction() {
	name := "Luana"
	var version float32 = 1.25

	// println("Hello, World!") // This also works
	// fmt.Println("Hello, World with Go!")
	fmt.Println("Hello, Mrs.", name)
	fmt.Println("This program is at version", version)

	// Import "reflect"
	// fmt.Println("The type of the variable name is:", reflect.TypeOf(name))
}

func showMenu() {
	fmt.Println("1- Start Monitoring")
	fmt.Println("2- Display Logs")
	fmt.Println("0- Exit Program")
}

func readCommand() int {
	var command int
	// fmt.Scanf("%d", &command)
	fmt.Scan(&command)

	// fmt.Println("The address of my variable is", &command)
	fmt.Println("The chosen command was:", command)
	fmt.Println("")

	return command
}

func initMonitoring() {
	fmt.Println("Monitoring...")
	// site := "https://www.alura.com.br/"
	// site := "https://httpbin.org/status/404" // 200

	// names := []string{"Luana", "Matheus", "Jucelene"}
	// names = append(names, "Irene")
	// fmt.Println(names)

	// var sites = []string{"https://www.alura.com.br/", "https://httpbin.org/status/404", "https://www.caelum.com.br/", "https://httpbin.org/status/500"}

	sites := readSitesFromFile()

	// for i, site := range sites {
	// 	fmt.Println("I'm at position", i, "of my slice and this position has the site:", site)
	// }

	for i := 0; i < monitorings; i++ {
		for i, site := range sites {
			fmt.Println("Testing site", i, ":", site)
			verifySite(site)
		}

		time.Sleep(delay * time.Second)
		fmt.Println("")
	}

	fmt.Println("")
}

func verifySite(site string) {
	response, err := http.Get(site)
	if err != nil {
		fmt.Println("Error accessing the site:", err)
		return
	}
	defer response.Body.Close()

	// fmt.Println("Response:", response)

	if response.StatusCode == 200 {
		fmt.Println("Site:", site, "successfully loaded!")
		registerLog(site, true)
	} else {
		fmt.Println("Site:", site, "has problems. Status Code:", response.StatusCode)
		registerLog(site, false)
	}
}

func readSitesFromFile() []string {
	var sites []string

	file, err := os.Open("sites.txt")
	// file, err := os.ReadFile("sites.txt")
	// fmt.Println(string(file))

	if err != nil {
		fmt.Println("An error occurred while opening the file:", err)
	}

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		// fmt.Println(line)

		sites = append(sites, line)

		if err == io.EOF {
			break
		}
	}

	file.Close()
	return sites
}

func registerLog(site string, status bool) {
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("An error occurred while opening the log file:", err)
	}

	file.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - online: " + strconv.FormatBool(status) + "\n")

	// fmt.Println(file)
	file.Close()
}

func printLogs() {
	fmt.Println("Displaying Logs...")

	file, err := os.ReadFile("log.txt")

	if err != nil {
		fmt.Println("An error occurred while reading the log file:", err)
	}

	fmt.Println(string(file))
}
