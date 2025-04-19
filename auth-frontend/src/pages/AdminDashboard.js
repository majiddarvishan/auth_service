import React, { useEffect, useState } from "react";
import axios from "axios";
import NewRoleForm from "../components/Admin/NewRoleForm";

const AdminDashboard = () => {
  const [users, setUsers] = useState([]);
  const [availableRoles, setAvailableRoles] = useState([]);
  const [selectedUser, setSelectedUser] = useState("");
  const [newRole, setNewRole] = useState("");

  // Fetch users and available roles once when the component mounts.
  useEffect(() => {
    async function fetchUsers() {
      try {
        const response = await axios.get("http://localhost:8080/admin", {
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        });
        // Response should have a "users" key that is an array of objects containing at least { username, role }
        setUsers(response.data.users);
      } catch (error) {
        alert("Failed to load users");
      }
    }

    async function fetchRoles() {
      try {
        const response = await axios.get("http://localhost:8080/roles", {
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        });
        // Expecting response.data.roles as an array of role objects, e.g., { ID, Name, Description }
        setAvailableRoles(response.data.roles);
      } catch (error) {
        alert("Failed to load roles");
      }
    }

    fetchUsers();
    fetchRoles();
  }, []);

  const handleRoleChange = async () => {
    if (!selectedUser || !newRole) {
      alert("Please select both a user and a new role.");
      return;
    }
    try {
      // Update the user's role by sending a PUT request to the backend.
      await axios.put(
        `http://localhost:8080/user/${selectedUser}/role`,
        { role: newRole },
        {
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        }
      );
      alert(`Role updated for ${selectedUser} to ${newRole}`);
      // Update the local user list optionally or refresh the data.
      setUsers((prevUsers) =>
        prevUsers.map((user) =>
          user.username === selectedUser ? { ...user, role: newRole } : user
        )
      );
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
      <div className="mb-3">
        <label className="form-label">Select User:</label>
        <select
          className="form-select"
          onChange={(e) => setSelectedUser(e.target.value)}
          defaultValue=""
        >
          <option value="" disabled>
            -- Select User --
          </option>
          {users.map((user) => (
            <option key={user.username} value={user.username}>
              {user.username}
            </option>
          ))}
        </select>
      </div>

      <div className="mb-3">
        <label className="form-label">Select New Role:</label>
        <select
          className="form-select"
          value={newRole}
          onChange={(e) => setNewRole(e.target.value)}
        >
          <option value="">-- Select Role --</option>
          {availableRoles.map((role) => (
            <option key={role.ID} value={role.Name}>
              {role.Name}
            </option>
          ))}
        </select>
      </div>

      <button className="btn btn-primary" onClick={handleRoleChange}>
        Update Role
      </button>

      <div className="container">
        <h2>Admin Dashboard</h2>
        <NewRoleForm />
      </div>
    </div>
  );
};

export default AdminDashboard;
