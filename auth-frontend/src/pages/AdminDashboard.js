import React, { useEffect, useState } from "react";
import axios from "axios";
import NewRoleForm from "../components/Admin/NewRoleForm";
import NewUserForm from "../components/Admin/NewUserForm";

const AdminDashboard = () => {
  // State to store users, roles and form selections.
  const [users, setUsers] = useState([]);
  const [availableRoles, setAvailableRoles] = useState([]);

  // States for updating a user's role:
  const [selectedUserRole, setSelectedUserRole] = useState("");
  const [newRole, setNewRole] = useState("");

  // States for setting charge:
  const [selectedUserCharge, setSelectedUserCharge] = useState("");
  const [chargeValue, setChargeValue] = useState("");

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

  // Handler for updating a user's role.
  const handleRoleChange = async () => {
    if (!selectedUserRole || !newRole) {
      alert("Please select both a user and a new role.");
      return;
    }
    try {
      await axios.put(
        `http://localhost:8080/user/${selectedUserRole}/role`,
        { role: newRole },
        {
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        }
      );
      alert(`Role updated for ${selectedUserRole} to ${newRole}`);
      setUsers((prevUsers) =>
        prevUsers.map((user) =>
          user.username === selectedUserRole ? { ...user, role: newRole } : user
        )
      );
    } catch (error) {
      alert("Failed to update role");
    }
  };

  // Handler for setting a charge for a user.
  const handleSetCharge = async () => {
    if (!selectedUserCharge || !chargeValue) {
      alert("Please select a user and specify a charge amount.");
      return;
    }
    try {
      await axios.put(
        `http://localhost:8080/accounting/users/${selectedUserCharge}/charge`,
        { charge: parseFloat(chargeValue) },
        { headers: { Authorization: `Bearer ${localStorage.getItem("token")}` } }
      );
      alert(`Charge of ${chargeValue} set for ${selectedUserCharge}`);
      // Optionally, update the local list for real-time feedback
      setSelectedUserCharge("");
      setChargeValue("");
    } catch (error) {
      alert("Failed to set charge");
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
          onChange={(e) => setSelectedUserRole(e.target.value)}
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

      {/* Set User Charge Section */}
      <h3 className="mt-4">Set User Charge</h3>
      <div className="mb-3">
        <label className="form-label">Select User:</label>
        <select
          className="form-select"
          value={selectedUserCharge}
          onChange={(e) => setSelectedUserCharge(e.target.value)}
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
        <label className="form-label">Enter Charge Amount:</label>
        <input
          type="number"
          className="form-control"
          placeholder="Enter charge value"
          value={chargeValue}
          onChange={(e) => setChargeValue(e.target.value)}
        />
      </div>
      <button className="btn btn-primary" onClick={handleSetCharge}>
        Set Charge
      </button>

      {/* New Role Definition Section */}
      <NewRoleForm />


      {/* New User Creation Section */}
      <NewUserForm />
    </div>
  );
};

export default AdminDashboard;
