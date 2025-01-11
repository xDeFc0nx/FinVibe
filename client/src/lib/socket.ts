export class WebSocketClient {
  public socket: WebSocket;
  public messageHandlers: ((msg: string) => void)[] = [];
  static send: any;

  constructor(url: string) {
    this.socket = new WebSocket(url);

    this.socket.onmessage = (event) => {
      this.messageHandlers.forEach((handler) => handler(event.data));
      if (event.data === "ping") {
        this.socket.send(JSON.stringify({ Action: "pong" }));
      }
    };
  }

  onMessage(handler: (msg: string) => void) {
    this.messageHandlers.push(handler);
  }

  send(message: string) {
    if (this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify({ Action: message}));
    } else {
      console.error("WebSocket is not open.");
    }
  }

  close() {
    this.socket.close();
  }
}
