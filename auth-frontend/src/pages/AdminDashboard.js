import React, { useEffect, useState } from "react";
import NewRoleForm from "../components/Admin/NewRoleForm";
import NewUserForm from "../components/Admin/NewUserForm";
import DynamicRouteForm from "../components/Admin/DynamicRouteForm";
import api from '../services/api';
import 'bootstrap/dist/js/bootstrap.bundle.min.js'; // Import Bootstrap JS
import { Tab } from 'bootstrap'; // Import Tab from bootstrap

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

    // Initialize Bootstrap tabs
    const tabEls = document.querySelectorAll('button[data-bs-toggle="tab"]');
    tabEls.forEach(tabEl => {
      new Tab(tabEl);
    });
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

      <ul className="nav nav-tabs mb-4" id="adminTabs" role="tablist">
        <li className="nav-item" role="presentation">
          <button className="nav-link active" id="users-tab" data-bs-toggle="tab" data-bs-target="#users" type="button" role="tab" aria-controls="users" aria-selected="true">Users</button>
        </li>
        <li className="nav-item" role="presentation">
          <button className="nav-link" id="new-user-tab" data-bs-toggle="tab" data-bs-target="#new-user" type="button" role="tab" aria-controls="new-user" aria-selected="false">New User</button>
        </li>
        <li className="nav-item" role="presentation">
          <button className="nav-link" id="set-charge-tab" data-bs-toggle="tab" data-bs-target="#set-charge" type="button" role="tab" aria-controls="set-charge" aria-selected="false">Set Charge</button>
        </li>
        <li className="nav-item" role="presentation">
          <button className="nav-link" id="update-role-tab" data-bs-toggle="tab" data-bs-target="#update-role" type="button" role="tab" aria-controls="update-role" aria-selected="false">Update Role</button>
        </li>
        <li className="nav-item" role="presentation">
          <button className="nav-link" id="new-role-tab" data-bs-toggle="tab" data-bs-target="#new-role" type="button" role="tab" aria-controls="new-role" aria-selected="false">New Role</button>
        </li>
        <li className="nav-item" role="presentation">
          <button className="nav-link" id="dynamic-route-tab" data-bs-toggle="tab" data-bs-target="#dynamic-route" type="button" role="tab" aria-controls="dynamic-route" aria-selected="false">Dynamic Route</button>
        </li>
      </ul>

      <div className="tab-content" id="adminTabsContent">
        {/* Users List Tab */}
        <div className="tab-pane fade show active" id="users" role="tabpanel" aria-labelledby="users-tab">
          <div className="card">
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
        </div>

        {/* Update User Role Tab */}
        <div className="tab-pane fade" id="update-role" role="tabpanel" aria-labelledby="update-role-tab">
          <div className="card">
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
        </div>

        {/* Set User Charge Tab */}
        <div className="tab-pane fade" id="set-charge" role="tabpanel" aria-labelledby="set-charge-tab">
          <div className="card">
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
        </div>

        {/* New Role Definition Tab */}
        <div className="tab-pane fade" id="new-role" role="tabpanel" aria-labelledby="new-role-tab">
          <div className="card">
            <div className="card-header">
              <h3 className="mb-0">Define New Role</h3>
            </div>
            <div className="card-body">
              <NewRoleForm />
            </div>
          </div>
        </div>
       
        {/* New User Creation Tab */}
        <div className="tab-pane fade" id="new-user" role="tabpanel" aria-labelledby="new-user-tab">
          <div className="card">
            <div className="card-header">
              <h3 className="mb-0">Create New User</h3>
            </div>
            <div className="card-body">
              <NewUserForm />
            </div>
          </div>
        </div>

        {/* Dynamic Route Creation Tab */}
        <div className="tab-pane fade" id="dynamic-route" role="tabpanel" aria-labelledby="dynamic-route-tab">
          <div className="card">
            <div className="card-header">
              <h3 className="mb-0">Define New Dynamic Route</h3>
            </div>
            <div className="card-body">
              <DynamicRouteForm />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AdminDashboard;
