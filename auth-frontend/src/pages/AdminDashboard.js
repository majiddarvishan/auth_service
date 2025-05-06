import React, { useEffect, useState } from "react";
import NewRoleForm from "../components/Admin/NewRoleForm";
import NewUserForm from "../components/Admin/NewUserForm";
import DynamicRouteForm from "../components/Admin/DynamicRouteForm";
import api from '../services/api';

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
        const response = await api.get("/admin", {
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
        const response = await api.get("/roles", {
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
      await api.put(
        `/users/${selectedUserRole}/role`,
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
      await api.put(
        `/accounting/users/${selectedUserCharge}/charge`,
        { charge: parseFloat(chargeValue) },
        {
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        }
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
    <div className="container py-4">
      <h2 className="mb-4 text-center">Admin Dashboard</h2>

      {/* Users List */}
      <div className="card mb-4">
        <div className="card-header">
          <h3 className="mb-0">User List</h3>
        </div>
        <div className="card-body">
          <table className="table table-bordered table-hover">
            <thead className="table-light">
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
        </div>
      </div>

      {/* Update User Role Section */}
      <div className="card mb-4">
        <div className="card-header">
          <h3 className="mb-0">Update User Role</h3>
        </div>
        <div className="card-body">
          <div className="mb-3">
            <label className="form-label">Select User:</label>
            <select
              className="form-select"
              onChange={(e) => setSelectedUserRole(e.target.value)}
              defaultValue=""
            >
              <option value="" disabled>-- Select User --</option>
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
        </div>
      </div>

      {/* Set User Charge Section */}
      <div className="card mb-4">
        <div className="card-header">
          <h3 className="mb-0">Set User Charge</h3>
        </div>
        <div className="card-body">
          <div className="mb-3">
            <label className="form-label">Select User:</label>
            <select
              className="form-select"
              value={selectedUserCharge}
              onChange={(e) => setSelectedUserCharge(e.target.value)}
              defaultValue=""
            >
              <option value="" disabled>-- Select User --</option>
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
        </div>
      </div>

      {/* New Role Definition Section */}
      <div className="card mb-4">
        <div className="card-header">
          <h3 className="mb-0">New Role Definition</h3>
        </div>
        <div className="card-body">
          <NewRoleForm />
        </div>
      </div>

      {/* New User Creation Section */}
      <div className="card mb-4">
        <div className="card-header">
          <h3 className="mb-0">New User Creation</h3>
        </div>
        <div className="card-body">
          <NewUserForm />
        </div>
      </div>

      {/* Dynamic Route Creation Section */}
      <div className="card mb-4">
        <div className="card-header">
          <h3 className="mb-0">Dynamic Route Creation</h3>
        </div>
        <div className="card-body">
          <DynamicRouteForm />
        </div>
      </div>
    </div>
  );
};

export default AdminDashboard;
