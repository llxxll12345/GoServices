import React, { useState, useEffect } from 'react';

const SimpleClient = () => {
  const [socket, setSocket] = useState(null);
  const [userPort, setUserPort] = useState(prompt('Enter your user port:'));
  const [targetUserPort, setTargetUserPort] = useState('');
  const [message, setMessage] = useState('');
  const [receivedMessages, setReceivedMessages] = useState([]);

  useEffect(() => {
    const newSocket = new WebSocket(`ws://localhost:8080/ws?port=${userPort}`);
    setSocket(newSocket);

    newSocket.onopen = (event) => {
      console.log('WebSocket connection opened:', event);
      newSocket.send(`Hello from User ${userPort}!`);
    };

    newSocket.onmessage = (event) => {
      console.log('Received message:', event.data);
      setReceivedMessages((prevMessages) => [...prevMessages, event.data]);
    };

    newSocket.onclose = (event) => {
      console.log('WebSocket connection closed:', event);
    };

    return () => {
      newSocket.close();
    };
  }, [userPort]);

  const sendMessage = () => {
    if (targetUserPort && message) {
      socket.send(`To: ${targetUserPort} ${message}`);
    }
  };

  return (
    <div>
      <h2>User {userPort}</h2>
      <div>
        <label>
          Target User Port:
          <input type="text" value={targetUserPort} onChange={(e) => setTargetUserPort(e.target.value)} />
        </label>
      </div>
      <div>
        <label>
          Message:
          <input type="text" value={message} onChange={(e) => setMessage(e.target.value)} />
        </label>
      </div>
      <div>
        <button onClick={sendMessage}>Send Message</button>
      </div>
      <div>
        <h3>Received Messages:</h3>
        <ul>
          {receivedMessages.map((msg, index) => (
            <li key={index}>{msg}</li>
          ))}
        </ul>
      </div>
    </div>
  );
};

export default SimpleClient;