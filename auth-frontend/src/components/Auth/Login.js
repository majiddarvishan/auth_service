import React, { useState, useEffect } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { jwtDecode } from "jwt-decode";

// Determine API base URL, supporting HTTP or HTTPS (self-signed backend on port 8443 by default)
const defaultHttp = "http://localhost:8080";
const defaultHttps = "https://localhost:8443"; // adjust port if needed
const API_BASE = process.env.REACT_APP_API_BASE_URL
  || (window.location.protocol === "https:" ? defaultHttps : defaultHttp);

const Login = () => {
  const [form, setForm] = useState({ username: "", password: "" });
  const [captchaData, setCaptchaData] = useState({ id: null, image: null });
  const [captchaInput, setCaptchaInput] = useState("");
  const navigate = useNavigate();

  // Optionally set axios base URL
  useEffect(() => {
    axios.defaults.baseURL = API_BASE;
  }, []);

  // Fetch a new captcha from backend
  const loadCaptcha = async () => {
    try {
      const response = await axios.get("/captcha/new");
      // Expecting { id: string, image: base64String }
      setCaptchaData(response.data);
      setCaptchaInput("");
    } catch (err) {
      console.error("Failed to load captcha", err);
    }
  };

  useEffect(() => {
    loadCaptcha();
  }, []);

  const handleLogin = async (e) => {
    e.preventDefault();

    if (!captchaInput) {
      alert("Please enter the captcha");
      return;
    }

    try {
      const payload = {
        ...form,
        captchaId: captchaData.id,
        captchaText: captchaInput,
      };

      const response = await axios.post("/login", payload);
      const { token } = response.data;
      localStorage.setItem("token", token);

      const decodedToken = jwtDecode(token);
      if (decodedToken.role === "admin") {
        navigate("/admin");
      } else {
        navigate("/sms");
      }
    } catch (error) {
      alert(error.response?.data?.message || "Login failed");
      // refresh captcha on failure
      loadCaptcha();
    }
  };

  return (
    <div className="container mt-5" style={{ maxWidth: '400px' }}>
      <h2>Login</h2>
      <form onSubmit={handleLogin}>
        <input
          type="text"
          placeholder="Username"
          value={form.username}
          onChange={(e) => setForm({ ...form, username: e.target.value })}
          required
          className="form-control mb-3"
        />
        <input
          type="password"
          placeholder="Password"
          value={form.password}
          onChange={(e) => setForm({ ...form, password: e.target.value })}
          required
          className="form-control mb-3"
        />

        {captchaData.image && (
          <div className="mb-3 text-center">
            <img
              src={`data:image/png;base64,${captchaData.image}`}
              alt="captcha"
              style={{ cursor: 'pointer' }}
              onClick={loadCaptcha}
              title="Click to refresh"
            />
          </div>
        )}

        <input
          type="text"
          placeholder="Enter Captcha"
          value={captchaInput}
          onChange={(e) => setCaptchaInput(e.target.value)}
          required
          className="form-control mb-3"
        />

        <button type="submit" className="btn btn-primary w-100">
          Login
        </button>
      </form>
    </div>
  );
};

export default Login;
