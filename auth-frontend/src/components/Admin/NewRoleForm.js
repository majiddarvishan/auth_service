import React, { useState } from "react";
import axios from "axios";

const NewRoleForm = () => {
  const [roleData, setRoleData] = useState({ name: "", description: "" });

  const handleCreateRole = async (e) => {
    e.preventDefault();
    try {
      await axios.post(
        "http://localhost:8080/roles",
        roleData,
        {
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        }
      );
      alert(`New role "${roleData.name}" created successfully!`);
      setRoleData({ name: "", description: "" }); // Reset form
    } catch (error) {
      alert("Failed to create role");
    }
  };

  return (
    <div className="container mt-4">
      <h3>Create a New Role</h3>
      <form onSubmit={handleCreateRole}>
        <div className="form-group">
          <label>Role Name</label>
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
        <div className="form-group mt-2">
          <label>Description</label>
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
        <button type="submit" className="btn btn-success mt-3">
          Create Role
        </button>
      </form>
    </div>
  );
};

export default NewRoleForm;
