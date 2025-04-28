import React, { useState } from "react";
import api from '../../services/api';

const DynamicRouteForm = () => {
  const [dynamicRoute, setDynamicRoute] = useState({
    path: "",
    method: "ANY",
    endpoint: "",
    needAccounting: false, // new flag to indicate if accounting check is required
  });

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
      setDynamicRoute({
        path: "",
        method: "ANY",
        endpoint: "",
        needAccounting: false,
      });
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
        <div className="mb-3">
          <label className="form-label">Endpoint:</label>
          <input
            type="text"
            className="form-control"
            placeholder="Enter target endpoint (e.g., https://api.external.com)"
            value={dynamicRoute.endpoint}
            onChange={(e) =>
              setDynamicRoute({ ...dynamicRoute, endpoint: e.target.value })
            }
            required
          />
        </div>
        <div className="mb-3 form-check">
          <input
            type="checkbox"
            className="form-check-input"
            id="needAccounting"
            checked={dynamicRoute.needAccounting}
            onChange={(e) =>
              setDynamicRoute({
                ...dynamicRoute,
                needAccounting: e.target.checked,
              })
            }
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
