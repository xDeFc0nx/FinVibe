export class WebSocketClient {
  public socket: WebSocket;
  public messageHandlers: ((msg: string) => void)[] = [];
  static send: any;

  constructor(url: string) {
    this.socket = new WebSocket(url);

    this.socket.onmessage = (event) => {
      this.messageHandlers.forEach((handler) => handler(event.data));
        const message = JSON.parse(event.data);
      if (message === "ping") {
            if (this.socket.readyState === WebSocket.OPEN) {
        this.socket.send(JSON.stringify({ Action: "pong" }));
            }else {
    console.error("Socket is not open. Cannot send pong.");
  }
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
