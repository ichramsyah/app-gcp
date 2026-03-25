import axios from 'axios';

const api = axios.create({
  // Fallback ke '/api' buat Production biar otomatis di-handle sama Traefik Reverse Proxy
  baseURL: import.meta.env.VITE_API_URL || '/api'
});

api.interceptors.request.use(config => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default api;
