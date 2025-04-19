import React, { useEffect, useState } from "react";
import axios from "axios";

const AdminDashboard = () => {
  const [users, setUsers] = useState([]);

  useEffect(() => {
    async function fetchUsers() {
      try {
        const response = await axios.get("http://localhost:8080/admin", {
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        });
        setUsers(response.data.users);
      } catch (error) {
        alert("Failed to load users");
      }
    }
    fetchUsers();
  }, []);

  return (
    <div className="container">
      <h2>Admin Dashboard</h2>
      <ul>
        {users.map((user) => (
          <li key={user.username}>{user.username} - {user.role}</li>
        ))}
      </ul>
    </div>
  );
};

export default AdminDashboard;
