export class WebSocketClient {
	public socket: WebSocket;
	public messageHandlers: ((msg: string) => void)[] = [];
	static send: any;

	constructor(url: string) {
		this.socket = new WebSocket(url);

		this.socket.onmessage = (event) => {
			this.messageHandlers.forEach((handler) => handler(event.data));
			const response = JSON.parse(event.data);
			console.log(response);
			if (response.message === "pong") {
				if (this.socket.readyState === WebSocket.OPEN) {
					this.socket.send(JSON.stringify({ Action: "pong" }));
				} else {
					console.error("Socket is not open. Cannot send pong.");
				}
			}
		};
	}

	onMessage(handler: (msg: string) => void) {
		this.messageHandlers.push(handler);
	}

	send(action: string, data: any = null) {
		if (this.socket.readyState === WebSocket.OPEN) {
			const message: { Action: string; Data?: any } = { Action: action };
			if (data) {
				message.Data = data;
			}
			this.socket.send(JSON.stringify(message));
			console.log(message);
		} else {
			console.error("WebSocket is not open.");
		}
	}
	close() {
		this.socket.close();
	}
}
