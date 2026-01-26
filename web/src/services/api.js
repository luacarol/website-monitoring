import axios from 'axios';

const API_BASE_URL = process.env.NODE_ENV === 'production' 
  ? '/api' 
  : 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
});

// Sites
export const sitesAPI = {
  getAll: () => api.get('/sites'),
  create: (data) => api.post('/sites', data),
  delete: (id) => api.delete(`/sites/${id}`),
  toggle: (id) => api.put(`/sites/${id}/toggle`),
  checkNow: (id) => api.post(`/monitor/check/${id}`),
};

// Logs
export const logsAPI = {
  getAll: (params = {}) => api.get('/logs', { params }),
};

// Stats
export const statsAPI = {
  get: () => api.get('/stats'),
};

// Monitor
export const monitorAPI = {
  getStatus: () => api.get('/monitor/status'),
};

export default api;