import React, { useState } from 'react';

const DeleteFile = () => {
  const [bucketName, setBucketName] = useState('');
  const [fileName, setFileName] = useState('');

  const deleteFile = () => {
    fetch(`http://localhost:8080/delete?bucket=${bucketName}&filename=${fileName}`, { method: 'DELETE' })
      .then(response => response.text())
      .then(data => alert(data))
      .catch(error => console.error('Error:', error));
  };

  return (
    <div>
      <h2>Delete File</h2>
      <input
        type="text"
        value={bucketName}
        onChange={(e) => setBucketName(e.target.value)}
        placeholder="Bucket Name"
      />
      <input
        type="text"
        value={fileName}
        onChange={(e) => setFileName(e.target.value)}
        placeholder="File Name"
      />
      <button onClick={deleteFile}>Delete</button>
    </div>
  );
};

export default DeleteFile;