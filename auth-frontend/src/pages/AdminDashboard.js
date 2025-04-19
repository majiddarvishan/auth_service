import React, { useEffect, useState } from "react";
import axios from "axios";

const AdminDashboard = () => {
  const [users, setUsers] = useState([]);
  const [rules, setRules] = useState([]);

  useEffect(() => {
    async function fetchDashboardData() {
      try {
        const response = await axios.get("http://localhost:8080/admin", {
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        });
        setUsers(response.data.users);
        setRules(response.data.rules);
      } catch (error) {
        alert("Failed to load dashboard data");
      }
    }
    fetchDashboardData();
  }, []);

  return (
    <div className="container">
      <h2>Admin Dashboard</h2>

      <h3>User List</h3>
      <table className="table table-bordered">
        <thead>
          <tr>
            <th>Username</th>
            <th>Role</th>
          </tr>
        </thead>
        <tbody>
          {users.map((user) => (
            <tr key={user.username}>
              <td>{user.username}</td>
              <td>{user.role}</td>
            </tr>
          ))}
        </tbody>
      </table>

      <h3>Accounting Rules</h3>
      <table className="table table-striped">
        <thead>
          <tr>
            <th>Endpoint</th>
            <th>Charge ($)</th>
          </tr>
        </thead>
        <tbody>
          {rules.map((rule) => (
            <tr key={rule.endpoint}>
              <td>{rule.Endpoint}</td>
              <td>{rule.Charge}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default AdminDashboard;
