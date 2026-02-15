# 🌐 Website Monitor

A complete website monitoring application with Go backend and React frontend, offering real-time monitoring with a modern visual interface.

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![React](https://img.shields.io/badge/react-%2320232a.svg?style=for-the-badge&logo=react&logoColor=%2361DAFB)
![SQLite](https://img.shields.io/badge/sqlite-%2307405e.svg?style=for-the-badge&logo=sqlite&logoColor=white)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)

## 🚀 Features

- **📊 Real-time Dashboard**: Modern web interface with React
- **🔄 Automatic Monitoring**: Continuous verification of multiple websites
- **📝 Logs History**: Detailed records with timestamps of all checks
- **⚡ REST API**: Endpoints to manage sites and view statistics
- **🔌 WebSocket**: Real-time updates for the dashboard
- **💾 Database**: Persistent storage with SQLite
- **🎨 Modern Interface**: Responsive design with TailwindCSS
- **📈 Statistics**: Uptime metrics and status for each website
- **🖥️ CLI Mode**: Command-line version also available

## 🏗️ Project Structure

```
website-monitoring/
├── cmd/
│   ├── cli/            # CLI application
│   │   └── main.go
│   └── server/         # Web server
│       └── main.go
├── internal/
│   ├── database/       # Database configuration
│   ├── handlers/       # REST API handlers
│   ├── models/         # Data models
│   └── services/       # Monitoring logic
├── web/                # React frontend
│   ├── public/
│   └── src/
│       ├── components/
│       ├── services/
│       └── App.js
├── Makefile           # Useful commands
├── go.mod             # Go dependencies
└── sites.txt          # List of sites to monitor
```

## 📋 Prerequisites

- **Go** 1.19+ installed
- **Node.js** 16+ and npm
- **Git** to clone the repository

## 🔧 Installation & Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/luacarol/website-monitoring.git
   cd website-monitoring
   ```

2. **Install Go dependencies**
   ```bash
   make deps
   ```

3. **Install Frontend dependencies**
   ```bash
   make setup-frontend
   ```

## 🚀 How to Run

### Option 1: Full Environment (Recommended)

Run the complete system with backend and frontend in separate terminals:

**Terminal 1 - Backend:**
```bash
make run-server
```
Backend server will be available at `http://localhost:8080`

**Terminal 2 - Frontend:**
```bash
make dev-frontend
```
Web dashboard will be available at `http://localhost:3000`

### Option 2: CLI Only

To use only the command-line version:
```bash
make run-cli
```

### Option 3: Production Build

**Backend Build:**
```bash
make build
./bin/monitor-server
```

**Frontend Build:**
```bash
make build-frontend
```
Static files will be generated in `web/build/`

## 📖 Using the Application

### Web Dashboard (http://localhost:3000)

1. **Dashboard**: View general statistics and status of all sites
   - Total monitored sites
   - Online/offline sites
   - Overall uptime
   - Individual status of each site

2. **Manage Sites**: Add or remove sites from monitoring
   - Click on "Manage Sites"
   - Add new sites with URL
   - Enable/disable monitoring
   - Delete sites from the list

3. **Logs**: View complete verification history
   - Filter by site and status
   - Detailed timestamps
   - HTTP status codes

### REST API (http://localhost:8080/api)

#### Available Endpoints:

**Sites:**
- `GET /api/sites` - List all sites
- `POST /api/sites` - Add new site
- `DELETE /api/sites/:id` - Remove site
- `PUT /api/sites/:id/toggle` - Enable/disable monitoring

**Logs:**
- `GET /api/logs` - Returns log history

**Statistics:**
- `GET /api/stats` - Returns general statistics

**Monitoring:**
- `POST /api/monitor/check/:id` - Check specific site
- `GET /api/monitor/status` - Monitor status

### CLI Mode

In CLI mode, you'll see an interactive menu:

```
Olá, sra. Luana
Este programa está na versão 1.25
1- Iniciar Monitoramento
2- Exibir Logs
0- Sair do Programa
```

**Options:**
- **1**: Start monitoring sites listed in `sites.txt`
- **2**: Display logs from previous checks
- **0**: Exit the program

## ⚙️ Makefile Commands

```bash
make help           # Show all available commands
make deps           # Install Go dependencies
make setup-frontend # Install frontend dependencies
make run-cli        # Run CLI version
make run-server     # Run backend server
make dev-frontend   # Start frontend in development mode
make build          # Compile binaries
make build-frontend # Build frontend for production
make clean          # Remove temporary files and database
```

## ⚙️ Configuration

### Adding Sites for Monitoring

**Via Web Interface:**
1. Access the dashboard at `http://localhost:3000`
2. Click on "Manage Sites"
3. Fill the form with the URL
4. Click "Add Site"

**Via File (for CLI):**
Edit the `sites.txt` file:
```txt
https://www.google.com
https://www.youtube.com
https://www.alura.com.br
```

### Environment Variables

```bash
PORT=8080          # Backend server port (default: 8080)
GIN_MODE=release   # Gin mode (debug/release)
```

## 📊 Log Format

Logs are stored in SQLite database and displayed in the format:

```
2026-02-15 11:25:12 - https://www.google.com - Status: 200 - Online: true
2026-02-15 11:25:15 - https://www.example.com - Status: 404 - Online: false
```

Each log entry includes:
- **Timestamp**: Exact date and time of the check
- **URL**: Monitored site
- **Status Code**: HTTP response code
- **Online**: Boolean availability indicator

## 🔍 How It Works

### Backend (Go)
1. **HTTP Server**: Gin framework serving REST API on port 8080
2. **Database**: GORM with SQLite for persistence
3. **Monitor Service**: Goroutine checking sites at regular intervals
4. **WebSocket**: Real-time connections for dashboard updates
5. **CORS**: Configured to accept frontend requests

### Frontend (React)
1. **Dashboard**: Main component with statistics cards
2. **Sites Manager**: Interface to add/remove sites
3. **Logs Viewer**: Monitoring history visualization
4. **API Service**: Axios for backend communication
5. **TailwindCSS**: Modern and responsive styling

### Monitoring Flow
1. Monitor Service executes periodic checks
2. Sends HTTP GET requests to each active site
3. Records status code and timestamp in database
4. Updates uptime statistics
5. Frontend queries API to display updated data

## 🛠️ Tech Stack

**Backend:**
- Go 1.19+
- Gin Web Framework
- GORM (ORM)
- SQLite
- Gorilla WebSocket

**Frontend:**
- React 18
- Axios
- React Router
- TailwindCSS
- Lucide React (icons)
- date-fns

## 📝 API Usage Example

```bash
# Add new site
curl -X POST http://localhost:8080/api/sites \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.github.com", "name": "GitHub"}'

# List all sites
curl http://localhost:8080/api/sites

# View statistics
curl http://localhost:8080/api/stats

# Check specific site now
curl -X POST http://localhost:8080/api/monitor/check/1
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 🎯 Roadmap

**v2.0** (Current)
- ✅ Web dashboard with React
- ✅ REST API
- ✅ Database

**v3.0** (Planned)
- ⏳ User authentication
- ⏳ Notifications
- ⏳ Multiple protocols

## 📸 Screenshots

The dashboard displays:
- Cards with main metrics (Total Sites, Online, Offline, Uptime)
- Site list with individual status and uptime
- Visual status indicators (green/red)
- Refresh buttons and quick actions

## 📈 Future Enhancements

- [ ] Email/SMS notifications for downtime
- [ ] User authentication and authorization
- [ ] Custom HTTP headers support
- [ ] Response time measurement
- [ ] Different monitoring protocols (ping, TCP)
- [ ] Report export (PDF, CSV)
- [ ] Uptime history charts
- [ ] Configurable alerts per site
- [ ] SSL certificate monitoring
- [ ] Multi-user dashboard

## 🚀 Implemented Features

- [x] Web dashboard with React
- [x] Complete REST API
- [x] SQLite database
- [x] Automatic background monitoring
- [x] WebSocket for real-time updates
- [x] Responsive interface
- [x] Site management via UI
- [x] Logs history
- [x] Uptime statistics
- [x] CLI and web server mode

## 🐛 Troubleshooting

### Port already in use
```bash
# If port 8080 or 3000 is in use, you can change it:
PORT=8081 make run-server  # Backend on port 8081

# For frontend, edit web/package.json and add:
# "start": "PORT=3001 react-scripts start"
```

### Corrupted database
```bash
make clean  # Remove database and temporary files
```

### Outdated dependencies
```bash
make deps           # Update Go dependencies
cd web && npm update  # Update npm dependencies
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**Luana Carol**
- GitHub: [@luacarol](https://github.com/luacarol)

## 🙏 Acknowledgments

- Project developed as part of the Go and React learning journey
- Inspired by the need for simple monitoring solutions
- Thanks to the Go and React communities for excellent documentation

## 📚 Useful Resources

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [React Documentation](https://react.dev/)
- [TailwindCSS](https://tailwindcss.com/)
- [GORM](https://gorm.io/)

---

⭐ **If this project was helpful to you, consider giving it a star!** ⭐

🚀 **Application running at:**
- Backend: http://localhost:8080
- Frontend: http://localhost:3000
- API Docs: http://localhost:8080/api
