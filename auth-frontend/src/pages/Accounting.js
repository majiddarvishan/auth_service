import React, { useEffect, useState } from "react";
import axios from "axios";

const Accounting = () => {
  const [rules, setRules] = useState([]);

  useEffect(() => {
    async function fetchAccountingRules() {
      try {
        const response = await axios.get("http://localhost:8082/accounting/rules");
        setRules(response.data.rules);
      } catch (error) {
        alert("Failed to load accounting rules");
      }
    }
    fetchAccountingRules();
  }, []);

  return (
    <div className="container">
      <h2>Accounting Rules</h2>
      <ul>
        {rules.map((rule) => (
          <li key={rule.endpoint}>{rule.endpoint} - ${rule.charge}</li>
        ))}
      </ul>
    </div>
  );
};

export default Accounting;
