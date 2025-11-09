'use client';

import { useState } from 'react';

export default function TestPage() {
  const [result, setResult] = useState('');

  const testConnection = async () => {
    try {
      const response = await fetch('http://localhost:9090', {
        method: 'GET',
        mode: 'cors',
      });
      setResult(`Success! Status: ${response.status}`);
    } catch (error: any) {
      setResult(`Error: ${error.message}`);
    }
  };

  return (
    <div className="p-8">
      <h1 className="text-2xl mb-4">Test Backend Connection</h1>
      <button 
        onClick={testConnection}
        className="px-4 py-2 bg-blue-500 text-white rounded"
      >
        Test Connection
      </button>
      <p className="mt-4">{result}</p>
    </div>
  );
}
