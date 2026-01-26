import React, { useState, useEffect } from 'react';
import { 
  Plus, 
  Trash2, 
  Globe, 
  RefreshCw,
  ExternalLink,
  CheckCircle,
  XCircle,
  Pause,
  Play
} from 'lucide-react';
import toast from 'react-hot-toast';

import { sitesAPI } from '../services/api';

const SitesManager = () => {
  const [sites, setSites] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showAddForm, setShowAddForm] = useState(false);
  const [formData, setFormData] = useState({ name: '', url: '' });
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    loadSites();
  }, []);

  const loadSites = async () => {
    try {
      const response = await sitesAPI.getAll();
      setSites(response.data.sites || []);
    } catch (error) {
      toast.error('Erro ao carregar sites');
      console.error('Erro:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!formData.name.trim() || !formData.url.trim()) {
      toast.error('Por favor, preencha todos os campos');
      return;
    }

    // Validação básica de URL
    if (!formData.url.startsWith('http://') && !formData.url.startsWith('https://')) {
      setFormData(prev => ({ ...prev, url: 'https://' + prev.url }));
    }

    setSubmitting(true);
    try {
      await sitesAPI.create(formData);
      toast.success('Site adicionado com sucesso!');
      setFormData({ name: '', url: '' });
      setShowAddForm(false);
      loadSites();
    } catch (error) {
      const errorMsg = error.response?.data?.error || 'Erro ao adicionar site';
      toast.error(errorMsg);
    } finally {
      setSubmitting(false);
    }
  };

  const handleDelete = async (id, name) => {
    if (!window.confirm(`Tem certeza que deseja remover "${name}"?`)) {
      return;
    }

    try {
      await sitesAPI.delete(id);
      toast.success('Site removido com sucesso!');
      loadSites();
    } catch (error) {
      toast.error('Erro ao remover site');
    }
  };

  const handleToggle = async (id, name, currentActive) => {
    try {
      await sitesAPI.toggle(id);
      const action = currentActive ? 'desativado' : 'ativado';
      toast.success(`${name} ${action} com sucesso!`);
      loadSites();
    } catch (error) {
      toast.error('Erro ao alterar status do site');
    }
  };

  const handleCheckNow = async (id, name) => {
    try {
      toast.loading(`Verificando ${name}...`, { id: 'check' });
      await sitesAPI.checkNow(id);
      toast.success(`${name} verificado!`, { id: 'check' });
      setTimeout(() => loadSites(), 2000);
    } catch (error) {
      toast.error('Erro ao verificar site', { id: 'check' });
    }
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
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">Manage Sites</h2>
          <p className="text-gray-600">
            Add, remove, and configure websites for monitoring
          </p>
        </div>
        <button
          onClick={() => setShowAddForm(true)}
          className="btn btn-primary"
        >
          <Plus className="h-4 w-4 mr-2" />
          Add Website
        </button>
      </div>

      {/* Add Site Form */}
      {showAddForm && (
        <div className="card">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="flex items-center justify-between">
              <h3 className="text-lg font-semibold text-gray-900">
                Add New Website
              </h3>
              <button
                type="button"
                onClick={() => {
                  setShowAddForm(false);
                  setFormData({ name: '', url: '' });
                }}
                className="text-gray-400 hover:text-gray-600"
              >
                <XCircle className="h-5 w-5" />
              </button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Website Name
                </label>
                <input
                  type="text"
                  className="input"
                  placeholder="e.g., Google"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Website URL
                </label>
                <input
                  type="url"
                  className="input"
                  placeholder="e.g., https://www.google.com"
                  value={formData.url}
                  onChange={(e) => setFormData({ ...formData, url: e.target.value })}
                  required
                />
              </div>
            </div>

            <div className="flex justify-end space-x-3">
              <button
                type="button"
                onClick={() => setShowAddForm(false)}
                className="btn btn-secondary"
                disabled={submitting}
              >
                Cancel
              </button>
              <button
                type="submit"
                className="btn btn-primary"
                disabled={submitting}
              >
                {submitting ? (
                  <>
                    <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
                    Adding...
                  </>
                ) : (
                  <>
                    <Plus className="h-4 w-4 mr-2" />
                    Add Website
                  </>
                )}
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Sites List */}
      <div className="space-y-4">
        {sites.map((site) => (
          <div key={site.id} className="card">
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-4 min-w-0 flex-1">
                <div className="min-w-0 flex-1">
                  <div className="flex items-center space-x-2">
                    <h3 className="text-lg font-semibold text-gray-900 truncate">
                      {site.name}
                    </h3>
                    {!site.active && (
                      <span className="badge badge-warning">Disabled</span>
                    )}
                  </div>
                  <div className="flex items-center space-x-2 text-sm text-gray-500">
                    <a
                      href={site.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="hover:text-primary-600 flex items-center space-x-1 truncate"
                    >
                      <span className="truncate">{site.url}</span>
                      <ExternalLink className="h-3 w-3 flex-shrink-0" />
                    </a>
                  </div>
                </div>
              </div>

              <div className="flex items-center space-x-3">
                {/* Actions */}
                <div className="flex items-center space-x-2">
                  <button
                    onClick={() => handleCheckNow(site.id, site.name)}
                    className="btn btn-secondary btn-sm"
                    title="Check now"
                  >
                    <RefreshCw className="h-4 w-4" />
                  </button>

                  <button
                    onClick={() => handleToggle(site.id, site.name, site.active)}
                    className={`btn btn-sm ${site.active ? 'btn-secondary' : 'btn-success'}`}
                    title={site.active ? 'Disable monitoring' : 'Enable monitoring'}
                  >
                    {site.active ? (
                      <Pause className="h-4 w-4" />
                    ) : (
                      <Play className="h-4 w-4" />
                    )}
                  </button>

                  <button
                    onClick={() => handleDelete(site.id, site.name)}
                    className="btn btn-danger btn-sm"
                    title="Delete site"
                  >
                    <Trash2 className="h-4 w-4" />
                  </button>
                </div>
              </div>
            </div>
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
            onClick={() => setShowAddForm(true)}
            className="btn btn-primary"
          >
            <Plus className="h-4 w-4 mr-2" />
            Add Website
          </button>
        </div>
      )}
    </div>
  );
};

export default SitesManager;