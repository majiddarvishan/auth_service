import React, { useState } from "react";
import axios from "axios";

const NewRoleForm = () => {
  const [roleData, setRoleData] = useState({ name: "", description: "" });

  const handleCreateRole = async (e) => {
    e.preventDefault();
    try {
      // POST to the backend endpoint (adjust the URL if needed)
      const response = await axios.post(
        "https://localhost:8443/roles",
        roleData,
        {
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        }
      );
      alert(`New role "${roleData.name}" created successfully!`);
      setRoleData({ name: "", description: "" }); // Clear the form
    } catch (error) {
      console.error(error);
      alert("Failed to create role");
    }
  };

  return (
    <div className="mt-5">
      <h3>Define New Role</h3>
      <form onSubmit={handleCreateRole}>
        <div className="mb-3">
          <label className="form-label">Role Name</label>
          <input
            type="text"
            className="form-control"
            placeholder="Enter role name"
            value={roleData.name}
            onChange={(e) =>
              setRoleData({ ...roleData, name: e.target.value })
            }
            required
          />
        </div>
        <div className="mb-3">
          <label className="form-label">Description</label>
          <textarea
            className="form-control"
            placeholder="Enter role description"
            value={roleData.description}
            onChange={(e) =>
              setRoleData({ ...roleData, description: e.target.value })
            }
            required
          />
        </div>
        <button type="submit" className="btn btn-success">
          Create Role
        </button>
      </form>
    </div>
  );
};

export default NewRoleForm;
