import React from 'react';

interface ControlPanelProps {
  onGetRequest: () => void;
  requestServerMap: { requestId: string; serverIp: string }[];
  serverRequestCounts: { ip: string; count: number }[];
}

const ControlPanel: React.FC<ControlPanelProps> = ({
  onGetRequest,
  requestServerMap,
  serverRequestCounts,
}) => {
  return (
    <div className="control-panel">
      <button onClick={onGetRequest}>Get</button>

      <h3>Request-Server Mapping</h3>
      <table>
        <thead>
          <tr>
            <th>Request ID</th>
            <th>Server IP</th>
          </tr>
        </thead>
        <tbody>
          {requestServerMap.map((entry) => (
            <tr key={entry.requestId}>
              <td>{entry.requestId}</td>
              <td>{entry.serverIp}</td>
            </tr>
          ))}
        </tbody>
      </table>

      <h3>Server Request Counts</h3>
      <table>
        <thead>
          <tr>
            <th>Server IP</th>
            <th>Request Count</th>
          </tr>
        </thead>
        <tbody>
          {serverRequestCounts.map((entry) => (
            <tr key={entry.ip}>
              <td>{entry.ip}</td>
              <td>{entry.count}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default ControlPanel;
