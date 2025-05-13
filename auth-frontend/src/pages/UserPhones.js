import React, { useState, useEffect } from "react";
import api from "../services/api";

const UserPhones = () => {
  const [phones, setPhones] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchPhones = async () => {
      try {
        const response = await api.get("/user/phones", {
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        });
        // expecting response.data.phones to be an array of strings
        setPhones(response.data.phones || []);
      } catch (error) {
        alert("Failed to load phone numbers");
      } finally {
        setLoading(false);
      }
    };

    fetchPhones();
  }, []);

  if (loading) {
    return (
      <div className="container mt-5 text-center">
        <div className="spinner-border text-primary" role="status" />
      </div>
    );
  }

  return (
    <div className="container mt-5">
      <h2 className="mb-4">👤 User Phone Numbers</h2>
      {phones.length === 0 ? (
        <p className="text-muted">No phone numbers found.</p>
      ) : (
        <ul className="list-group">
          {phones.map((num, idx) => (
            <li key={idx} className="list-group-item">
              {num}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default UserPhones;
