import React, { useState } from "react";
import { motion } from "framer-motion";
import axios from "axios";

const TOTAL_POSITIONS = 64; 
const CIRCLE_RADIUS = 300;

const Ring = () => {
  const [servers, setServers] = useState([]);
  const [keys, setKeys] = useState([]);
  const [disableAddServer, setDisableAddServer] = useState(false);

  const getCoordinates = (index) => {
    const angle = (2 * Math.PI * index) / TOTAL_POSITIONS;
    const x = CIRCLE_RADIUS * Math.cos(angle);
    const y = CIRCLE_RADIUS * Math.sin(angle);
    return { x, y };
  };

  const addServer = async () => {
    setDisableAddServer(true);
    try {
      const res = await axios.get("http://localhost:8085/init");
      const position = res.data.hash;
      setServers([...servers, { id: servers.length, position }]);
    } catch (err) {
      console.error("Error adding server:", err);
      setDisableAddServer(false);
    }
  };

  const addKey = async () => {
    try {
      const res = await axios.get("http://localhost:8085/");
      const position = res.data.hash;
      setKeys([...keys, { id: keys.length, position }]);
    } catch (err) {
      console.error("Error adding key:", err);
    }
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen">
      <h1 className="text-3xl font-bold mb-6">Consistent Hashing Ring</h1>
      <div className="flex gap-4 mb-4">
        <button
          onClick={addServer}
          disabled={disableAddServer}
          className={`px-4 py-2 rounded bg-green-600 text-white font-semibold ${
            disableAddServer ? "opacity-50 cursor-not-allowed" : "hover:bg-green-700"
          }`}
        >
          Add Server
        </button>
        <button
          onClick={addKey}
          className="px-4 py-2 rounded bg-blue-600 text-white font-semibold hover:bg-blue-700"
        >
          Add Key
        </button>
      </div>

      <svg width={700} height={700}>
        <g transform={`translate(${350}, ${350})`}>
          {/* Draw all 127 segments */}
          {[...Array(TOTAL_POSITIONS)].map((_, index) => {
            const { x, y } = getCoordinates(index);
            return (
              <g key={`segment-${index}`} transform={`translate(${x}, ${y})`}>
                <circle r={14} fill="#f0f0f0" stroke="#ccc" strokeWidth={2} />
                <text
                  x={0}
                  y={4}
                  textAnchor="middle"
                  fontSize={10}
                  fontWeight="bold"
                  fill="#000"
                >
                  {index}
                </text>
              </g>
            );
          })}

          {/* Draw servers */}
          {servers.map((server) => {
            const { x, y } = getCoordinates(server.position);
            return (
              <motion.image
                key={`server-${server.id}`}
                href="https://cdn-icons-png.flaticon.com/512/3665/3665923.png"
                width={30}
                height={30}
                x={x - 15}
                y={y - 15}
                initial={{ scale: 0 }}
                animate={{ scale: 1 }}
                transition={{ type: "spring", stiffness: 200 }}
              />
            );
          })}

          {/* Draw keys */}
          {keys.map((key) => {
            const { x, y } = getCoordinates(key.position);
            return (
              <motion.image
                key={`key-${key.id}`}
                href="https://cdn-icons-png.flaticon.com/512/2910/2910791.png"
                width={25}
                height={25}
                x={x - 12.5}
                y={y - 12.5}
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.5 }}
              />
            );
          })}
        </g>
      </svg>
    </div>
  );
};

export default Ring;