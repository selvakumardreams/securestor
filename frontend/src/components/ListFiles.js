import React, { useState } from 'react';

const ListFiles = () => {
  const [bucketName, setBucketName] = useState('');
  const [files, setFiles] = useState([]);

  const listFiles = () => {
    fetch(`http://localhost:8080/list?bucket=${bucketName}`)
      .then(response => response.text())
      .then(data => {
        setFiles(data.split('\n').filter(file => file));
      })
      .catch(error => console.error('Error:', error));
  };

  return (
    <div>
      <h2>List Files</h2>
      <input
        type="text"
        value={bucketName}
        onChange={(e) => setBucketName(e.target.value)}
        placeholder="Bucket Name"
      />
      <button onClick={listFiles}>List</button>
      <ul>
        {files.map(file => (
          <li key={file}>{file}</li>
        ))}
      </ul>
    </div>
  );
};

export default ListFiles;