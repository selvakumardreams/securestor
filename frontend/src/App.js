import React from 'react';
import './App.css';
import CreateBucket from './components/CreateBucket';
import UploadFile from './components/UploadFile';
import DownloadFile from './components/DownloadFile';
import ListFiles from './components/ListFiles';
import DeleteFile from './components/DeleteFile';
import HealthCheck from './components/HealthCheck';

function App() {
  return (
    <div className="App">
      <h1>Bluenoise</h1>
      <CreateBucket />
      <UploadFile />
      <DownloadFile />
      <ListFiles />
      <DeleteFile />
      <HealthCheck />
    </div>
  );
}

export default App;