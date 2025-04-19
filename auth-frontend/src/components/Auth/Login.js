import React, { useState } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { jwtDecode } from "jwt-decode";

const Login = () => {
  const [form, setForm] = useState({ username: "", password: "" });
  const navigate = useNavigate();

  const handleLogin = async (e) => {
    e.preventDefault();
    try {
      const response = await axios.post("http://localhost:8080/login", form);
      const token = response.data.token;
      localStorage.setItem("token", token);

      // Decode the token to extract role information
      const decodedToken = jwtDecode(token);
      // Assume token includes a claim "role" (e.g., { ..., role: "admin", ... })
      if (decodedToken.role === "admin") {
        navigate("/admin"); // route for admin dashboard
      } else {
        navigate("/sms"); // route for non-admin SMS page
      }
    } catch (error) {
      alert("Login failed");
    }
  };

  return (
    <div className="container mt-5">
      <h2>Login</h2>
      <form onSubmit={handleLogin}>
        <input
          type="text"
          placeholder="Username"
          value={form.username}
          onChange={(e) => setForm({ ...form, username: e.target.value })}
        />
        <input
          type="password"
          placeholder="Password"
          value={form.password}
          onChange={(e) => setForm({ ...form, password: e.target.value })}
        />
        <button type="submit" className="btn btn-primary">Login</button>
      </form>
    </div>
  );
};

export default Login;
