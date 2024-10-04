import React, { useState } from 'react';

const CreateBucket = () => {
  const [bucketName, setBucketName] = useState('');

  const createBucket = () => {
    fetch(`http://localhost:8080/create-bucket?bucket=${bucketName}`, { method: 'POST' })
      .then(response => response.text())
      .then(data => alert(data))
      .catch(error => console.error('Error:', error));
  };

  return (
    <div>
      <h2>Create Bucket</h2>
      <input
        type="text"
        value={bucketName}
        onChange={(e) => setBucketName(e.target.value)}
        placeholder="Bucket Name"
      />
      <button onClick={createBucket}>Create</button>
    </div>
  );
};

export default CreateBucket;