import React, { useState, useEffect } from 'react';
import { 
  FileText, 
  RefreshCw, 
  Filter, 
  Download,
  CheckCircle,
  XCircle,
  Clock
} from 'lucide-react';
import { format } from 'date-fns';
import toast from 'react-hot-toast';

import { logsAPI, sitesAPI } from '../services/api';

const LogsViewer = () => {
  const [logs, setLogs] = useState([]);
  const [sites, setSites] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadSites();
    loadLogs();
  }, []);

  const loadSites = async () => {
    try {
      const response = await sitesAPI.getAll();
      setSites(response.data.sites || []);
    } catch (error) {
      console.error('Erro ao carregar sites:', error);
    }
  };

  const loadLogs = async () => {
    setLoading(true);
    try {
      const response = await logsAPI.getAll({ limit: 25 });
      setLogs(response.data.logs || []);
    } catch (error) {
      toast.error('Erro ao carregar logs');
      console.error('Erro:', error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusIcon = (isOnline, statusCode) => {
    if (statusCode === 0) {
      return <Clock className="h-4 w-4 text-gray-400" />;
    }
    return isOnline ? (
      <CheckCircle className="h-4 w-4 text-success-500" />
    ) : (
      <XCircle className="h-4 w-4 text-danger-500" />
    );
  };

  const getStatusBadge = (isOnline, statusCode) => {
    if (statusCode === 0) {
      return <span className="badge badge-warning">Pending</span>;
    }
    return (
      <span className={`badge ${isOnline ? 'badge-success' : 'badge-danger'}`}>
        {statusCode}
      </span>
    );
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">Monitoring Logs</h2>
          <p className="text-gray-600">
            View detailed monitoring history and statistics
          </p>
        </div>
        
        <button
          onClick={loadLogs}
          className="btn btn-secondary"
          disabled={loading}
        >
          <RefreshCw className={`h-4 w-4 mr-2 ${loading ? 'animate-spin' : ''}`} />
          Refresh
        </button>
      </div>

      {/* Logs Table */}
      <div className="card overflow-hidden">
        {loading ? (
          <div className="flex items-center justify-center py-8">
            <RefreshCw className="h-8 w-8 animate-spin text-primary-500" />
          </div>
        ) : logs.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Status
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Website
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Response Time
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Checked At
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {logs.map((log) => (
                  <tr key={log.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center space-x-2">
                        {getStatusIcon(log.is_online, log.status_code)}
                        {getStatusBadge(log.is_online, log.status_code)}
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <div>
                        <div className="font-medium text-gray-900">
                          {log.site.name}
                        </div>
                        <div className="text-sm text-gray-500 truncate max-w-xs">
                          {log.site.url}
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className="text-sm font-medium text-gray-900">
                        {log.response_time}ms
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {format(new Date(log.checked_at), 'MMM dd, HH:mm:ss')}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <div className="text-center py-12">
            <FileText className="h-12 w-12 text-gray-400 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-gray-900 mb-2">
              No logs found
            </h3>
            <p className="text-gray-600">
              Add some websites to start collecting monitoring data
            </p>
          </div>
        )}
      </div>
    </div>
  );
};

export default LogsViewer;