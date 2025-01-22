import type React from 'react';
import { WebSocketClient } from '@/lib/socket';
import { createContext, useContext, useEffect, useMemo, useState } from 'react';

interface WebSocketContextType {
  socket: WebSocketClient | null;
  isReady: boolean;
}
const WebSocketContext = createContext<WebSocketContextType | undefined>(
  undefined,
);

interface WebSocketProviderProps {
  children: React.ReactNode;
}

export const WebSocketProvider: React.FC<WebSocketProviderProps> = ({
  children,
}) => {
  const [socket, setSocket] = useState<WebSocketClient | null>(null);
  const [isReady, setReady] = useState(false);

  useEffect(() => {
    const newSocket = new WebSocketClient('ws://localhost:3001/ws');

    setSocket(newSocket);

    newSocket.socket.onopen = () => {
      setReady(true);
    };

    return () => {
      newSocket.close();
    };
  }, []);

  const value = useMemo(() => ({ socket, isReady }), [socket, isReady]);
  return (
    <WebSocketContext.Provider value={value}>
      {children} {/* Render children here */}
    </WebSocketContext.Provider>
  );
};

export const useWebSocket = (): WebSocketContextType => {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error('useWebSocket must be used within a WebSocketProvider');
  }
  return context;
};
