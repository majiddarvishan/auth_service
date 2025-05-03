import React, { useState } from "react";
import api from '../../services/api';

const DynamicRouteForm = () => {
  const [dynamicRoute, setDynamicRoute] = useState({
    path: "",
    method: "ANY",
    endpoints: [],   // array of target endpoints
    needAccounting: false, // flag to indicate if accounting check is required
  });
  const [newEndpoint, setNewEndpoint] = useState(""); // temp input for adding endpoints

  // Add a new endpoint to the array
  const handleAddEndpoint = () => {
    const url = newEndpoint.trim();
    if (url && !dynamicRoute.endpoints.includes(url)) {
      setDynamicRoute({
        ...dynamicRoute,
        endpoints: [...dynamicRoute.endpoints, url],
      });
      setNewEndpoint("");
    }
  };

  // Remove an endpoint by index
  const handleRemoveEndpoint = (index) => {
    setDynamicRoute({
      ...dynamicRoute,
      endpoints: dynamicRoute.endpoints.filter((_, i) => i !== index),
    });
  };

  const handleCreateRoute = async (e) => {
    e.preventDefault();
    try {
      await api.post(
        "/admin/customendpoints",
        dynamicRoute,
        {
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        }
      );
      alert(`Dynamic route created successfully: ${dynamicRoute.path}`);
      // reset form
      setDynamicRoute({ path: "", method: "ANY", endpoints: [], needAccounting: false });
      setNewEndpoint("");
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
            placeholder="Enter path (e.g., /sms)"
            value={dynamicRoute.path}
            onChange={(e) => setDynamicRoute({ ...dynamicRoute, path: e.target.value })}
            required
          />
        </div>

        <div className="mb-3">
          <label className="form-label">HTTP Method:</label>
          <select
            className="form-select"
            value={dynamicRoute.method}
            onChange={(e) => setDynamicRoute({ ...dynamicRoute, method: e.target.value })}
          >
            <option value="ANY">ANY</option>
            <option value="GET">GET</option>
            <option value="POST">POST</option>
            <option value="PUT">PUT</option>
            <option value="DELETE">DELETE</option>
          </select>
        </div>

        <div className="mb-3">
          <label className="form-label">Add Endpoint:</label>
          <div className="d-flex">
            <input
              type="text"
              className="form-control"
              placeholder="https://api.external.com"
              value={newEndpoint}
              onChange={(e) => setNewEndpoint(e.target.value)}
            />
            <button
              type="button"
              className="btn btn-secondary ms-2"
              onClick={handleAddEndpoint}
            >
              Add
            </button>
          </div>
          {dynamicRoute.endpoints.length > 0 && (
            <ul className="list-group mt-2">
              {dynamicRoute.endpoints.map((url, idx) => (
                <li key={idx} className="list-group-item d-flex justify-content-between align-items-center">
                  {url}
                  <button
                    type="button"
                    className="btn btn-sm btn-danger"
                    onClick={() => handleRemoveEndpoint(idx)}
                  >
                    Remove
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>

        <div className="mb-3 form-check">
          <input
            type="checkbox"
            className="form-check-input"
            id="needAccounting"
            checked={dynamicRoute.needAccounting}
            onChange={(e) => setDynamicRoute({ ...dynamicRoute, needAccounting: e.target.checked })}
          />
          <label className="form-check-label" htmlFor="needAccounting">
            Check Accounting
          </label>
        </div>

        <button type="submit" className="btn btn-success">
          Create Route
        </button>
      </form>
    </div>
  );
};

export default DynamicRouteForm;
