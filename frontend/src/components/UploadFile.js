import React, { useState } from 'react';

const UploadFile = () => {
  const [bucketName, setBucketName] = useState('');
  const [file, setFile] = useState(null);

  const handleFileUpload = (event) => {
    setFile(event.target.files[0]);
  };

  const uploadFile = () => {
    const formData = new FormData();
    formData.append('file', file);

    fetch(`http://localhost:8080/upload?bucket=${bucketName}`, {
      method: 'POST',
      body: formData
    })
      .then(response => response.text())
      .then(data => alert(data))
      .catch(error => console.error('Error:', error));
  };

  return (
    <div>
      <h2>Upload File</h2>
      <input
        type="text"
        value={bucketName}
        onChange={(e) => setBucketName(e.target.value)}
        placeholder="Bucket Name"
      />
      <input type="file" onChange={handleFileUpload} />
      <button onClick={uploadFile}>Upload</button>
    </div>
  );
};

export default UploadFile;