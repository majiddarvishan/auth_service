import React, { useState, useEffect } from "react";
import api from '../../services/api';

const Register = () => {
  const [form, setForm] = useState({ username: "", password: "", role: "" });
  const [availableRoles, setAvailableRoles] = useState([]);

  // Fetch roles from API when the component mounts
  useEffect(() => {
    async function fetchRoles() {
      try {
        const response = await api.get("/roles", {
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        });
        // Assuming response.data.roles is an array of role objects
        setAvailableRoles(response.data.roles);
      } catch (error) {
        alert("Failed to load roles");
      }
    }

    fetchRoles();
  }, []); // Empty dependency array to run only once after the component mounts

  const handleRegister = async (e) => {
    e.preventDefault();
    try {
      await api.post("/users", form);
      alert("Registered successfully!");
    } catch (error) {
      alert("Registration failed");
    }
  };

  return (
    <div className="container mt-5">
      <div className="card shadow-lg p-5" style={{ maxWidth: "500px", margin: "0 auto" }}>
        <h2 className="text-center mb-4">Create an Account</h2>
        <form onSubmit={handleRegister}>
          <div className="mb-3">
            <label htmlFor="username" className="form-label fw-bold">Username</label>
            <input
              type="text"
              id="username"
              className="form-control"
              placeholder="Enter your username"
              value={form.username}
              onChange={(e) => setForm({ ...form, username: e.target.value })}
              required
            />
          </div>

          <div className="mb-3">
            <label htmlFor="password" className="form-label fw-bold">Password</label>
            <input
              type="password"
              id="password"
              className="form-control"
              placeholder="Enter your password"
              value={form.password}
              onChange={(e) => setForm({ ...form, password: e.target.value })}
              required
            />
          </div>

          <div className="mb-3">
            <label htmlFor="role" className="form-label fw-bold">Role</label>
            <select
              id="role"
              className="form-select"
              value={form.role}
              onChange={(e) => setForm({ ...form, role: e.target.value })}
              required
            >
              <option value="" disabled>Select a role</option>
              {availableRoles.map((role) => (
                <option key={role.ID} value={role.Name}>
                  {role.Name}
                </option>
              ))}
            </select>
          </div>

          <button type="submit" className="btn btn-success btn-lg w-100">Register</button>
        </form>
      </div>
    </div>
  );
};

export default Register;
