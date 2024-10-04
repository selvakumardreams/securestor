import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

const Login = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const navigate = useNavigate();
  
    const handleLogin = (e) => {
      e.preventDefault(); // Prevent the default form submission behavior
  
      // Dummy credentials
      const dummyUsername = 'admin';
      const dummyPassword = 'password';
  
      if (username === dummyUsername && password === dummyPassword) {
        navigate('/dashboard', { replace: true }); // Replace the current entry in the history stack
      } else {
        alert('Invalid username or password');
      }
    };
  
    return (
      <div className="login-container">
        <h2>Login</h2>
        <form onSubmit={handleLogin}>
          <input
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="Username"
          />
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Password"
          />
          <button type="submit">Login</button>
        </form>
      </div>
    );
  };

  export default Login;