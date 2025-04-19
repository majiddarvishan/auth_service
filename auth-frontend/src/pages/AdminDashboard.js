import React, { useEffect, useState } from "react";
import axios from "axios";
import NewRoleForm from "../components/Admin/NewRoleForm";

const AdminDashboard = () => {
  // State to store users, roles and form selections.
  const [users, setUsers] = useState([]);
  const [availableRoles, setAvailableRoles] = useState([]);
  const [selectedUser, setSelectedUser] = useState("");
  const [newRole, setNewRole] = useState("");

  // Fetch the list of users and roles when the component mounts.
  useEffect(() => {
    async function fetchUsers() {
      try {
        const response = await axios.get("http://localhost:8080/admin", {
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        });
        // Expecting response.data.users to be an array of objects with at least { username, role }
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
        // Expecting response.data.roles to be an array of role objects, e.g., { ID, Name, Description }
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
      await axios.put(
        `http://localhost:8080/user/${selectedUser}/role`,
        { role: newRole },
        {
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        }
      );
      alert(`Role updated for ${selectedUser} to ${newRole}`);
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

      {/* Users List */}
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

      {/* Update User Role Section */}
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

      {/* New Role Definition Section */}
      <NewRoleForm />
    </div>
  );
};

export default AdminDashboard;
