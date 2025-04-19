import React, { useEffect, useState } from "react";
import axios from "axios";

const AdminDashboard = () => {
  const [users, setUsers] = useState([]);
  const [selectedUser, setSelectedUser] = useState("");
  const [newRole, setNewRole] = useState("user");

  useEffect(() => {
    async function fetchDashboardData() {
      try {
        const response = await axios.get("http://localhost:8080/admin", {
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        });
        setUsers(response.data.users);
      } catch (error) {
        alert("Failed to load dashboard data");
      }
    }
    fetchDashboardData();
  }, []);

  const handleRoleChange = async () => {
    try {
      await axios.put(`http://localhost:8080/user/${selectedUser}/role`, { role: newRole }, {
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
      });
      alert(`Role updated for ${selectedUser} to ${newRole}`);
    } catch (error) {
      alert("Failed to update role");
    }
  };

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

      <h3>Update User Role</h3>
      <div>
        <select onChange={(e) => setSelectedUser(e.target.value)}>
          <option value="">Select User</option>
          {users.map((user) => (
            <option key={user.username} value={user.username}>
              {user.username}
            </option>
          ))}
        </select>

        <select value={newRole} onChange={(e) => setNewRole(e.target.value)}>
          <option value="user">User</option>
          <option value="admin">Admin</option>
        </select>

        <button className="btn btn-primary" onClick={handleRoleChange}>
          Update Role
        </button>
      </div>
    </div>
  );
};

export default AdminDashboard;
