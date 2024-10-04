import React, { useState } from 'react';

const ListFiles = () => {
  const [bucketName, setBucketName] = useState('');
  const [files, setFiles] = useState([]);
  const [showTable, setShowTable] = useState(false);

  const listFiles = () => {
    console.log(`Fetching files for bucket: ${bucketName}`);
    fetch(`http://localhost:8080/list?bucket=${bucketName}`)
      .then(response => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.text();
      })
      .then(data => {
        console.log('Response data:', data);
        const fileList = data.split('\n').filter(file => file);
        console.log('Parsed file list:', fileList);
        setFiles(fileList);
        setShowTable(true); // Show the table after fetching the files
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
      {showTable && files.length > 0 && (
        <table>
          <thead>
            <tr>
              <th>File Name</th>
            </tr>
          </thead>
          <tbody>
            {files.map(file => (
              <tr key={file}>
                <td>{file}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
};

export default ListFiles;