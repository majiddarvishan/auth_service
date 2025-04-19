import React, { useEffect, useState } from "react";
import axios from "axios";

const NewUserForm = () => {
  const [userData, setUserData] = useState({
    username: "",
    password: "",
    role: "" // will be set after roles are fetched
  });
  const [availableRoles, setAvailableRoles] = useState([]);

  // Fetch available roles when the component mounts.
  useEffect(() => {
    const fetchRoles = async () => {
      try {
        const response = await axios.get("http://localhost:8080/roles", {
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` }
        });
        // Assume response.data.roles is an array of role objects, e.g., { ID, Name, Description }
        setAvailableRoles(response.data.roles);
        // If roles exist, set the default role to the first one.
        if (response.data.roles.length > 0) {
          setUserData(prev => ({ ...prev, role: response.data.roles[0].Name }));
        }
      } catch (error) {
        console.error("Failed to fetch roles", error);
      }
    };

    fetchRoles();
  }, []);

  const handleCreateUser = async (e) => {
    e.preventDefault();
    try {
      await axios.post("http://localhost:8080/admin/users", userData, {
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` }
      });
      alert(`New user "${userData.username}" created successfully!`);
      // Reset form. If available roles exist, set role to first one.
      setUserData({
        username: "",
        password: "",
        role: availableRoles.length > 0 ? availableRoles[0].Name : ""
      });
    } catch (error) {
      console.error(error);
      alert("Failed to create user");
    }
  };

  return (
    <div className="mt-5">
      <h3>Create New User</h3>
      <form onSubmit={handleCreateUser}>
        <div className="mb-3">
          <label className="form-label">Username</label>
          <input
            type="text"
            className="form-control"
            placeholder="Enter username"
            value={userData.username}
            onChange={(e) =>
              setUserData({ ...userData, username: e.target.value })
            }
            required
          />
        </div>
        <div className="mb-3">
          <label className="form-label">Password</label>
          <input
            type="password"
            className="form-control"
            placeholder="Enter password"
            value={userData.password}
            onChange={(e) =>
              setUserData({ ...userData, password: e.target.value })
            }
            required
          />
        </div>
        <div className="mb-3">
          <label className="form-label">Role</label>
          <select
            className="form-select"
            value={userData.role}
            onChange={(e) =>
              setUserData({ ...userData, role: e.target.value })
            }
            required
          >
            {availableRoles.length === 0 && (
              <option value="">Loading roles...</option>
            )}
            {availableRoles.map((role) => (
              <option key={role.ID} value={role.Name}>
                {role.Name}
              </option>
            ))}
          </select>
        </div>
        <button type="submit" className="btn btn-success">
          Create User
        </button>
      </form>
    </div>
  );
};

export default NewUserForm;
