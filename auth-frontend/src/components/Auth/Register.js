import React, { useState } from "react";
import axios from "axios";

const Register = () => {
  const [form, setForm] = useState({ username: "", password: "", role: "user" });

  const handleRegister = async (e) => {
    e.preventDefault();
    try {
      await axios.post("http://localhost:8080/register", form);
      alert("Registered successfully!");
    } catch (error) {
      alert("Registration failed");
    }
  };

  return (
    <div className="container mt-5">
      <h2>Register</h2>
      <form onSubmit={handleRegister}>
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

        {/* Role Dropdown */}
        <select
          value={form.role}
          onChange={(e) => setForm({ ...form, role: e.target.value })}
        >
          <option value="user">User</option>
          <option value="admin">Admin</option>
        </select>

        <button type="submit" className="btn btn-success">Register</button>
      </form>
    </div>
  );
};

export default Register;
