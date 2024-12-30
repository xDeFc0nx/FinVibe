import { WebSocketClient } from "@/libs/socket";
import { createContext, createSignal, useContext, type JSX } from "solid-js";

interface WebSocketContextType {
  socket: WebSocketClient | null;
}

const WebSocketContext = createContext<WebSocketContextType | undefined>(
  undefined
);

export const WebSocketProvider = (props: { children: JSX.Element }) => {
  const [socket, setSocket] = createSignal<WebSocketClient | null>(null);

  if (!socket()) {
    const newSocket = new WebSocketClient("ws://your-websocket-url");
    setSocket(newSocket); // Initialize WebSocket connection
  }

  return (
    <WebSocketContext.Provider value={{ socket: socket() }}>
      {props.children}
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
