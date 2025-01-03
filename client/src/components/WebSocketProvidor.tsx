import { WebSocketClient } from "@/lib/socket";
import type React from "react";
import { createContext, useContext, useEffect, useState } from "react";

interface WebSocketContextType {
  socket: WebSocketClient | null;
}

const WebSocketContext = createContext<WebSocketContextType | undefined>(
  undefined
);

export const WebSocketProvider: React.FC = () => {
  const [socket, setSocket] = useState<WebSocketClient | null>(null);

  useEffect(() => {
    const newSocket = new WebSocketClient("ws://localhost:3001/ws");
    setSocket(newSocket);

    return () => {
      newSocket.close();
    };
  }, []);

  return (
    <WebSocketContext.Provider value={{ socket }}>
      {/* No need to render children */}
    </WebSocketContext.Provider>
  );
};

export const useWebSocket = (): WebSocketClient | null => {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error("useWebSocket must be used within a WebSocketProvider");
  }
  return context.socket;
};
