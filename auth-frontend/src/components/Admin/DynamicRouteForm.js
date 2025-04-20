import React, { useState } from "react";
import axios from "axios";

const DynamicRouteForm = () => {
  const [dynamicRoute, setDynamicRoute] = useState({
    path: "",
    handler: "",
    method: "ANY"
  });

  // Handler for creating new dynamic routes
  const handleCreateRoute = async (e) => {
    e.preventDefault();
    try {
      await axios.post(
        "http://localhost:8080/admin/customendpoints",
        dynamicRoute,
        {
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        }
      );
      alert(`Dynamic route created successfully: ${dynamicRoute.path}`);
      setDynamicRoute({ path: "", handler: "", method: "ANY" });
    } catch (error) {
      console.error("Failed to create route", error);
      alert("Failed to create route");
    }
  };

  return (
    <div className="mt-4">
      <h3>Define New Dynamic Route</h3>
      <form onSubmit={handleCreateRoute}>
        <div className="mb-3">
          <label className="form-label">Path:</label>
          <input
            type="text"
            className="form-control"
            placeholder="Enter path (e.g., /sms/*path)"
            value={dynamicRoute.path}
            onChange={(e) =>
              setDynamicRoute({ ...dynamicRoute, path: e.target.value })
            }
            required
          />
        </div>
        <div className="mb-3">
          <label className="form-label">Handler Name:</label>
          <input
            type="text"
            className="form-control"
            placeholder="Enter handler name (e.g., SMSProxyRequest)"
            value={dynamicRoute.handler}
            onChange={(e) =>
              setDynamicRoute({ ...dynamicRoute, handler: e.target.value })
            }
            required
          />
        </div>
        <div className="mb-3">
          <label className="form-label">HTTP Method:</label>
          <select
            className="form-select"
            value={dynamicRoute.method}
            onChange={(e) =>
              setDynamicRoute({ ...dynamicRoute, method: e.target.value })
            }
          >
            <option value="ANY">ANY</option>
            <option value="GET">GET</option>
            <option value="POST">POST</option>
            <option value="PUT">PUT</option>
            <option value="DELETE">DELETE</option>
          </select>
        </div>
        <button type="submit" className="btn btn-success">Create Route</button>
      </form>
    </div>
  );
};

export default DynamicRouteForm;
