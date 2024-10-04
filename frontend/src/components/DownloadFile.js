import React, { useState } from 'react';

const DownloadFile = () => {
  const [bucketName, setBucketName] = useState('');
  const [fileName, setFileName] = useState('');

  const downloadFile = () => {
    window.location.href = `http://localhost:8080/download?bucket=${bucketName}&filename=${fileName}`;
  };

  return (
    <div>
      <h2>Download File</h2>
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
      <button onClick={downloadFile}>Download</button>
    </div>
  );
};

export default DownloadFile;