import React, { useState } from 'react';
import Ring from './components/Ring';
// import ControlPanel from './components/ControlPanel';
import axios from 'axios';
import './App.css';

interface Request {
  position: number;
  id: string;
}

interface Server {
  position: number;
  ip: string;
}

interface ServerCount {
  ip: string;
  count: number;
}

function App() {
  const [requests, setRequests] = useState<Request[]>([]);
  const [servers, setServers] = useState<Server[]>([]);
  const [requestServerMap, setRequestServerMap] = useState<
    { requestId: string; serverIp: string }[]
  >([]);
  const [serverRequestCounts, setServerRequestCounts] = useState<ServerCount[]>([]);

  const handleGetRequest = async () => {
    try {
      // Call your backend (replace the URL with your actual backend endpoint)
      const res = await axios.get('http://localhost:8085');
      const { clientIP, node, position } = res.data;

      const newRequest: Request = {
        position,
        id: clientIP,
      };

      const updatedRequests = [...requests, newRequest];
      setRequests(updatedRequests);

      // Update request-server map
      setRequestServerMap((prev) => [...prev, { requestId: clientIP, serverIp: node }]);

      // Update server request counts
      setServerRequestCounts((prevCounts) => {
        const existing = prevCounts.find((entry) => entry.ip === node);
        if (existing) {
          return prevCounts.map((entry) =>
            entry.ip === node ? { ...entry, count: entry.count + 1 } : entry
          );
        } else {
          return [...prevCounts, { ip: node, count: 1 }];
        }
      });

      // Optionally update known servers if needed
      if (!servers.find((s) => s.ip === node)) {
        setServers((prev) => [...prev, { ip: node, position }]);
      }
    } catch (err) {
      console.error('Error fetching server info:', err);
    }
  };

  return (
    <div className="App">
      <h1>Consistent Hash Ring Visualization</h1>
      <Ring  />
      {/* <ControlPanel
        onGetRequest={handleGetRequest}
        requestServerMap={requestServerMap}
        serverRequestCounts={serverRequestCounts}
      /> */}
    </div>
  );
}

export default App;
