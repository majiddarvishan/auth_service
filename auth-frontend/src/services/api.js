import axios from "axios";

// Base URL configuration
const API = axios.create({
  baseURL: process.env.REACT_APP_API_BASE_URL,
});

// Attach Authorization token automatically to requests
API.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
}, (error) => Promise.reject(error));

// API functions
export const login = (data) => API.post("/login", data);
export const register = (data) => API.post("/register", data);
export const fetchAdminDashboard = () => API.get("/admin");
export const fetchAccountingRules = () => axios.get("http://localhost:8082/accounting/rules"); // Call accounting service
export const sendSMS = (data) => API.post("/sms", data);

export default API;
