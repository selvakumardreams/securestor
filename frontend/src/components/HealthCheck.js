import React, { useState } from 'react';

const HealthCheck = () => {
  const [healthStatus, setHealthStatus] = useState('');

  const checkHealth = () => {
    fetch('http://localhost:8080/health')
      .then(response => response.text())
      .then(data => {
        setHealthStatus(data);
      })
      .catch(error => console.error('Error:', error));
  };

  return (
    <div>
      <h2>Health Check</h2>
      <button onClick={checkHealth}>Check</button>
      <p>{healthStatus}</p>
    </div>
  );
};

export default HealthCheck;