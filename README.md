# 🌐 Website Monitor

A simple yet powerful command-line website monitoring tool built in Go that helps you track the availability and status of your websites in real-time.

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Terminal](https://img.shields.io/badge/Terminal-%23054020?style=for-the-badge&logo=gnu-bash&logoColor=white)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)

## 🚀 Features

- **📊 Real-time Website Monitoring**: Monitor multiple websites simultaneously
- **📝 Automatic Logging**: Keep detailed logs with timestamps of all monitoring activities
- **⚡ Fast HTTP Checks**: Quick response time validation using HTTP status codes
- **📁 File-based Configuration**: Easy website management through a simple text file
- **🖥️ Interactive CLI**: User-friendly command-line interface with menu options
- **⏰ Configurable Intervals**: Set custom monitoring frequency and cycles
- **📈 Status Tracking**: Track website uptime and downtime with boolean status indicators

## 🏗️ Project Structure

```
website-monitor/
├── main.go          # Main application code
├── sites.txt        # List of websites to monitor
├── log.txt          # Monitoring logs (auto-generated)
└── main             # Compiled binary
```

## 📋 Prerequisites

- Go 1.19+ installed on your system
- Basic understanding of command-line interfaces

## 🔧 Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/luacarol/website-monitoring.git
   cd website-monitoring
   ```

2. **Build the application**
   ```bash
   go build -o main main.go
   ```

3. **Run the application**
   ```bash
   ./main
   ```

## 📖 Usage

### Starting the Monitor

When you run the application, you'll see an interactive menu:

```
Olá, sra. Luana
Este programa está na versão 1.25
1- Iniciar Monitoramento
2- Exibir Logs
0- Sair do Programa
```

### Menu Options

- **Option 1**: Start monitoring all websites listed in `sites.txt`
- **Option 2**: Display monitoring logs from previous sessions
- **Option 0**: Exit the application

### Configuration

Edit the `sites.txt` file to add or remove websites you want to monitor:

```txt
https://www.alura.com.br/
https://httpbin.org/status/404
https://www.caelum.com.br/
https://httpbin.org/status/500
```

## ⚙️ Configuration Options

You can modify these constants in `main.go` to customize the monitoring behavior:

```go
const monitorings = 3  // Number of monitoring cycles
const delay = 5        // Delay between cycles (seconds)
```

## 📊 Log Format

The application generates logs in the following format:

```
02/01/2006 15:04:05 - https://www.example.com - online: true
02/01/2006 15:04:10 - https://www.example.com - online: false
```

Each log entry includes:
- **Timestamp**: Exact date and time of the check
- **Website URL**: The monitored website
- **Status**: Boolean indicating if the site is accessible (HTTP 200)

## 🔍 How It Works

1. **Website Loading**: Reads website URLs from `sites.txt`
2. **HTTP Requests**: Sends GET requests to each website
3. **Status Validation**: Checks if the response status code is 200 (OK)
4. **Logging**: Records results with timestamps in `log.txt`
5. **Cycle Management**: Repeats the process based on configured intervals

## 🛠️ Technical Details

- **Language**: Go (Golang)
- **HTTP Client**: Built-in `net/http` package
- **File I/O**: Uses `bufio` for efficient file reading
- **Error Handling**: Comprehensive error handling for network and file operations
- **Concurrency**: Sequential monitoring with configurable delays

## 📝 Sample Output

```
Monitorando...
Testando o site 0 : https://www.alura.com.br/
Site: https://www.alura.com.br/ foi carregado com sucesso!
Testando o site 1 : https://httpbin.org/status/404
Site: https://httpbin.org/status/404 está com problemas. Status Code: 404
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📈 Future Enhancements

- [ ] Add email/SMS notifications for downtime
- [ ] Implement concurrent monitoring for better performance
- [ ] Add support for custom HTTP headers and authentication
- [ ] Create a web dashboard for monitoring results
- [ ] Add response time measurement
- [ ] Implement different monitoring protocols (ping, TCP, etc.)
- [ ] Add configuration file support (JSON/YAML)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**Luana Carol**
- GitHub: [@luacarol](https://github.com/luacarol)

## 🙏 Acknowledgments

- Built as part of Go learning journey
- Inspired by the need for simple website monitoring solutions
- Special thanks to the Go community for excellent documentation

---

⭐ **If you found this project helpful, please give it a star!** ⭐
