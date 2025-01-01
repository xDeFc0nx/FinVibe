import { WebSocketClient } from "@/lib/socket";
import {
  createContext,
  createSignal,
  useContext,
  onCleanup,
  createEffect,
  JSX,
} from "solid-js";

interface WebSocketContextType {
  socket: WebSocketClient | null;
}

const WebSocketContext = createContext<WebSocketContextType | undefined>(
  undefined
);

export const WebSocketProvider = (props: { children?: JSX.Element }) => {
  const [socket, setSocket] = createSignal<WebSocketClient | null>(null);

  createEffect(() => {
    if (socket() === null) {
      const newSocket = new WebSocketClient("ws://localhost:3001/ws");
      setSocket(newSocket);
    }
  });

  onCleanup(() => {
    if (socket()) {
      socket()?.close();
    }
  });

  return (
    <WebSocketContext.Provider value={{ socket: socket() }}>
      {<></>}
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
