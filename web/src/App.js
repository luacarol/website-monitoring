import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, NavLink } from 'react-router-dom';
import { Toaster } from 'react-hot-toast';
import { 
  Monitor, 
  Activity, 
  FileText, 
  Plus,
  Wifi,
  WifiOff
} from 'lucide-react';

import Dashboard from './components/Dashboard';
import SitesManager from './components/SitesManager';
import LogsViewer from './components/LogsViewer';
import { monitorAPI } from './services/api';

function App() {
  const [monitorStatus, setMonitorStatus] = useState({ running: false });

  useEffect(() => {
    checkMonitorStatus();
    const interval = setInterval(checkMonitorStatus, 30000);
    return () => clearInterval(interval);
  }, []);

  const checkMonitorStatus = async () => {
    try {
      const response = await monitorAPI.getStatus();
      setMonitorStatus(response.data);
    } catch (error) {
      console.error('Erro ao verificar status do monitor:', error);
    }
  };

  return (
    <Router>
      <div className="min-h-screen bg-gray-50">
        <Toaster
          position="top-right"
          toastOptions={{
            duration: 4000,
            style: {
              background: '#363636',
              color: '#fff',
            },
            success: {
              duration: 3000,
            },
            error: {
              duration: 5000,
            },
          }}
        />
        
        {/* Header */}
        <header className="bg-white border-b border-gray-200 sticky top-0 z-40">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex justify-between items-center h-16">
              {/* Logo */}
              <div className="flex items-center space-x-3">
                <div className="bg-blue-600 p-2 rounded-lg">
                  <Monitor className="h-6 w-6 text-white" />
                </div>
                <div>
                  <h1 className="text-xl font-semibold text-gray-900">
                    Website Monitor
                  </h1>
                  <p className="text-xs text-gray-500">
                    Real-time monitoring
                  </p>
                </div>
              </div>

              {/* Monitor Status */}
              <div className="flex items-center space-x-2">
                <div className="flex items-center space-x-1">
                  {monitorStatus.running ? (
                    <>
                      <Wifi className="h-4 w-4 text-green-500" />
                      <span className="text-sm text-green-600 font-medium">
                        Online
                      </span>
                      <div className="h-2 w-2 bg-green-500 rounded-full animate-pulse"></div>
                    </>
                  ) : (
                    <>
                      <WifiOff className="h-4 w-4 text-red-500" />
                      <span className="text-sm text-red-600 font-medium">
                        Offline
                      </span>
                      <div className="h-2 w-2 bg-red-500 rounded-full"></div>
                    </>
                  )}
                </div>
              </div>
            </div>
          </div>
        </header>

        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="flex flex-col lg:flex-row lg:space-x-8">
            {/* Sidebar Navigation */}
            <aside className="lg:w-64 mb-8 lg:mb-0">
              <nav className="bg-white rounded-xl shadow-sm border border-gray-200 p-4">
                <ul className="space-y-2">
                  <li>
                    <NavLink
                      to="/"
                      className={({ isActive }) =>
                        `flex items-center space-x-3 px-3 py-2 rounded-lg transition-colors ${
                          isActive
                            ? 'bg-blue-100 text-blue-700 font-medium'
                            : 'text-gray-600 hover:bg-gray-100'
                        }`
                      }
                      end
                    >
                      <Activity className="h-5 w-5" />
                      <span>Dashboard</span>
                    </NavLink>
                  </li>
                  <li>
                    <NavLink
                      to="/sites"
                      className={({ isActive }) =>
                        `flex items-center space-x-3 px-3 py-2 rounded-lg transition-colors ${
                          isActive
                            ? 'bg-blue-100 text-blue-700 font-medium'
                            : 'text-gray-600 hover:bg-gray-100'
                        }`
                      }
                    >
                      <Plus className="h-5 w-5" />
                      <span>Manage Sites</span>
                    </NavLink>
                  </li>
                  <li>
                    <NavLink
                      to="/logs"
                      className={({ isActive }) =>
                        `flex items-center space-x-3 px-3 py-2 rounded-lg transition-colors ${
                          isActive
                            ? 'bg-blue-100 text-blue-700 font-medium'
                            : 'text-gray-600 hover:bg-gray-100'
                        }`
                      }
                    >
                      <FileText className="h-5 w-5" />
                      <span>Logs</span>
                    </NavLink>
                  </li>
                </ul>
              </nav>
            </aside>

            {/* Main Content */}
            <main className="flex-1 min-w-0">
              <Routes>
                <Route path="/" element={<Dashboard />} />
                <Route path="/sites" element={<SitesManager />} />
                <Route path="/logs" element={<LogsViewer />} />
              </Routes>
            </main>
          </div>
        </div>
      </div>
    </Router>
  );
}

export default App;