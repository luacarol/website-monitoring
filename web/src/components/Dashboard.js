import React, { useState, useEffect } from 'react';
import { 
  Globe, 
  CheckCircle, 
  XCircle, 
  Clock, 
  TrendingUp,
  RefreshCw,
  AlertTriangle
} from 'lucide-react';
import { format } from 'date-fns';
import toast from 'react-hot-toast';

import { sitesAPI, statsAPI } from '../services/api';

const Dashboard = () => {
  const [sites, setSites] = useState([]);
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadData();
    const interval = setInterval(loadData, 30000); // Atualizar a cada 30s
    return () => clearInterval(interval);
  }, []);

  const loadData = async () => {
    try {
      const [sitesResponse, statsResponse] = await Promise.all([
        sitesAPI.getAll(),
        statsAPI.get()
      ]);
      setSites(sitesResponse.data.sites || []);
      setStats(statsResponse.data);
    } catch (error) {
      toast.error('Erro ao carregar dados');
      console.error('Erro:', error);
    } finally {
      setLoading(false);
    }
  };

  const checkSiteNow = async (siteId, siteName) => {
    try {
      toast.loading(`Verificando ${siteName}...`, { id: 'check' });
      await sitesAPI.checkNow(siteId);
      toast.success(`${siteName} verificado!`, { id: 'check' });
      // Recarregar dados após alguns segundos
      setTimeout(() => loadData(), 2000);
    } catch (error) {
      toast.error('Erro ao verificar site', { id: 'check' });
    }
  };

  const getStatusIcon = (status, isOnline) => {
    if (status === 0) {
      return <Clock className="h-5 w-5 text-gray-400" />;
    }
    return isOnline ? (
      <CheckCircle className="h-5 w-5 text-success-500" />
    ) : (
      <XCircle className="h-5 w-5 text-danger-500" />
    );
  };

  const getStatusColor = (status, isOnline) => {
    if (status === 0) return 'gray';
    return isOnline ? 'success' : 'danger';
  };

  const getUptimeColor = (uptime) => {
    if (uptime >= 95) return 'text-success-600';
    if (uptime >= 80) return 'text-warning-600';
    return 'text-danger-600';
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <RefreshCw className="h-8 w-8 animate-spin text-primary-500" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">Dashboard</h2>
          <p className="text-gray-600">
            Monitor the status of your websites in real-time
          </p>
        </div>
        <button
          onClick={loadData}
          className="btn btn-secondary"
          disabled={loading}
        >
          <RefreshCw className={`h-4 w-4 mr-2 ${loading ? 'animate-spin' : ''}`} />
          Refresh
        </button>
      </div>

      {/* Stats Cards */}
      {stats && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <div className="card">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Total Sites</p>
                <p className="text-3xl font-bold text-gray-900">{stats.total_sites}</p>
              </div>
              <Globe className="h-8 w-8 text-primary-500" />
            </div>
          </div>

          <div className="card">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Online</p>
                <p className="text-3xl font-bold text-success-600">{stats.online_sites}</p>
              </div>
              <CheckCircle className="h-8 w-8 text-success-500" />
            </div>
          </div>

          <div className="card">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Offline</p>
                <p className="text-3xl font-bold text-danger-600">{stats.offline_sites}</p>
              </div>
              <XCircle className="h-8 w-8 text-danger-500" />
            </div>
          </div>

          <div className="card">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Uptime</p>
                <p className={`text-3xl font-bold ${getUptimeColor(stats.overall_uptime)}`}>
                  {stats.overall_uptime.toFixed(1)}%
                </p>
              </div>
              <TrendingUp className="h-8 w-8 text-primary-500" />
            </div>
          </div>
        </div>
      )}

      {/* Sites Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {sites.map((site) => (
          <div key={site.id} className="card hover:shadow-md transition-shadow">
            <div className="flex items-start justify-between">
              <div className="flex items-center space-x-3 min-w-0 flex-1">
                {getStatusIcon(site.last_status, site.last_status >= 200 && site.last_status < 400)}
                <div className="min-w-0 flex-1">
                  <h3 className="text-lg font-semibold text-gray-900 truncate">
                    {site.name}
                  </h3>
                  <p className="text-sm text-gray-500 truncate">
                    {site.url}
                  </p>
                </div>
              </div>
              
              <button
                onClick={() => checkSiteNow(site.id, site.name)}
                className="btn btn-secondary btn-sm ml-3 flex-shrink-0"
                title="Check now"
              >
                <RefreshCw className="h-4 w-4" />
              </button>
            </div>

            <div className="mt-4 grid grid-cols-3 gap-4 text-sm">
              <div>
                <p className="text-gray-500">Status</p>
                <span className={`badge badge-${getStatusColor(site.last_status, site.last_status >= 200 && site.last_status < 400)}`}>
                  {site.last_status === 0 ? 'Pending' : `${site.last_status}`}
                </span>
              </div>
              
              <div>
                <p className="text-gray-500">Uptime</p>
                <p className={`font-semibold ${getUptimeColor(site.uptime)}`}>
                  {site.uptime?.toFixed(1)}%
                </p>
              </div>
              
              <div>
                <p className="text-gray-500">Last Check</p>
                <p className="text-gray-900">
                  {site.last_check && site.last_check !== '0001-01-01T00:00:00Z' 
                    ? format(new Date(site.last_check), 'HH:mm:ss')
                    : 'Never'
                  }
                </p>
              </div>
            </div>

            {!site.active && (
              <div className="mt-3 flex items-center space-x-2 text-warning-600">
                <AlertTriangle className="h-4 w-4" />
                <span className="text-sm font-medium">Monitoring Disabled</span>
              </div>
            )}
          </div>
        ))}
      </div>

      {sites.length === 0 && (
        <div className="text-center py-12">
          <Globe className="h-12 w-12 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">
            No sites configured
          </h3>
          <p className="text-gray-600 mb-4">
            Add your first website to start monitoring
          </p>
          <button
            onClick={() => window.location.href = '/sites'}
            className="btn btn-primary"
          >
            Add Website
          </button>
        </div>
      )}
    </div>
  );
};

export default Dashboard;