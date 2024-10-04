import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import './App.css';
import CreateBucket from './components/CreateBucket';
import UploadFile from './components/UploadFile';
import DownloadFile from './components/DownloadFile';
import ListFiles from './components/ListFiles';
import DeleteFile from './components/DeleteFile';
import HealthCheck from './components/HealthCheck';
import Login from './components/Login';

function App() {
  return (
    <Router>
      <div className="App">
        <Routes>
          <Route path="/" element={<Login />} />
          <Route path="/dashboard" element={<Dashboard />} />
        </Routes>
      </div>
    </Router>
  );
}

const Dashboard = () => (
  <div>
    <h1>Bluenoise</h1>
    <CreateBucket />
    <UploadFile />
    <DownloadFile />
    <ListFiles />
    <DeleteFile />
    <HealthCheck />
  </div>
);

export default App;